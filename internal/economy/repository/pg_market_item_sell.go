package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func (s *PgStore) SellItem(
	ctx context.Context,
	kingdomID, sellerID string,
	itemID uuid.UUID,
	price int,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// lock item
	tag, err := tx.Exec(ctx, `
		UPDATE user_items
		SET listed=true
		WHERE id=$1 AND user_id=$2 AND equipped=false AND listed=false
	`, itemID, sellerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("item cannot be sold")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO market_item_orders (id, kingdom_id, seller_id, item_id, price)
		VALUES ($1,$2,$3,$4,$5)
	`, uuid.New(), kingdomID, sellerID, itemID, price)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
