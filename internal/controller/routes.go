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
	Auth         Auth
}

func (impl *Holder) RegisterRoutes() {
	var app = impl.Dependencies.Echo

	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	app.POST("/", impl.Test.Index)
	app.GET("/auth", impl.Auth.handleAuth)
	app.GET("/callback", impl.Auth.handleAuthCallback)
	app.GET("/userinfo", impl.Auth.handleUserInfo)
}
