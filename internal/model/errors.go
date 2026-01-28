package model

import (
	"errors"
)

var (
	ErrInternalServerError = errors.New("Internal server error")

	//Voucher error
	ErrNotFound       = errors.New("Requested item not found")
	ErrConflict       = errors.New("Item already exists")
	ErrBadParamInput  = errors.New("Given param is not valid")
	ErrExpiredVoucher = errors.New("Voucher is expired")
	ErrQuotaExhausted = errors.New("Voucher quota insufficient")

	//Redemption error
	ErrAlreadyRedeemed = errors.New("voucher already redeemed")
)
