package controller

import (
	"jadwalin-auth-events-integrator/internal/shared"

	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

type Holder struct {
	dig.In
	Dependencies shared.Dependencies
	Test         Test
}

func (impl *Holder) RegisterRoutes() {
	var app = impl.Dependencies.Echo

	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	app.POST("/", impl.Test.Index)

	//TODO: Add other routes
}
