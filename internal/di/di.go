package di

import (
	"jadwalin-auth-events-integrator/internal/repository"
	"jadwalin-auth-events-integrator/internal/service"
	"jadwalin-auth-events-integrator/internal/shared/config"

	"go.uber.org/dig"
)

var (
	Container = dig.New()
)

func init() {
	if err := Container.Provide(config.NewConfig); err != nil {
		panic(err)
	}

	if err := Container.Provide(NewLogger); err != nil {
		panic(err)
	}

	// TODO: Add another providers

	if err := Container.Provide(NewEcho); err != nil {
		panic(err)
	}

	if err := repository.Register(Container); err != nil {
		panic(err)
	}

	if err := service.Register(Container); err != nil {
		panic(err)
	}
}
