package repository

import "go.uber.org/dig"

type Holder struct {
	dig.In
	Test Test
}

func Register(container *dig.Container) error {
	if err := container.Provide(NewTest); err != nil {
		return err
	}

	return nil
}
