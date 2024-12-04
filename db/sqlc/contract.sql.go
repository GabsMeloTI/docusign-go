// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: contract.sql

package db

import (
	"context"
)

const createContract = `-- name: CreateContract :one
INSERT INTO public.contracts
(id, provider_name, document_url, status, created_at, is_signed, date_signed)
VALUES(nextval('contracts_id_seq'::regclass), $1, $2, 'pending'::character varying, now(), false, NULL)
    RETURNING id, provider_name, document_url, status, created_at, is_signed, date_signed
`

type CreateContractParams struct {
	ProviderName string `json:"provider_name"`
	DocumentUrl  string `json:"document_url"`
}

func (q *Queries) CreateContract(ctx context.Context, arg CreateContractParams) (Contract, error) {
	row := q.db.QueryRowContext(ctx, createContract, arg.ProviderName, arg.DocumentUrl)
	var i Contract
	err := row.Scan(
		&i.ID,
		&i.ProviderName,
		&i.DocumentUrl,
		&i.Status,
		&i.CreatedAt,
		&i.IsSigned,
		&i.DateSigned,
	)
	return i, err
}