package controller

import (
	"jadwalin-auth-events-integrator/internal/service"
	"jadwalin-auth-events-integrator/internal/shared/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type Event struct {
	dig.In
	Service service.Holder
	Logger  log.Logger
}

func (impl *Event) handleSync(c echo.Context) error {
	userID := c.Param("userID")
	if userID == "" {
		impl.Logger.Error("User ID is empty")
		return c.JSON(http.StatusBadRequest, dto.Response{
			Status: http.StatusBadRequest,
			Error:  "User ID cannot be empty",
		})
	}

	token, err := impl.Service.Auth.GetToken(userID)
	if err != nil {
		impl.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, dto.Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	impl.Logger.Info("Syncing Events for User ID: ", userID)

	if err := impl.Service.Event.SyncAPIWithDB(token, userID); err != nil {
		impl.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, dto.Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return nil
}

func (impl *Event) handleGetEventsInHour(c echo.Context) error {
	hourString := c.Param("hour")
	if hourString == "" {
		impl.Logger.Error("Hour is empty")
		return c.JSON(http.StatusBadRequest, dto.Response{
			Status: http.StatusBadRequest,
			Error:  "Hour cannot be empty",
		})
	}

	hour, err := strconv.Atoi(hourString)
	if err != nil {
		impl.Logger.Error(err)
		return c.JSON(http.StatusBadRequest, dto.Response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		})
	}

	impl.Logger.Info("Getting Events in Hour: ", hour)

	events, err := impl.Service.Event.GetEventsInHour(hour)
	if err != nil {
		impl.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, dto.Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Status: http.StatusOK,
		Data:   events,
	})
}

func (impl *Event) handleGetUserEvents(c echo.Context) error {
	userID := c.Param("userID")
	if userID == "" {
		impl.Logger.Error("User ID is empty")
		return c.JSON(http.StatusBadRequest, dto.Response{
			Status: http.StatusBadRequest,
			Error:  "User ID cannot be empty",
		})
	}

	impl.Logger.Info("Getting Events for User ID: ", userID)

	//TODO: get user events from DB(Calling services)

	return nil
}
