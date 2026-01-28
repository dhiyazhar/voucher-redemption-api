package entity

import (
	"time"
)

type Redemption struct {
	ID         int64     `json:"id"`
	VoucherID  int64     `json:"voucher_id"`
	UserID     int64     `json:"user_id"`
	RedeemedAt time.Time `json:"redeemed_at"`
}
