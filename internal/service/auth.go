package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jadwalin-auth-events-integrator/internal/repository"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var (
	config *oauth2.Config
)

type (
	Auth interface {
		IndexAuth(echo.Context) (string, error)
		ManageAPICode(echo.Context, string) error
	}

	authService struct {
		logger     log.Logger
		repository repository.Holder
	}
)

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (service *authService) IndexAuth(c echo.Context) (string, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		service.logger.Errorf("Unable to read client secret file: %v", err)
		return "", err
	}

	config, err = google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		service.logger.Errorf("Unable to parse client secret file to config: %v", err)
		return "", err
	}

	// tokFile := "token.json"
	// tok, err := tokenFromFile(tokFile)
	// service.logger.Info(authURL)

	// if err != nil {
	// 	service.logger.Info("dlm if sebelum redirect")
	// 	return authURL, nil
	// service.logger.Info("dlm if setelah redirect")
	// tok = getTokenFromWeb(config)
	// saveToken(tokFile, tok)
	// }

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL, nil
}

func (service *authService) ManageAPICode(c echo.Context, authCode string) error {
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	// ctx := context.Background()

	// client := config.Client(ctx, tok)

	// srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve Calendar client: %v", err)
	// }

	// t := time.Now().Format(time.RFC3339)
	// events, err := srv.Events.List("primary").ShowDeleted(false).
	// 	SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	// }
	// fmt.Println("Upcoming events:")
	// if len(events.Items) == 0 {
	// 	fmt.Println("No upcoming events found.")
	// } else {
	// 	for _, item := range events.Items {
	// 		date := item.Start.DateTime
	// 		if date == "" {
	// 			date = item.Start.Date
	// 		}
	// 		fmt.Printf("%v (%v)\n", item.Summary, date)
	// 	}
	// }

	return nil
}

func OAuth(logger log.Logger, repo repository.Holder) (Auth, error) {
	return &authService{
		logger:     logger,
		repository: repo,
	}, nil
}
