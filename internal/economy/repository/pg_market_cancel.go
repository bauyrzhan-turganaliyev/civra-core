package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s *PgStore) CancelSellOrder(
	ctx context.Context,
	orderID uuid.UUID,
	sellerID string,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1) lock order
	var resource string
	var quantity int
	err = tx.QueryRow(ctx,
		`SELECT resource, quantity FROM market_orders
		 WHERE id=$1 AND seller_id=$2
		 FOR UPDATE`,
		orderID, sellerID,
	).Scan(&resource, &quantity)

	if err != nil {
		return err // not found or not owner
	}

	// 2) return resource to seller
	_, err = tx.Exec(ctx,
		`INSERT INTO personal_inventory (user_id, resource, quantity)
		 VALUES ($1,$2,$3)
		 ON CONFLICT (user_id, resource)
		 DO UPDATE SET quantity = personal_inventory.quantity + EXCLUDED.quantity`,
		sellerID, resource, quantity,
	)
	if err != nil {
		return err
	}

	// 3) delete order
	_, err = tx.Exec(ctx,
		`DELETE FROM market_orders WHERE id=$1`,
		orderID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
