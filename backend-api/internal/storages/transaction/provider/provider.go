package provider

import (
	"context"
	"database/sql"

	txmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
)

type Provider struct {
	db *sql.DB
}

func New(db *sql.DB) *Provider {
	return &Provider{
		db: db,
	}
}

func (p *Provider) ProvideTransaction(ctx context.Context) txmsql.Tr {
	return txmsql.DefaultCtxGetter.DefaultTrOrDB(ctx, p.db)
}
