package repository

import (
	"context"
	"database/sql"

	"github.com/dhiyazhar/voucher-redemption-api/internal/entity"
)

type (
	VoucherRepository interface {
		Create(ctx context.Context, tx *sql.Tx, voucher *entity.Voucher) error
		GetByCode(ctx context.Context, code string) (*entity.Voucher, error)
		GetByCodeForUpdate(ctx context.Context, tx *sql.Tx, code string) (*entity.Voucher, error)
		UpdateQuota(ctx context.Context, tx *sql.Tx, voucherId int64, delta int) (int64, error)
	}

	RedemptionRepository interface {
		Create(ctx context.Context, tx *sql.Tx, redemption *entity.Redemption) error
		GetByVoucherCode(ctx context.Context, code string) ([]entity.Redemption, error)
	}
)
