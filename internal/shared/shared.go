package shared

import (
	"context"
	"jadwalin-auth-events-integrator/internal/service"
	"jadwalin-auth-events-integrator/internal/shared/config"

	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
)

type Dependencies struct {
	dig.In

	Logger   log.Logger
	Services service.Holder

	Config    *config.Config
	Echo      *echo.Echo
	Scheduler *gocron.Scheduler
	DB        *mongo.Database
	Cache     *redis.Client

	//TODO: Add other dependencies
}

func (d *Dependencies) Close() {
	if err := d.Echo.Close(); err != nil {
		d.Logger.Errorf("failed to close echo server: %v", err)
	}

	if err := d.DB.Client().Disconnect(context.Background()); err != nil {
		d.Logger.Errorf("failed to close mongodb: %v", err)
	}

	if err := d.Cache.Close(); err != nil {
		d.Logger.Errorf("failed to close redis: %v", err)
	}

	// TODO: Close other dependencies
}
