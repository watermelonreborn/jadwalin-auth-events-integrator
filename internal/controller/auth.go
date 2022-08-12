package controller

import (
	"jadwalin-auth-events-integrator/internal/service"
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

func (controller *Auth) IndexAuth(c echo.Context) error {

	controller.Logger.Info("masuk controller")
	// return c.Redirect(301, "https://accounts.google.com/o/oauth2/auth?access_type=offline\u0026client_id=304586193738-vfrl77vb8laoqh738tqku7fepfp5mi5c.apps.googleusercontent.com\u0026redirect_uri=http%3A%2F%2Flocalhost%3A8080%2F\u0026response_type=code\u0026scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcalendar.readonly\u0026state=state-token")
	url, err := controller.Service.Auth.IndexAuth(c)
	controller.Logger.Info("setelah service")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// if err := c.Redirect(301, url); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
	return c.Redirect(301, url)
}

func (controller *Auth) ParsingCode(c echo.Context) error {
	code := c.QueryParams().Get("code")
	controller.Logger.Info(code)

	return c.JSON(http.StatusOK, nil)
}
