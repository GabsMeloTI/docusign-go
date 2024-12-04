package contract

import (
	"context"
	db "docusign/db/sqlc"
	"docusign/internal/model"
	"docusign/pkg"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

type ContractInterfaceService interface {
	CreateContractService(context.Context, ContractRequestCreate) (ContractResponse, error)
}

type ContractService struct {
	ContractInterface ContractInterfaceRepository
}

func NewContractService(ContractInterface ContractInterfaceRepository) *ContractService {
	return &ContractService{ContractInterface}
}

func (s *ContractService) CreateContractService(ctx context.Context, data ContractRequestCreate) (ContractResponse, error) {
	unicId := uuid.New()

	arg := db.CreateContractParams{
		ProviderName: data.ProviderName,
		DocumentUrl:  "https://contracts-guapi.s3.us-east-1.amazonaws.com/terms_of_authorization/contract-" + unicId.String() + ".pdf",
	}

	result, err := s.ContractInterface.CreateContractRepository(ctx, arg)
	if err != nil {
		return ContractResponse{}, err
	}

	err = pkg.GeneratePDF(arg.ProviderName, fmt.Sprintf("contract-%s.pdf", unicId.String()))
	if err != nil {
		return ContractResponse{}, fmt.Errorf("error generating PDF: %v", err)
	}

	jwtToken, err := pkg.GenerateJWT(os.Getenv("DOCUSIGN_APIKEY"), os.Getenv("DOCUSIGN_USERNAME"), strings.ReplaceAll(os.Getenv("DOCUSIGN_RSA_PRIVATE_KEY"), "\\n", "\n"))
	if err != nil {
		return ContractResponse{}, nil
	}

	accessToken, err := pkg.GetAccessToken(jwtToken)
	if err != nil {
		return ContractResponse{}, err
	}

	_, err = pkg.UploadFileToS3(fmt.Sprintf("contract-%s.pdf", unicId.String()), "contracts-guapi/terms_of_authorization/")
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
					Email:        "bielmelodossantos8@gmail.com",
					Name:         "Guilherme Zoio Da Sul",
					RecipientID:  "1",
					RoutingOrder: "1",
					Tabs: model.Tabs{
						SignHereTabs: []model.SignHereTab{
							{
								DocumentID:  "1",
								PageNumber:  "1",
								RecipientID: "1",
								TabLabel:    "AssinaturaAqui",
								XPosition:   "100",
								YPosition:   "150",
							},
						},
					},
				},
			},
		},
		Status: "sent",
	}

	_, err = pkg.SendEnvelope(accessToken, os.Getenv("DOCUSIGN_ACCTID"), dataEnvelope)
	if err != nil {
		return ContractResponse{}, err
	}

	createContract := ContractResponse{}
	createContract.ParseFromContractObject(result)

	return createContract, err
}

//func SendContract(configs config.Config) {
//	jwtToken, err := generateJWT(configs.DocuSignApiKey, configs.DocuSignUsername, configs.DocuSignRSAKey)
//
//	if err != nil {
//		fmt.Println("Erro ao gerar JWT:", err)
//		return
//	}
//	fmt.Println("JWT Gerado:", jwtToken)
//
//	accessToken, err := getAccessToken(jwtToken)
//	if err != nil {
//		fmt.Println("Erro ao obter token de acesso:", err)
//		return
//	}
//
//	encodedContent, err := encodeFileToBase64("document.pdf")
//	if err != nil {
//		fmt.Println("Erro ao codificar documento:", err)
//		return
//	}
//
//	envelope := model.EnvelopeDefinition{
//		EmailSubject: "Contrato",
//		Documents: []model.Document{
//			{
//				DocumentBase64: encodedContent,
//				Name:           "Contrato.pdf",
//				FileExtension:  "pdf",
//				DocumentID:     "1",
//			},
//		},
//		Recipients: model.Recipients{
//			Signers: []model.Signer{
//				{
//					Email:        "bielmelodossantos8@gmail.com",
//					Name:         "Guilherme Zoio Da Sul",
//					RecipientID:  "1",
//					RoutingOrder: "1",
//					Tabs: model.Tabs{
//						SignHereTabs: []model.SignHereTab{
//							{
//								DocumentID:  "1",
//								PageNumber:  "1",
//								RecipientID: "1",
//								TabLabel:    "AssinaturaAqui",
//								XPosition:   "100",
//								YPosition:   "150",
//							},
//						},
//					},
//				},
//			},
//		},
//		Status: "sent",
//	}
//
//	envelopeID, err := sendEnvelope(accessToken, configs.DocuSignAccountId, envelope)
//	if err != nil {
//		fmt.Println("Erro ao enviar envelope:", err)
//		return
//	}
//
//	fmt.Println("Envelope enviado com sucesso. ID do Envelope:", envelopeID)
//}
