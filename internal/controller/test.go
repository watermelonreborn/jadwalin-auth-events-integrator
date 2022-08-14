package controller

import (
	"jadwalin-auth-events-integrator/internal/service"
	"jadwalin-auth-events-integrator/internal/shared/dto"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type Test struct {
	dig.In
	Service service.Holder
	Logger  log.Logger
}

func (controller *Test) Index(c echo.Context) error {
	var (
		request = dto.TestRequest{}
	)
	controller.Logger.Info("masuk controller test")
	if err := c.Bind(&request); err != nil {
		controller.Logger.Errorf("Error binding request: %s", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	response, err := controller.Service.Test.Index(c, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, response)
}
