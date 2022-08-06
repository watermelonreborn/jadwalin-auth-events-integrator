package dto

type TestRequest struct {
	Name string `json:"name"`
}

type TestResponse struct {
	Message string `json:"message"`
}
