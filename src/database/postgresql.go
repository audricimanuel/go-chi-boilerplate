package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go-chi-boilerplate/src/tools"
)

const (
	DRIVER_POSTGRES = "postgres"
)

func InitializePostgresqlDatabase(ctx context.Context, dsn string) *sqlx.DB {
	db := tools.NewSqlxDsn(ctx, DRIVER_POSTGRES, dsn)
	return db
}
