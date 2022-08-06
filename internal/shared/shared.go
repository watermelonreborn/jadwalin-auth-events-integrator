package shared

import (
	"jadwalin-auth-events-integrator/internal/shared/config"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type Dependencies struct {
	dig.In
	Config *config.Config
	Logger log.Logger
	Echo   *echo.Echo

	//TODO: Add other dependencies
}

func (d *Dependencies) Close() {
	if err := d.Echo.Close(); err != nil {
		d.Logger.Errorf("failed to close echo server: %v", err)
	}

	// TODO: Close other dependencies
}
