package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func (s *PgStore) BuyItem(
	ctx context.Context,
	orderID uuid.UUID,
	buyerID string,
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
		WHERE id=$1
		FOR UPDATE
	`, orderID).Scan(&itemID)
	if err != nil {
		return err
	}

	// transfer ownership
	tag, err := tx.Exec(ctx, `
		UPDATE user_items
		SET user_id=$2, listed=false, equipped=false
		WHERE id=$1 AND listed=true
	`, itemID, buyerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("item already sold")
	}

	_, err = tx.Exec(ctx, `DELETE FROM market_item_orders WHERE id=$1`, orderID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
