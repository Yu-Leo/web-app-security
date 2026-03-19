package resources

import (
	"context"

	txmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
)

type database interface {
	ProvideTransaction(ctx context.Context) txmsql.Tr
}

//nolint:unused
type transaction interface {
	txmsql.Tr
}
