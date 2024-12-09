package contract

import (
	"context"
	"database/sql"
	db "docusign/db/sqlc"
	"docusign/internal/model"
	"docusign/pkg"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
)

type ContractInterfaceService interface {
	CreateContractService(context.Context, ContractRequestCreate) (ContractResponse, error)
	AssignContractService(context.Context, string) (ContractResponse, error)
}

type ContractService struct {
	ContractInterface ContractInterfaceRepository
}

func NewContractService(ContractInterface ContractInterfaceRepository) *ContractService {
	return &ContractService{ContractInterface}
}

func (s *ContractService) CreateContractService(ctx context.Context, data ContractRequestCreate) (ContractResponse, error) {
	unicId := uuid.New()

	var valueBatch uuid.UUID
	if data.ContractType == "term_of_assignment" {
		valueBatch = uuid.Nil
	} else {
		if data.IdBatchControl == uuid.Nil {
			return ContractResponse{}, errors.New("id_batch_control must be provided for this contract type")
		}
		valueBatch = data.IdBatchControl
	}

	templatePath := "assets/" + data.ContractType + ".html"
	outputFile := fmt.Sprintf("contract-%s.pdf", unicId.String())
	dataPdf := map[string]interface{}{
		"name_enterprise": data.ProviderName,
		"date":            time.Now().Format("02 de Janeiro de 2006"),
	}
	err := pkg.GeneratePDF(templatePath, dataPdf, outputFile)
	if err != nil {
		fmt.Printf("Erro ao gerar o PDF: %v\n", err)
	} else {
		fmt.Println("PDF gerado com sucesso:", outputFile)
	}

	jwtToken, err := pkg.GenerateJWT(os.Getenv("DOCUSIGN_APIKEY"), os.Getenv("DOCUSIGN_USERNAME"), strings.ReplaceAll(os.Getenv("DOCUSIGN_RSA_PRIVATE_KEY"), "\\n", "\n"))
	if err != nil {
		return ContractResponse{}, nil
	}

	accessToken, err := pkg.GetAccessToken(jwtToken)
	if err != nil {
		return ContractResponse{}, err
	}

	_, err = pkg.UploadFileToS3(fmt.Sprintf("contract-%s.pdf", unicId.String()), "contracts-guapi/"+data.ContractType+"/")
	if err != nil {
		return ContractResponse{}, fmt.Errorf("error uploading to S3: %v", err)
	}

	encodedContent, err := pkg.EncodeFileToBase64(fmt.Sprintf("contract-%s.pdf", unicId.String()))
	if err != nil {
		return ContractResponse{}, err
	}

	dataEnvelope := model.EnvelopeDefinition{
		EmailSubject: "Contrato",
		Documents: []model.Document{
			{
				DocumentBase64: encodedContent,
				Name:           "Contrato.pdf",
				FileExtension:  "pdf",
				DocumentID:     "1",
			},
		},
		Recipients: model.Recipients{
			Signers: []model.Signer{
				{
					Email:        data.ProviderEmail,
					Name:         data.ProviderName,
					RecipientID:  "1",
					RoutingOrder: "1",
					Tabs: model.Tabs{
						SignHereTabs: []model.SignHereTab{
							{
								DocumentID:  "1",
								PageNumber:  "1",
								RecipientID: "1",
								TabLabel:    "AssinaturaAqui",
								XPosition:   "30",
								YPosition:   "480",
							},
						},
					},
				},
			},
		},
		Status: "sent",
	}

	resultEnv, err := pkg.SendEnvelope(accessToken, os.Getenv("DOCUSIGN_ACCTID"), dataEnvelope)
	if err != nil {
		return ContractResponse{}, err
	}

	arg := db.CreateContractParams{
		ProviderName:  data.ProviderName,
		ProviderEmail: data.ProviderEmail,
		ContractType:  data.ContractType,
		IDControlBatch: uuid.NullUUID{
			UUID:  valueBatch,
			Valid: true,
		},
		DocumentUrl: "https://contracts-guapi.s3.us-east-1.amazonaws.com/terms_of_authorization/contract-" + unicId.String() + ".pdf",
		EnvelopID:   resultEnv,
		AccessID: sql.NullInt64{
			Int64: data.AccessId,
			Valid: true,
		},
		TenantID: uuid.NullUUID{
			UUID:  data.TenantId,
			Valid: true,
		},
	}

	result, err := s.ContractInterface.CreateContractRepository(ctx, arg)
	if err != nil {
		return ContractResponse{}, err
	}

	createContract := ContractResponse{}
	createContract.ParseFromContractObject(result)

	return createContract, err
}

func (s *ContractService) AssignContractService(ctx context.Context, envelopeID string) (ContractResponse, error) {
	_, err := s.ContractInterface.AssignedContract(ctx, envelopeID)
	return ContractResponse{}, err
}
