package repository

import (
	"context"
	"jadwalin-auth-events-integrator/internal/entity"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CollectionName = "test"
)

type (
	Test interface {
		IndexDB(entity.Test) error
		IndexRedis(interface{}) error
	}

	testRepo struct {
		logger log.Logger
		db     *mongo.Database
		redis  *redis.Client
	}
)

func (repo *testRepo) IndexDB(test entity.Test) error {
	_, err := repo.db.Collection(CollectionName).InsertOne(context.Background(), test)
	if err != nil {
		repo.logger.Errorf("Error index db: %s", err)
		return err
	}

	repo.logger.Info("Index db success")

	return nil
}

func (repo *testRepo) IndexRedis(value interface{}) error {
	op := repo.redis.Set(context.Background(), "key", value, time.Minute)
	if err := op.Err(); err != nil {
		repo.logger.Errorf("Error set index value redis: %s", err)
		return err
	}

	op2 := repo.redis.Get(context.Background(), "key")
	if err := op2.Err(); err != nil {
		repo.logger.Errorf("Error get index value redis: %s", err)
		return err
	}

	repo.logger.Infof("Redis value: %s", op2.Val())

	return nil
}

func NewTest(logger log.Logger, db *mongo.Database, redis *redis.Client) (Test, error) {
	return &testRepo{logger: logger, db: db, redis: redis}, nil
}
