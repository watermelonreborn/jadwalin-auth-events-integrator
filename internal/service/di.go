package service

import (
	"go.uber.org/dig"
)

type Holder struct {
	dig.In
	Test Test
	Auth Auth
}

func Register(container *dig.Container) error {
	if err := container.Provide(NewTest); err != nil {
		return err
	}

	if err := container.Provide(OAuth); err != nil {
		return err
	}

	return nil
}
