package entity

type User struct {
	ID           string `bson:"_id"`
	RefreshToken string `bson:"refresh_token"`
}
