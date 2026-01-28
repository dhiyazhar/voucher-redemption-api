package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dhiyazhar/voucher-redemption-api/internal/entity"
	"github.com/dhiyazhar/voucher-redemption-api/internal/model"
	"github.com/lib/pq"
)

type redemptionRepository struct {
	db *sql.DB
}

func NewRedemptionRepository(db *sql.DB) RedemptionRepository {
	return &redemptionRepository{
		db: db,
	}
}

func (r *redemptionRepository) Create(ctx context.Context, tx *sql.Tx, rh *entity.Redemption) error {
	query := "INSERT INTO redemption_history(voucher_id, user_id, redeemed_at) VALUES($1, $2, $3)"

	_, err := tx.ExecContext(ctx, query, rh.VoucherID, rh.UserID, rh.RedeemedAt)
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

func (r *redemptionRepository) GetByVoucherCode(ctx context.Context, code string) ([]entity.Redemption, error) {
	query := `
        SELECT 
            rh.id, 
            rh.voucher_id, 
            rh.user_id, 
            rh.redeemed_at
        FROM redemption_history rh
        INNER JOIN vouchers v ON v.id = rh.voucher_id
        WHERE v.code = $1
        ORDER BY rh.redeemed_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redemptions []entity.Redemption
	for rows.Next() {
		var redemption entity.Redemption
		err := rows.Scan(
			&redemption.ID,
			&redemption.VoucherID,
			&redemption.UserID,
			&redemption.RedeemedAt,
		)
		if err != nil {
			return nil, err
		}
		redemptions = append(redemptions, redemption)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return redemptions, nil
}
