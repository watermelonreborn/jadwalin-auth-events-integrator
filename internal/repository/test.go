package repository

import (
	"jadwalin-auth-events-integrator/internal/entity"

	log "github.com/sirupsen/logrus"
)

type (
	Test interface {
		Index(entity.Test) (entity.Test, error)
	}

	testRepo struct {
		logger log.Logger
		// TODO: Add db connection
	}
)

func (repo *testRepo) Index(test entity.Test) (entity.Test, error) {
	repo.logger.Infof("Indexing test: %s", test.Name)
	return test, nil
}

func NewTest(logger log.Logger) (Test, error) {
	return &testRepo{logger: logger}, nil
}
