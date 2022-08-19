package service

import (
	"jadwalin-auth-events-integrator/internal/repository"

	log "github.com/sirupsen/logrus"
)

type (
	Event interface {
		SyncAPIWithDB() error
	}

	eventService struct {
		logger     log.Logger
		repository repository.Holder
	}
)

func (service *eventService) SyncAPIWithDB() error {
	return nil
}

func NewEvent(logger log.Logger, repo repository.Holder) (Event, error) {
	return &eventService{
		logger:     logger,
		repository: repo,
	}, nil
}
