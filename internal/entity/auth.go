package entity

type User struct {
	Name  string `bson:"name"`
	Token string `bson:"token"`
}
