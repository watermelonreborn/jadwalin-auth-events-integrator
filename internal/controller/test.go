package controller

import (
	"jadwalin-auth-events-integrator/internal/service"
	"jadwalin-auth-events-integrator/internal/shared/dto"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type Test struct {
	dig.In
	Service service.Holder
}

func (controller *Test) Index(c echo.Context) error {
	var (
		request = dto.TestRequest{}
	)

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	response, err := controller.Service.Test.Index(c, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, response)
}
