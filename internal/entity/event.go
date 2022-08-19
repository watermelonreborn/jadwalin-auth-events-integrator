package entity

type Event struct {
}

type UserEvents struct {
	Id     string  `bson: "_id"`
	Events []Event `bson:"events"`
}
