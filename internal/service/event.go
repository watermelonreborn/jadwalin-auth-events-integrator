package service

import (
	"context"
	"jadwalin-auth-events-integrator/internal/entity"
	"jadwalin-auth-events-integrator/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type (
	Event interface {
		SyncAPIWithDB(*oauth2.Token, string) error
		GetEventsInHour(int) ([]entity.UserEvents, error)
		GetUserEvents(string) ([]entity.Event, error)
	}

	eventService struct {
		logger     log.Logger
		repository repository.Holder
	}
)

func (service *eventService) SyncAPIWithDB(token *oauth2.Token, userID string) error {
	client := config.Client(context.Background(), token)

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		service.logger.Errorf("Unable to retrieve Calendar client: %v", err)
		return err
	}

	t := time.Now().Format(time.RFC3339)
	tMonth := time.Now().AddDate(0, 1, 0).Format(time.RFC3339)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).TimeMax(tMonth).OrderBy("startTime").Do()
	if err != nil {
		service.logger.Errorf("Unable to retrieve next ten of the user's events: %v", err)
		return err
	}

	userEvents := entity.UserEvents{
		ID:     userID,
		Events: make([]entity.Event, 0),
	}

	for _, item := range events.Items {
		var conferenceName, conferenceLink string
		if item.ConferenceData != nil {
			conferenceName = item.ConferenceData.ConferenceSolution.Name
			conferenceLink = item.ConferenceData.EntryPoints[0].Uri
		}

		userEvents.Events = append(userEvents.Events, entity.Event{
			Description: item.Description,
			Organizer:   item.Organizer.Email,
			Summary:     item.Summary,
			UpdatedAt:   item.Updated,
			StartTime: entity.EventTime{
				DateTime: item.Start.DateTime,
				TimeZone: item.Start.TimeZone,
			},
			EndTime: entity.EventTime{
				DateTime: item.End.DateTime,
				TimeZone: item.End.TimeZone,
			},
			ConferenceData: entity.ConferenceData{
				Name: conferenceName,
				URI:  conferenceLink,
			},
		})
	}

	if err := service.repository.Event.UpsertEvents(userEvents); err != nil {
		service.logger.Errorf("Unable to upsert events to database: %v", err)
		return err
	}

	service.logger.Info("Sync Events API with DB success")

	return nil
}

func (service *eventService) GetEventsInHour(hour int) ([]entity.UserEvents, error) {
	timeNow := time.Now().Format(time.RFC3339)
	timeHour := time.Now().Add(time.Duration(hour) * time.Hour).Format(time.RFC3339)
	events, err := service.repository.Event.GetAllEvents(timeNow, timeHour)
	if err != nil {
		service.logger.Errorf("Unable to get events from database: %v", err)
		return nil, err
	}

	service.logger.Infof("Get Events In %d Hour success", hour)

	return events, nil
}

func (service *eventService) GetUserEvents(userId string) ([]entity.Event, error) {
	events, err := service.repository.Event.GetUserEvents(userId)
	if err != nil {
		service.logger.Errorf("Unable to get user events with user ID %v from database: %v", userId, err)
		return nil, err
	}

	service.logger.Infof("Get user events with user ID %v", userId)

	return events, nil
}

func NewEvent(logger log.Logger, repo repository.Holder) (Event, error) {
	return &eventService{
		logger:     logger,
		repository: repo,
	}, nil
}
