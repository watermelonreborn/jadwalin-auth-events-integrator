package service

import (
	"fmt"
	"jadwalin-auth-events-integrator/internal/entity"
	"jadwalin-auth-events-integrator/internal/repository"
	"jadwalin-auth-events-integrator/internal/shared/dto"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type (
	Test interface {
		Index(echo.Context, dto.TestRequest) (dto.TestResponse, error)
	}

	testService struct {
		logger     log.Logger
		repository repository.Holder
	}
)

func (service *testService) Index(c echo.Context, request dto.TestRequest) (dto.TestResponse, error) {
	var (
		response = dto.TestResponse{}
	)

	err := service.repository.Test.IndexDB(entity.Test{Name: request.Name})
	if err != nil {
		return response, err
	}

	err = service.repository.Test.IndexRedis(request.Name)
	if err != nil {
		return response, err
	}

	response.Message = fmt.Sprintf("Hello %s!", request.Name)

	return response, nil
}

func NewTest(logger log.Logger, repo repository.Holder) (Test, error) {
	return &testService{
		logger:     logger,
		repository: repo,
	}, nil
}
