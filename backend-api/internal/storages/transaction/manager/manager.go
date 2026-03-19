package manager

import (
	"context"
	"database/sql"

	txmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
	txm "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type Manager struct {
	txm *txm.Manager
}

func New(db *sql.DB) *Manager {
	return &Manager{
		txm: txm.Must(txmsql.NewDefaultFactory(db)),
	}
}

func (m *Manager) Do(ctx context.Context, fn func(context.Context) error) error {
	return m.txm.Do(ctx, fn)
}
