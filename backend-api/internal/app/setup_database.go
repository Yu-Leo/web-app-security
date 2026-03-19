package app

import (
	"database/sql"

	txmanager "github.com/Yu-Leo/web-app-security/backend-api/internal/storages/transaction/manager"
	txprovider "github.com/Yu-Leo/web-app-security/backend-api/internal/storages/transaction/provider"
)

func MustTxManager(db *sql.DB) *txmanager.Manager {
	return txmanager.New(db)
}

func MustTxProvider(db *sql.DB) *txprovider.Provider {
	return txprovider.New(db)
}
