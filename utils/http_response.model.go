package utils

type HttpResponse struct {
	IsSuccessful bool   `json:"isSuccessful"`
	Status       int    `json:"status"`
	Code         string `json:"code"`
	Data         any    `json:"data"`
}
