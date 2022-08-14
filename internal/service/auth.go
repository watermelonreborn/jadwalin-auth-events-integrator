package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jadwalin-auth-events-integrator/internal/repository"
	"jadwalin-auth-events-integrator/internal/shared/dto"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const (
	SESSION_ID         = "jadwalin-session-id"
	OAUTH_STATE_STRING = "state-token"
)

var (
	store  *sessions.CookieStore
	config *oauth2.Config
)

type (
	Auth interface {
		URL() string
		GenerateToken(echo.Context, string, string) error
		GetToken(echo.Context) (dto.TokenResponse, error)
		GetUserInfo(string) (dto.UserInfoResponse, error)
	}

	authService struct {
		logger     log.Logger
		repository repository.Holder
	}
)

func init() {
	bytes, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		panic(err)
	}

	config, err = google.ConfigFromJSON(bytes, "https://www.googleapis.com/auth/userinfo.email", calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		panic(err)
	}

	store = newCookieStore()
}

func newCookieStore() *sessions.CookieStore {
	authKey := []byte("my-auth-key-very-secret")
	encryptKey := []byte("my-encryption-key-very-secret123")
	store := sessions.NewCookieStore(authKey, encryptKey)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	return store
}

func (service *authService) URL() string {
	return config.AuthCodeURL(OAUTH_STATE_STRING, oauth2.AccessTypeOffline)
}

func (service *authService) GenerateToken(c echo.Context, state string, code string) error {
	if state != OAUTH_STATE_STRING {
		return fmt.Errorf("invalid oauth state")
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("code exchange failed: %s", err.Error())
	}

	service.logger.Info(token.AccessToken)
	service.logger.Info(token.RefreshToken)

	session, err := store.Get(c.Request(), SESSION_ID)
	if err != nil {
		return fmt.Errorf("failed to get session: %s", err.Error())
	}

	session.Values["token"] = token.AccessToken
	session.Values["refresh"] = token.RefreshToken
	session.Values["expiry"] = token.Expiry.Unix()

	if err := session.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("failed to save session: %s", err.Error())
	}

	return nil
}

func (service *authService) GetToken(c echo.Context) (dto.TokenResponse, error) {
	session, err := store.Get(c.Request(), SESSION_ID)
	if err != nil {
		return dto.TokenResponse{}, fmt.Errorf("failed to get session: %s", err.Error())
	}

	if len(session.Values) == 0 {
		return dto.TokenResponse{}, fmt.Errorf("no session")
	}

	token := session.Values["token"].(string)
	refresh := session.Values["refresh"].(string)
	expiry := session.Values["expiry"].(int64)

	return dto.TokenResponse{
		AccessToken:  token,
		RefreshToken: refresh,
		Expiry:       expiry,
	}, nil
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

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return dto.UserInfoResponse{}, fmt.Errorf("failed to read user info: %s", err.Error())
	}

	service.logger.Info(string(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &userInfoResponse); err != nil {
		return dto.UserInfoResponse{}, fmt.Errorf("failed to unmarshal user info: %s", err.Error())
	}

	return userInfoResponse, nil

}

func OAuth(logger log.Logger, repo repository.Holder) (Auth, error) {
	return &authService{
		logger:     logger,
		repository: repo,
	}, nil
}
