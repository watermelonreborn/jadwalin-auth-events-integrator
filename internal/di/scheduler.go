package di

import (
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func NewScheduler(logger log.Logger) *gocron.Scheduler {
	timeLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logger.Error("Unable load location from jakarta: %s", err.Error())
		timeLocation = time.UTC
	}

	return gocron.NewScheduler(timeLocation)
}
