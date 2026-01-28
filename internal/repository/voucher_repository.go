package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dhiyazhar/voucher-redemption-api/internal/entity"
	"github.com/dhiyazhar/voucher-redemption-api/internal/model"
	"github.com/lib/pq"
)

type voucherRepository struct {
	db *sql.DB
}

func NewVoucherRepository(db *sql.DB) VoucherRepository {
	return &voucherRepository{
		db: db,
	}
}

func (r *voucherRepository) Create(ctx context.Context, tx *sql.Tx, v *entity.Voucher) error {
	query := "INSERT INTO vouchers(code, quota, valid_until) VALUES($1, $2, $3) RETURNING id"

	err := tx.QueryRowContext(ctx, query, v.Code, v.Quota, v.ValidUntil).Scan(&v.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return model.ErrConflict
			}
		}

		return err
	}

	return nil
}

func (r *voucherRepository) GetByCode(ctx context.Context, code string) (*entity.Voucher, error) {
	query := "SELECT id, code, quota, valid_until FROM vouchers WHERE code = $1"

	var v entity.Voucher
	err := r.db.QueryRowContext(ctx, query, code).Scan(&v.ID, &v.Code, &v.Quota, &v.ValidUntil)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}

		return nil, err
	}

	return &v, nil
}

func (r *voucherRepository) GetByCodeForUpdate(ctx context.Context, tx *sql.Tx, code string) (*entity.Voucher, error) {
	query := "SELECT id, code, quota, valid_until FROM vouchers WHERE code = $1 FOR UPDATE"

	var v entity.Voucher
	err := tx.QueryRowContext(ctx, query, code).Scan(&v.ID, &v.Code, &v.Quota, &v.ValidUntil)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}

		return nil, err
	}

	return &v, nil
}

func (r *voucherRepository) UpdateQuota(ctx context.Context, tx *sql.Tx, voucherId int64, delta int) (int64, error) {
	query := "UPDATE vouchers SET quota = quota + $2 WHERE id = $1 AND quota + $2 >= 0 RETURNING quota"

	var quota int64
	err := tx.QueryRowContext(ctx, query, voucherId, delta).Scan(&quota)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, model.ErrNotFound
		}
		return 0, fmt.Errorf("failed to update quota: %w", err)
	}

	return quota, nil
}
