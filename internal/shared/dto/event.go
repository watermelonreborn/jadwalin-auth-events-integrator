package dto

type SummaryRequest struct {
	UserId    string `json:"user_id"`
	Days      int    `json:"days"`
	StartHour int    `json:"start_hour"`
	EndHour   int    `json:"end_hour"`
}

type SummaryResponse struct {
	Date         string `json:"date"`
	Availibility []TimeSpan
}

type TimeSpan struct {
	StartHour int `json:"start_hour"`
	EndHour   int `json:"end_hour"`
}
