package repository

import (
	"context"
	"jadwalin-auth-events-integrator/internal/entity"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ctx = context.Background()
)

const (
	UserTokenCollection = "user_token"
)

type (
	Auth interface {
		AddUser(entity.User) error
		IsUserExist(string) (bool, error)
		GetToken(string) (string, error)
		GetAllUserToken() ([]entity.User, error)
	}

	authRepo struct {
		logger log.Logger
		db     *mongo.Database
		redis  *redis.Client
	}
)

func (repo *authRepo) AddUser(user entity.User) error {
	_, err := repo.db.Collection(UserTokenCollection).InsertOne(ctx, user)
	if err != nil {
		repo.logger.Errorf("Error add user : %s", err)
		return err
	}

	repo.logger.Info("Add user success")

	return nil
}

func (repo *authRepo) IsUserExist(userId string) (bool, error) {
	result := repo.db.Collection(UserTokenCollection).FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userId}})
	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			repo.logger.Infof("User is not exist: %s", err)
			return false, nil
		}
		repo.logger.Errorf("Error: %s", err)
		return false, err
	}

	repo.logger.Info("Get user success")
	return true, nil
}

func (repo *authRepo) GetToken(userId string) (string, error) {
	result := repo.db.Collection(UserTokenCollection).FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userId}})
	if err := result.Err(); err != nil {
		repo.logger.Errorf("Error get token from db: %s", err)
		return "", err
	}

	var user entity.User
	if err := result.Decode(&user); err != nil {
		repo.logger.Errorf("Error Decode to user: %s", err)
		return "", err
	}

	repo.logger.Info("Get user success")
	return user.RefreshToken, nil
}

func (repo *authRepo) GetAllUserToken() ([]entity.User, error) {
	cursor, err := repo.db.Collection(UserTokenCollection).Find(ctx, bson.D{})
	if err != nil {
		repo.logger.Errorf("Error get all user token from db: %s", err)
		return nil, err
	}

	var usersToken []entity.User
	if err := cursor.All(ctx, &usersToken); err != nil {
		repo.logger.Errorf("Error decode to define struct: %s", err)
		return nil, err
	}

	return usersToken, nil
}

func NewAuth(logger log.Logger, db *mongo.Database, redis *redis.Client) (Auth, error) {
	return &authRepo{
		logger: logger,
		db:     db,
		redis:  redis}, nil
}
