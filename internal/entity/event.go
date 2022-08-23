package entity

type UserEvents struct {
	ID     string  `json:"user_id" bson:"_id"`
	Events []Event `json:"events" bson:"events"`
}

type Event struct {
	Description string    `json:"description" bson:"description"`
	Organizer   string    `json:"organizer" bson:"organizer"`
	Summary     string    `json:"summary" bson:"summary"`
	UpdatedAt   string    `json:"updated_at" bson:"updated_at"`
	StartTime   EventTime `json:"start_time" bson:"start_time"`
	EndTime     EventTime `json:"end_time" bson:"end_time"`
	URI         string    `json:"uri" bson:"uri"`
}

type EventTime struct {
	DateTime string `json:"date_time" bson:"date_time"`
	TimeZone string `json:"time_zone" bson:"time_zone"`
}
