package service

import (
	"context"
	"go-chi-boilerplate/src/internals/repository"
	"go-chi-boilerplate/src/model"
)

type (
	ExampleService interface {
		GetExample(ctx context.Context) model.ExampleResponse
	}

	ExampleServiceImpl struct {
		exampleRepo repository.ExampleRepository
	}
)

func NewExampleService(e repository.ExampleRepository) ExampleService {
	return &ExampleServiceImpl{
		exampleRepo: e,
	}
}

func (s *ExampleServiceImpl) GetExample(ctx context.Context) model.ExampleResponse {
	return s.exampleRepo.GetExample(ctx)
}
