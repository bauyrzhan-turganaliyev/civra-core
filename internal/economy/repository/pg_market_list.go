package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MarketOrder struct {
	ID        uuid.UUID `json:"id"`
	KingdomID string    `json:"kingdomId"`
	SellerID  string    `json:"sellerId"`
	Resource  string    `json:"resource"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *PgStore) ListMarketOrders(ctx context.Context, kingdomID string, limit int) ([]MarketOrder, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, kingdom_id, seller_id, resource, quantity, price, created_at
		 FROM market_orders
		 WHERE kingdom_id=$1
		 ORDER BY created_at DESC
		 LIMIT $2`,
		kingdomID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]MarketOrder, 0, limit)
	for rows.Next() {
		var o MarketOrder
		if err := rows.Scan(&o.ID, &o.KingdomID, &o.SellerID, &o.Resource, &o.Quantity, &o.Price, &o.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
