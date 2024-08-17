package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-chi-boilerplate/src/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	DatabaseCollection struct {
		MongoDB      *mongo.Database
		PostgresqlDB *sqlx.DB
	}
)

func NewDatabaseCollection(cfg config.Config) DatabaseCollection {
	ctx := context.Background()

	// mongodb
	mongoDBConfig := cfg.DataSource.MongoDBConfig
	mongoDB := InitializeMongoDatabase(ctx, mongoDBConfig.ConnectionString, mongoDBConfig.DatabaseName)

	// postgres
	postgresDBConfig := cfg.DataSource.PostgresDBConfig
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		postgresDBConfig.Host, postgresDBConfig.User, postgresDBConfig.Password, postgresDBConfig.Name, postgresDBConfig.Port, postgresDBConfig.SSLMode, postgresDBConfig.Timezone,
	)
	postgresDB := InitializePostgresqlDatabase(ctx, dsn)

	return DatabaseCollection{
		MongoDB:      mongoDB,
		PostgresqlDB: postgresDB,
	}
}
