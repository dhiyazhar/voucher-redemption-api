package converter

import (
	"github.com/dhiyazhar/voucher-redemption-api/internal/entity"
	"github.com/dhiyazhar/voucher-redemption-api/internal/model"
)

func RedemptionsToResponse(redemptions []entity.Redemption) []model.RedemptionResponse {
	if len(redemptions) == 0 {
		return []model.RedemptionResponse{}
	}

	responses := make([]model.RedemptionResponse, len(redemptions))
	for i, r := range redemptions {
		responses[i] = model.RedemptionResponse{
			ID:         r.ID,
			VoucherID:  r.VoucherID,
			UserID:     r.UserID,
			RedeemedAt: r.RedeemedAt,
		}
	}

	return responses
}
