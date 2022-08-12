package repository

import (
	"context"
	"jadwalin-auth-events-integrator/internal/entity"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Auth interface {
		IndexDBAuth(entity.User) error
	}

	authRepo struct {
		logger log.Logger
		db     *mongo.Database
	}
)

func (repo *authRepo) IndexDBAuth(test entity.User) error {
	_, err := repo.db.Collection(CollectionName).InsertOne(context.Background(), test)
	if err != nil {
		repo.logger.Errorf("Error index db: %s", err)
		return err
	}

	repo.logger.Info("Index db success")

	return nil
}

func OAuth(logger log.Logger, db *mongo.Database) (Auth, error) {
	return &authRepo{logger: logger, db: db}, nil
}
