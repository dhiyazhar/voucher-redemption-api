package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/dhiyazhar/voucher-redemption-api/internal/entity"
	"github.com/dhiyazhar/voucher-redemption-api/internal/model"
	"github.com/dhiyazhar/voucher-redemption-api/internal/model/converter"
	"github.com/dhiyazhar/voucher-redemption-api/internal/repository"
)

type VoucherUsecase struct {
	DB                   *sql.DB
	VoucherRepository    repository.VoucherRepository
	RedemptionRepository repository.RedemptionRepository
	Logger               *slog.Logger
}

func NewVoucherUsecase(db *sql.DB,
	voucherRepository repository.VoucherRepository,
	redemptionRepository repository.RedemptionRepository,
	logger *slog.Logger,
) *VoucherUsecase {
	return &VoucherUsecase{
		DB:                   db,
		VoucherRepository:    voucherRepository,
		RedemptionRepository: redemptionRepository,
		Logger:               logger,
	}
}

func (c *VoucherUsecase) CreateVoucher(ctx context.Context, request *model.CreateVoucherRequest) (*model.VoucherResponse, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin a new transaction: %w", err)
	}
	defer tx.Rollback()

	voucher := &entity.Voucher{
		Code:       request.Code,
		Quota:      request.Quota,
		ValidUntil: time.Now().AddDate(0, 0, request.ValidUntilDays),
	}

	if err := c.VoucherRepository.Create(ctx, tx, voucher); err != nil {
		return nil, fmt.Errorf("create voucher failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("create voucher failed: %w", err)
	}

	return converter.VoucherToResponse(voucher), nil
}

func (c *VoucherUsecase) GetByCode(ctx context.Context, code string) (*model.VoucherResponse, error) {
	voucher, err := c.VoucherRepository.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get voucher failed: %w", err)
	}

	return converter.VoucherToResponse(voucher), nil
}

func (c *VoucherUsecase) ClaimVoucher(ctx context.Context, code string, userId int64) (*model.VoucherResponse, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin a new transaction claim voucher: %w", err)
	}
	defer tx.Rollback()

	voucher, err := c.VoucherRepository.GetByCodeForUpdate(ctx, tx, code)
	if err != nil {
		return nil, fmt.Errorf("claim voucher failed: %w", err)
	}

	if time.Now().After(voucher.ValidUntil) {
		return nil, model.ErrExpiredVoucher
	}

	if voucher.Quota <= 0 {
		return nil, model.ErrQuotaExhausted
	}

	newQuota, err := c.VoucherRepository.UpdateQuota(ctx, tx, voucher.ID, -1)
	if err != nil {
		return nil, err
	}

	redemption := &entity.Redemption{
		VoucherID:  voucher.ID,
		UserID:     userId,
		RedeemedAt: time.Now(),
	}

	err = c.RedemptionRepository.Create(ctx, tx, redemption)
	if err != nil {
		return nil, model.ErrAlreadyRedeemed
	}
	tx.Commit()

	voucher.Quota = newQuota

	return converter.VoucherToResponse(voucher), nil
}
