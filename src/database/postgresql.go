package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go-chi-boilerplate/src/tools"
	"gorm.io/gorm"
)

const (
	DRIVER_POSTGRES = "postgres"
)

func InitializePostgresqlDatabaseSqlx(ctx context.Context, dsn string) *sqlx.DB {
	db := tools.NewSqlxDsn(ctx, DRIVER_POSTGRES, dsn)
	return db
}

func InitializePostgresqlDatabaseGorm(ctx context.Context, dsn string) *gorm.DB {
	db := tools.NewGormDB(ctx, dsn)
	return db
}
