package repository

import "go.uber.org/dig"

type Holder struct {
	dig.In
	Test Test

	//TODO: Add other repositories
}

func Register(container *dig.Container) error {
	if err := container.Provide(NewTest); err != nil {
		return err
	}

	// TODO: Add another repository providers

	return nil
}
