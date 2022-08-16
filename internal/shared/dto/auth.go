package dto

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry"`
}

type UserInfoResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
