package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-chi-boilerplate/src/database"
	"go-chi-boilerplate/src/model"
)

type (
	ExampleRepository interface {
		GetExample(ctx context.Context) model.ExampleResponse
	}

	ExampleRepositoryImpl struct {
		db database.DBCollection
	}
)

func NewExampleRepository(db database.DBCollection) ExampleRepository {
	return &ExampleRepositoryImpl{
		db: db,
	}
}

func (e *ExampleRepositoryImpl) GetExample(ctx context.Context) model.ExampleResponse {
	logrus.WithContext(ctx).Info("GetExample")
	return model.ExampleResponse{
		AppName: "Example",
		Env:     "DEV",
	}
}
