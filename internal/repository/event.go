package repository

import (
	"jadwalin-auth-events-integrator/internal/entity"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UserEventsCollection = "user_events"
)

type (
	Event interface {
		UpsertEvents(entity.UserEvents) error
		GetAllEvents(string, string) ([]entity.UserEvents, error)
		GetUserEvents(string) ([]entity.Event, error)
	}

	eventRepo struct {
		logger log.Logger
		db     *mongo.Database
		redis  *redis.Client
	}
)

func (repo *eventRepo) UpsertEvents(userEvents entity.UserEvents) error {
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"_id": userEvents.ID}

	_, err := repo.db.Collection(UserEventsCollection).ReplaceOne(ctx, filter, userEvents, opts)
	if err != nil {
		repo.logger.Errorf("Error update user events : %s", err)
		return err
	}

	repo.logger.Info("Update user events success")

	return nil
}

func (repo *eventRepo) GetAllEvents(timeNow, timeHour string) ([]entity.UserEvents, error) {
	pipeline := mongo.Pipeline{
		{primitive.E{Key: "$unwind", Value: "$events"}},
		{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "events.start_time.date_time", Value: bson.D{primitive.E{Key: "$gt", Value: timeNow}, primitive.E{Key: "$lte", Value: timeHour}}}}}},
		{primitive.E{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, primitive.E{Key: "events", Value: bson.D{primitive.E{Key: "$push", Value: "$events"}}}}}},
	}

	cursor, err := repo.db.Collection(UserEventsCollection).Aggregate(ctx, pipeline)

	if err != nil {
		repo.logger.Errorf("Error find user events : %s", err)
		return nil, err
	}

	var result []entity.UserEvents
	if err = cursor.All(ctx, &result); err != nil {
		repo.logger.Errorf("Error to binding user events : %s", err)
		return nil, err
	}

	return result, nil
}

func (repo *eventRepo) GetUserEvents(userID string) ([]entity.Event, error) {
	result := repo.db.Collection(UserEventsCollection).FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userID}})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			repo.logger.Infof("User is not exist: %s", err)
			return nil, err
		}

		repo.logger.Errorf("Error: %s", err)
		return nil, err
	}

	var user_events entity.UserEvents
	if err := result.Decode(&user_events); err != nil {
		repo.logger.Errorf("Error decode result query to user_events: %s", err)
		return nil, err
	}

	return user_events.Events, nil
}

func NewEvent(logger log.Logger, db *mongo.Database, redis *redis.Client) (Event, error) {
	return &eventRepo{
		logger: logger,
		db:     db,
		redis:  redis}, nil
}
