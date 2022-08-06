package di

import (
	log "github.com/sirupsen/logrus"
)

func NewLogger() (log.Logger, error) {
	var (
		logger = log.New()
	)
	logger.SetFormatter(&log.JSONFormatter{})

	return *logger, nil
}
