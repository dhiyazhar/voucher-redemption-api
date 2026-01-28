package model

type WebResponse[T any] struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   T      `json:"data,omitempty"`
	Errors string `json:"errors,omitempty"`
}
