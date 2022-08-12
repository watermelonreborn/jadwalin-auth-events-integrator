package service

import (
	"go.uber.org/dig"
)

type Holder struct {
	dig.In
	Test Test
	Auth Auth
	//TODO: Add other services
}

func Register(container *dig.Container) error {
	if err := container.Provide(NewTest); err != nil {
		return err
	}

	if err := container.Provide(OAuth); err != nil {
		return err
	}

	// TODO: Add another service providers

	return nil
}
