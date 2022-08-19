package service

import (
	"go.uber.org/dig"
)

type Holder struct {
	dig.In
	Test  Test
	Auth  Auth
	Event Event
}

func Register(container *dig.Container) error {
	if err := container.Provide(NewTest); err != nil {
		return err
	}

	if err := container.Provide(NewAuth); err != nil {
		return err
	}

	if err := container.Provide(NewEvent); err != nil {
		return err
	}

	return nil
}
