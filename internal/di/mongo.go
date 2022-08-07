package di

import (
	"context"
	"jadwalin-auth-events-integrator/internal/shared/config"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo(config *config.Config, logger log.Logger) (*mongo.Database, error) {
	credential := options.Credential{
		Username: config.Database.Username,
		Password: config.Database.Password,
	}

	clientOptions := options.Client()
	clientOptions.ApplyURI(config.Database.ConnectionString)

	if config.Database.Username != "" && config.Database.Password != "" {
		clientOptions.SetAuth(credential)
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		logger.Errorf("Error new mongo client: %s", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logger.Errorf("Error connect mongo client: %s", err)
		return nil, err
	}

	return client.Database("jadwalin"), nil
}
