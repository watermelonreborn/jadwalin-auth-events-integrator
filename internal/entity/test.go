package entity

type Test struct {
	Name string `bson:"name"`
}

type Token struct {
	AccessToken  string `bson:"access_token"`
	TokenType    string `bson:"token_type"`
	RefreshToken string `bson:"refresh_token"`
}
