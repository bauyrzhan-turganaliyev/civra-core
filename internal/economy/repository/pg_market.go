package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s *PgStore) CreateSellOrder(
	ctx context.Context,
	kingdomID, sellerID, resource string,
	quantity, price int,
) (uuid.UUID, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1) lock seller inventory
	var available int
	err = tx.QueryRow(ctx,
		`SELECT quantity FROM personal_inventory
		 WHERE user_id=$1 AND resource=$2
		 FOR UPDATE`,
		sellerID, resource,
	).Scan(&available)
	if err != nil {
		return uuid.Nil, err
	}

	if available < quantity {
		return uuid.Nil, ErrNotEnoughResource
	}

	// 2) subtract from seller
	_, err = tx.Exec(ctx,
		`UPDATE personal_inventory
		 SET quantity = quantity - $3
		 WHERE user_id=$1 AND resource=$2`,
		sellerID, resource, quantity,
	)
	if err != nil {
		return uuid.Nil, err
	}

	orderID := uuid.New()

	// 3) insert order
	_, err = tx.Exec(ctx,
		`INSERT INTO market_orders
		 (id, kingdom_id, seller_id, resource, quantity, price)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		orderID, kingdomID, sellerID, resource, quantity, price,
	)
	if err != nil {
		return uuid.Nil, err
	}

	err = tx.Commit(ctx)
	return orderID, err
}
