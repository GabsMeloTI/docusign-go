package contract

import (
	"context"
	"database/sql"
	db "docusign/db/sqlc"
)

type ContractInterfaceRepository interface {
	CreateContractRepository(ctx context.Context, arg db.CreateContractParams) (db.Contract, error)
}

type ContractRepository struct {
	DBtx    db.DBTX
	Queries *db.Queries
}

func NewContractRepository(sqlDB *sql.DB) *ContractRepository {
	q := db.New(sqlDB)
	return &ContractRepository{
		DBtx:    sqlDB,
		Queries: q,
	}
}

func (r *ContractRepository) CreateContractRepository(ctx context.Context, arg db.CreateContractParams) (db.Contract, error) {
	return r.Queries.CreateContract(ctx, arg)
}
