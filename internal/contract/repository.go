package contract

import (
	"context"
	"database/sql"
	db "docusign/db/sqlc"
)

type ContractInterfaceRepository interface {
	CreateContractRepository(context.Context, db.CreateContractParams) (db.Contract, error)
	AssignedContract(context.Context, string) (db.Contract, error)
	GetContractAll(ctx context.Context, arg db.GetContractAllParams) (db.Contract, error)
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

func (r *ContractRepository) AssignedContract(ctx context.Context, envelopID string) (db.Contract, error) {
	return r.Queries.AssignContract(ctx, envelopID)
}

func (r *ContractRepository) GetContractAll(ctx context.Context, arg db.GetContractAllParams) (db.Contract, error) {
	return r.Queries.GetContractAll(ctx, arg)
}
