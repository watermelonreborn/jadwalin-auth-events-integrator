package dto

type Response struct {
	Status int         `json:"status"`
	Error  interface{} `json:"error"`
	Data   interface{} `json:"data"`
}
