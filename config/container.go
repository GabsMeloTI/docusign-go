package config

import (
	"database/sql"
	"docusign/infra/database"
	"docusign/infra/database/db_postgresql"
	"docusign/internal/contract"
)

type ContainerDI struct {
	Config             Config
	Conn               *sql.DB
	HandlerContract    *contract.ContractHandler
	WebhookHandler     *contract.WebhookHandler
	ServiceContract    *contract.ContractService
	RepositoryContract *contract.ContractRepository
}

func NewContainerDI(config Config) *ContainerDI {
	container := &ContainerDI{Config: config}

	container.db()
	container.buildRepository()
	container.buildService()
	container.buildHandlers()

	return container
}

func (c *ContainerDI) db() {
	dbConfig := database.Config{
		Host:        c.Config.DBHost,
		Port:        c.Config.DBPort,
		User:        c.Config.DBUser,
		Password:    c.Config.DBPassword,
		Database:    c.Config.DBDatabase,
		SSLMode:     c.Config.DBSSLMode,
		Driver:      c.Config.DBDriver,
		Environment: c.Config.Environment,
	}
	c.Conn = db_postgresql.NewConnection(&dbConfig)
}

func (c *ContainerDI) buildRepository() {
	c.RepositoryContract = contract.NewContractRepository(c.Conn)
}

func (c *ContainerDI) buildService() {
	c.ServiceContract = contract.NewContractService(c.RepositoryContract)
}

func (c *ContainerDI) buildHandlers() {
	c.HandlerContract = contract.NewContractHandler(c.ServiceContract)
	c.WebhookHandler = contract.NewWebhookHandler(*c.ServiceContract)
}
