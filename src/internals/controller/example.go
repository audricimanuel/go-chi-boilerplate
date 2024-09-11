package controller

import (
	"go-chi-boilerplate/src/internals/service"
	"go-chi-boilerplate/utils/httputils"
	"net/http"
)

type (
	ExampleController interface {
		GetExample(w http.ResponseWriter, r *http.Request)
	}

	ExampleControllerImpl struct {
		exampleService service.ExampleService
	}
)

func NewExampleController(e service.ExampleService) ExampleController {
	return &ExampleControllerImpl{
		exampleService: e,
	}
}

// @Tags			Example
// @Summary		Example API
// @Description	"Just an example"
// @Accept			json
// @Produce		json
// @Success		200	{object}	httputils.BaseResponse
// @Router			/example [get]
func (e *ExampleControllerImpl) GetExample(w http.ResponseWriter, r *http.Request) {
	data := e.exampleService.GetExample(r.Context())
	httputils.MapBaseResponse(w, r, data, nil, nil)
}
