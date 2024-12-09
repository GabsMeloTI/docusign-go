package contract

import (
	db "docusign/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type ContractRequestCreate struct {
	ProviderName   string    `json:"provider_name"`
	ProviderEmail  string    `json:"provider_email"`
	ContractType   string    `json:"contract_type" validate:"required,oneof=assignment_contract term_of_assignment"`
	IdBatchControl uuid.UUID `json:"id_batch_control"`
	TenantId       uuid.UUID `json:"tenant_id"`
	AccessId       int64     `json:"access_id"`
}

type ContractResponse struct {
	ID           int64     `json:"id"`
	ProviderName string    `json:"provider_name"`
	DocumentUrl  string    `json:"document_url"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	IsSigned     bool      `json:"is_signed"`
	ContractType string    `json:"contract_type"`
}

func (c *ContractResponse) ParseFromContractObject(result db.Contract) {
	c.ID = result.ID
	c.ProviderName = result.ProviderName
	c.DocumentUrl = result.DocumentUrl
	c.Status = result.Status
	c.CreatedAt = result.CreatedAt.Time
	c.IsSigned = result.IsSigned.Bool
	c.ContractType = result.ContractType
}
