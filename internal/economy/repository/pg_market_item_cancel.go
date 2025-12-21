package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s *PgStore) CancelItemSale(
	ctx context.Context,
	orderID uuid.UUID,
	sellerID string,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var itemID uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT item_id
		FROM market_item_orders
		WHERE id=$1 AND seller_id=$2
		FOR UPDATE
	`, orderID, sellerID).Scan(&itemID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM market_item_orders WHERE id=$1`, orderID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE user_items SET listed=false WHERE id=$1
	`, itemID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
