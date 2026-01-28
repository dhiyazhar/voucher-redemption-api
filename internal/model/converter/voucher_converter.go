package converter

import (
	"github.com/dhiyazhar/voucher-redemption-api/internal/entity"
	"github.com/dhiyazhar/voucher-redemption-api/internal/model"
)

func VoucherToResponse(v *entity.Voucher) *model.VoucherResponse {
	if v == nil {
		return &model.VoucherResponse{}
	}

	return &model.VoucherResponse{
		ID:         v.ID,
		Code:       v.Code,
		Quota:      v.Quota,
		ValidUntil: v.ValidUntil,
	}
}
