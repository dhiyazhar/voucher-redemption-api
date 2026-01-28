package usecase

import (
	"context"
	"log/slog"

	"github.com/dhiyazhar/voucher-redemption-api/internal/model"
	"github.com/dhiyazhar/voucher-redemption-api/internal/model/converter"
	"github.com/dhiyazhar/voucher-redemption-api/internal/repository"
)

type RedemptionUsecase struct {
	RedemptionRepository repository.RedemptionRepository
	VoucherRepository    repository.VoucherRepository
	Logger               *slog.Logger
}

func NewRedemptionUsecase(
	redemptionRepository repository.RedemptionRepository,
	voucherRepository repository.VoucherRepository,
	logger *slog.Logger,
) *RedemptionUsecase {
	return &RedemptionUsecase{
		RedemptionRepository: redemptionRepository,
		VoucherRepository:    voucherRepository,
		Logger:               logger,
	}
}

func (u *RedemptionUsecase) GetByVoucherCode(ctx context.Context, code string) ([]model.RedemptionResponse, error) {
	_, err := u.VoucherRepository.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	redemptions, err := u.RedemptionRepository.GetByVoucherCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return converter.RedemptionsToResponse(redemptions), nil
}
