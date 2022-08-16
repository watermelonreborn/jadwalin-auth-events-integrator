package controller

import (
	"jadwalin-auth-events-integrator/internal/service"
	"jadwalin-auth-events-integrator/internal/shared/dto"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type Auth struct {
	dig.In
	Service service.Holder
	Logger  log.Logger
}

func (impl *Auth) handleAuth(c echo.Context) error {
	url := impl.Service.Auth.URL()
	impl.Logger.Info("Auth URL Generated: ", url)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (impl *Auth) handleAuthCallback(c echo.Context) error {
	state := c.QueryParam("state")
	code := c.QueryParam("code")

	impl.Logger.Info("Auth Callback: ", state, code)

	token, err := impl.Service.Auth.GenerateToken(state, code)
	if err != nil {
		impl.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, dto.Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	impl.Logger.Info("Auth Token Generated")

	return c.JSON(http.StatusOK, dto.Response{
		Status: http.StatusOK,
		Data:   token,
	})
}

func (impl *Auth) handleUserInfo(c echo.Context) error {
	var (
		request = dto.UserInfoRequest{}
	)

	if err := c.Bind(&request); err != nil {
		impl.Logger.Errorf("Error binding request: %s", err)
		return c.JSON(http.StatusBadRequest, dto.Response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	response, err := impl.Service.Auth.GetUserInfo(request.Token)
	if err != nil {
		impl.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, dto.Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	impl.Logger.Info("Got the user info from the token")

	return c.JSON(http.StatusOK, dto.Response{
		Status: http.StatusOK,
		Data:   response,
	})
}
