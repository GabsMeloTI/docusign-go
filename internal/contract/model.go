package contract

import (
	db "docusign/db/sqlc"
	"time"
)

type ContractRequestCreate struct {
	ProviderName string `json:"provider_name"`
}

type ContractResponse struct {
	ID           int32     `json:"id"`
	ProviderName string    `json:"provider_name"`
	DocumentUrl  string    `json:"document_url"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	IsSigned     bool      `json:"is_signed"`
}

func (c *ContractResponse) ParseFromContractObject(result db.Contract) {
	c.ID = result.ID
	c.ProviderName = result.ProviderName
	c.DocumentUrl = result.DocumentUrl
	c.Status = result.Status
	c.CreatedAt = result.CreatedAt.Time
	c.IsSigned = result.IsSigned.Bool
}
