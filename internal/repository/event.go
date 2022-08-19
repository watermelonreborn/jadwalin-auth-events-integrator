package repository

import (
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Event interface {
	}

	eventRepo struct {
		logger log.Logger
		db     *mongo.Database
		redis  *redis.Client
	}
)

func NewEvent(logger log.Logger, db *mongo.Database, redis *redis.Client) (Event, error) {
	return &eventRepo{
		logger: logger,
		db:     db,
		redis:  redis}, nil
}
