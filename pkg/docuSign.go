package pkg

import (
	"bytes"
	"docusign/internal/model"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/golang-jwt/jwt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func EncodeFileToBase64(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(content), nil
}

func GenerateJWT(integrationKey, userID, rsaPrivateKey string) (string, error) {
	claims := jwt.MapClaims{
		"iss":   integrationKey,
		"sub":   userID,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
		"aud":   "account-d.docusign.com",
		"scope": "signature impersonation",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(rsaPrivateKey))
	if err != nil {
		return "", fmt.Errorf("erro ao parsear chave privada: %v", err)
	}
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("erro ao assinar o JWT: %v", err)
	}
	return signedToken, nil
}

func GetAccessToken(jwtToken string) (string, error) {
	url := "https://account-d.docusign.com/oauth/token"
	data := map[string]string{
		"grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer",
		"assertion":  jwtToken,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("erro ao criar JSON de dados: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler corpo da resposta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erro na resposta do DocuSign: HTTP %d, corpo: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta JSON: %v", err)
	}

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("falha ao obter o token de acesso, resposta: %v", result)
}

func GeneratePDF(templateFilePath string, templateData interface{}, outputFile string) error {
	tmplContent, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		return fmt.Errorf("error reading template file: %v", err)
	}

	tmpl, err := template.New("contract").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	var renderedHTML bytes.Buffer
	err = tmpl.Execute(&renderedHTML, templateData)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("error initializing PDF generator: %v", err)
	}

	page := wkhtmltopdf.NewPageReader(strings.NewReader(renderedHTML.String()))
	pdfg.AddPage(page)
	pdfg.OutputFile = outputFile

	if err = pdfg.Create(); err != nil {
		return fmt.Errorf("error creating PDF: %v", err)
	}

	return nil
}

func SendEnvelope(accessToken, accountID string, envelope model.EnvelopeDefinition) (string, error) {
	url := "https://demo.docusign.net/restapi/v2.1/accounts/" + accountID + "/envelopes"
	jsonData, err := json.Marshal(envelope)
	if err != nil {
		return "", fmt.Errorf("failed to marshal envelope: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if len(body) == 0 {
		return "", fmt.Errorf("response body is empty")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	var envelopeID string
	if envelopeID, ok := result["envelopeId"].(string); ok {
		return envelopeID, nil
	}
	return envelopeID, fmt.Errorf("envelopeId not found in response: %s", string(body))
}
