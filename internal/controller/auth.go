package controller

import (
	"jadwalin-auth-events-integrator/internal/entity"
	"jadwalin-auth-events-integrator/internal/service"
	"jadwalin-auth-events-integrator/internal/shared/dto"
	"net/http"

	"strings"

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

	userInfoDTO, err := impl.Service.Auth.GetUserInfo(token.AccessToken)
	if err != nil {
		impl.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, dto.Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	isUserExist, err := impl.Service.Auth.IsUserExist(userInfoDTO.ID)
	if err != nil {
		impl.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, dto.Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	if !isUserExist {
		impl.Service.Auth.AddUser(entity.User{
			ID:           userInfoDTO.ID,
			RefreshToken: token.RefreshToken,
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Status: http.StatusOK,
		Data:   token,
	})
}

func (impl *Auth) handleUserInfo(c echo.Context) error {
	authorizationHeaderValue := c.Request().Header.Get("Authorization")
	if authorizationHeaderValue == "" {
		errorMessage := "Error request: Authorization hasn't found on request header"
		impl.Logger.Errorf(errorMessage)
		return c.JSON(http.StatusBadRequest, dto.Response{
			Status: http.StatusBadRequest,
			Error:  errorMessage,
		})
	}

	tokenType := "Bearer"
	tokenStartIndex := strings.Index(authorizationHeaderValue, tokenType)
	token := authorizationHeaderValue[tokenStartIndex+(len(tokenType)+1):]

	response, err := impl.Service.Auth.GetUserInfo(token)
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
