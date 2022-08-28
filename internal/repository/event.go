package repository

import (
	"jadwalin-auth-events-integrator/internal/entity"
	"jadwalin-auth-events-integrator/internal/shared/dto"
	"sort"
	"strings"
	"time"

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
		GetEventsInHour(string, string) ([]entity.UserEvents, error)
		GetUserEvents(string) ([]entity.Event, error)
		GetUserSummary(dto.SummaryRequest) ([]dto.SummaryResponse, error)
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

func (repo *eventRepo) GetEventsInHour(timeNow, timeHour string) ([]entity.UserEvents, error) {
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
			repo.logger.Infof("User with user ID %s is not exist when geting user events from database: %s", userID, err)
			return nil, err
		}

		repo.logger.Errorf("Getting user events from database error: %s", err)
		return nil, err
	}

	var user_events entity.UserEvents
	if err := result.Decode(&user_events); err != nil {
		repo.logger.Errorf("Error decode query result to user_events: %s", err)
		return nil, err
	}

	return user_events.Events, nil
}

func (repo *eventRepo) GetUserSummary(request dto.SummaryRequest) ([]dto.SummaryResponse, error) {
	repo.logger.Info(request)
	var response []dto.SummaryResponse

	// Get user events
	userEvents, err := repo.GetUserEvents(request.UserId)
	if err != nil {
		repo.logger.Errorf("Failed to get user events for making user summary: %v", err)
		return nil, err
	}

	// Build map which value is slice of hour from user events
	userEventsInMapShape := make(map[string][]int)
	for _, event := range userEvents {
		// Build slice of hour to map for start time
		startTime, err := time.Parse(time.RFC3339, event.StartTime.DateTime)
		if err != nil {
			repo.logger.Errorf("Failed to parse startTime of event when getting summary: %v", err)
			continue
		}
		endTime, err := time.Parse(time.RFC3339, event.EndTime.DateTime)
		if err != nil {
			repo.logger.Errorf("Failed to parse endTime of event when getting summary: %v", err)
			continue
		}

		startTimeDate := strings.Split(startTime.String(), " ")[0]
		endTimeDate := strings.Split(endTime.String(), " ")[0]
		if startTimeDate == endTimeDate {
			startTimeHour := startTime.Hour()
			_, startTimeDateHoursExist := userEventsInMapShape[startTimeDate]
			if !startTimeDateHoursExist {
				userEventsInMapShape[startTimeDate] = make([]int, 0)
			}
			if endTime.Minute() == 0 {
				for hour := startTimeHour; hour < endTime.Hour(); hour++ {
					userEventsInMapShape[startTimeDate] = append(userEventsInMapShape[startTimeDate], hour)
				}
			} else {
				for hour := startTimeHour; hour <= endTime.Hour(); hour++ {
					userEventsInMapShape[startTimeDate] = append(userEventsInMapShape[startTimeDate], hour)
				}
			}
		}
	}
	repo.logger.Info("Succesfully build map which value is slice of hour from user events: %s", userEventsInMapShape)

	// Getting current time in Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		repo.logger.Errorf("Failed to get time location when getting user summary from db: %v", err)
		return nil, err
	}
	currentTime := time.Now().In(loc)
	repo.logger.Info("CURRENT TIME: %s", currentTime)

	// Iterate for every days in request to create SummaryResponse
	for i := 0; i <= request.Days; i++ {

		// Create slice of hour from request. Range value is from 0 - 24.
		currentAvailability := make([]int, 0)
		var nowStartHour int
		if i == 0 {
			if currentTime.Hour() > request.StartHour {
				if currentTime.Minute() != 0 && currentTime.Second() != 0 {
					nowStartHour = currentTime.Hour() + 1
				} else {
					nowStartHour = currentTime.Hour()
				}
			} else {
				nowStartHour = request.StartHour
			}
		} else {
			nowStartHour = request.StartHour
		}

		for k := nowStartHour; k <= request.EndHour; k++ {
			currentAvailability = append(currentAvailability, k)
		}
		sort.Sort(sort.IntSlice(currentAvailability))
		repo.logger.Info("Succesfully create slice of availability hour from request: %s", currentAvailability)

		currentRequestTime := currentTime.AddDate(0, 0, i)
		currentRequestDate := strings.Split(currentRequestTime.String(), " ")[0]
		var availabilityResult []dto.TimeSpan

		currentEventHours, currentEventHoursIsExist := userEventsInMapShape[currentRequestDate]
		if currentEventHoursIsExist {
			sort.Sort(sort.IntSlice(currentEventHours))
			substractResult := substract(currentAvailability, currentEventHours)
			repo.logger.Info("Succesfully substract currentAvailability slice with currentEventHours: %s", substractResult)

			substractResultLength := len(substractResult)
			if substractResultLength != 0 {
				var startAvailabilityBoundary int
				continuityTimeFlag := false
				for j := 0; j < substractResultLength-1; j++ {
					if substractResult[j+1] == (substractResult[j] + 1) {
						if continuityTimeFlag {
							continue
						} else {
							continuityTimeFlag = true
							startAvailabilityBoundary = substractResult[j]
						}
					} else {
						if continuityTimeFlag {
							continuityTimeFlag = false
							availabilityResult = append(availabilityResult, dto.TimeSpan{
								StartHour: startAvailabilityBoundary,
								EndHour:   substractResult[j] + 1,
							})
						} else {
							availabilityResult = append(availabilityResult, dto.TimeSpan{
								StartHour: substractResult[j],
								EndHour:   substractResult[j] + 1,
							})
						}
					}
				}

				substractResultLastElement := substractResult[substractResultLength-1]
				if continuityTimeFlag {
					if substractResultLastElement == request.EndHour {
						availabilityResult = append(availabilityResult, dto.TimeSpan{
							StartHour: startAvailabilityBoundary,
							EndHour:   substractResultLastElement,
						})
					} else {
						availabilityResult = append(availabilityResult, dto.TimeSpan{
							StartHour: startAvailabilityBoundary,
							EndHour:   substractResultLastElement + 1,
						})
					}
				} else {
					if substractResultLastElement != request.EndHour {
						availabilityResult = append(availabilityResult, dto.TimeSpan{
							StartHour: substractResultLastElement,
							EndHour:   substractResultLastElement + 1,
						})
					}
				}
			}
		} else {
			if len(currentAvailability) != 1 {
				availabilityResult = append(availabilityResult, dto.TimeSpan{
					StartHour: currentAvailability[0],
					EndHour:   currentAvailability[len(currentAvailability)-1],
				})
			}
		}

		if len(availabilityResult) != 0 {
			response = append(response, dto.SummaryResponse{
				Date:         currentRequestDate,
				Availability: availabilityResult,
			})
		}
		repo.logger.Info("SummaryResponse added: %s", response)
	}

	return response, nil
}

// Return firstSlice without element that also exist in secondSlice
func substract(firstSlice []int, secondSlice []int) []int {
	var result []int
	for _, s1 := range firstSlice {
		found := false
		for _, s2 := range secondSlice {
			if s1 == s2 {
				found = true
				break
			}
		}

		if !found {
			result = append(result, s1)
		}
	}

	return result
}

func NewEvent(logger log.Logger, db *mongo.Database, redis *redis.Client) (Event, error) {
	return &eventRepo{
		logger: logger,
		db:     db,
		redis:  redis}, nil
}
