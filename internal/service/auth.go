package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"jadwalin-auth-events-integrator/internal/entity"
	"jadwalin-auth-events-integrator/internal/repository"
	"jadwalin-auth-events-integrator/internal/shared/dto"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"

	log "github.com/sirupsen/logrus"
)

const (
	SESSION_ID         = "jadwalin-session-id"
	OAUTH_STATE_STRING = "state-token"
)

var (
	config *oauth2.Config
)

type (
	Auth interface {
		URL() string
		GenerateToken(string, string) (dto.TokenResponse, error)
		GetToken(string) (*oauth2.Token, error)
		GetUserInfo(string) (dto.UserInfoResponse, error)
		UpsertUser(entity.User) error
	}

	authService struct {
		logger     log.Logger
		repository repository.Holder
	}
)

func init() {
	bytes, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		panic(err)
	}

	config, err = google.ConfigFromJSON(bytes, "https://www.googleapis.com/auth/userinfo.email", calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		panic(err)
	}
}

func (service *authService) URL() string {
	return config.AuthCodeURL(OAUTH_STATE_STRING, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (service *authService) GenerateToken(state string, code string) (dto.TokenResponse, error) {
	if state != OAUTH_STATE_STRING {
		return dto.TokenResponse{}, fmt.Errorf("invalid oauth state")
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return dto.TokenResponse{}, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	service.logger.Info(token.AccessToken)
	service.logger.Info(token.RefreshToken)

	return dto.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.Unix(),
	}, nil
}

func (service *authService) GetToken(userID string) (*oauth2.Token, error) {
	refreshToken, err := service.repository.Auth.GetToken(userID)
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{RefreshToken: refreshToken}, nil
}

func (service *authService) GetUserInfo(accessToken string) (dto.UserInfoResponse, error) {
	var (
		userInfoResponse dto.UserInfoResponse
	)

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return dto.UserInfoResponse{}, fmt.Errorf("failed to get user info: %s", err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return dto.UserInfoResponse{}, fmt.Errorf("failed to get user info: %s", response.Status)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return dto.UserInfoResponse{}, fmt.Errorf("failed to read user info: %s", err.Error())
	}

	service.logger.Info(string(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &userInfoResponse); err != nil {
		return dto.UserInfoResponse{}, fmt.Errorf("failed to unmarshal user info: %s", err.Error())
	}

	return userInfoResponse, nil

}

func (service *authService) UpsertUser(user entity.User) error {
	err := service.repository.Auth.UpsertUser(user)
	if err != nil {
		service.logger.Error("Error: %s", err)
		return err
	}

	return nil
}

func NewAuth(logger log.Logger, repo repository.Holder) (Auth, error) {
	return &authService{
		logger:     logger,
		repository: repo,
	}, nil
}
