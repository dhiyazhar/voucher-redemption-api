package model

import "time"

type VoucherResponse struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	Quota      int64     `json:"quota"`
	ValidUntil time.Time `json:"valid_until"`
}

type CreateVoucherRequest struct {
	Code           string `json:"code" validate:"required,max=100"`
	Quota          int64  `json:"quota" validate:"required,min=1"`
	ValidUntilDays int    `json:"valid_until_days" validate:"required,min=1"`
}
