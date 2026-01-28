package entity

import (
	"time"
)

type Voucher struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	Quota      int64     `json:"quota"`
	ValidUntil time.Time `json:"valid_until"`
}
