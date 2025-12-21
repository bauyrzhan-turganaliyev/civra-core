package repository

import (
	"context"

	"github.com/google/uuid"
)

type MarketItemOrder struct {
	OrderID  uuid.UUID `json:"orderId"`
	ItemID   uuid.UUID `json:"itemId"`
	SellerID string    `json:"sellerId"`
	Price    int       `json:"price"`
	Tier     int       `json:"tier"`
	Dur      int       `json:"durability"`
	Bonus    int       `json:"bonusPct"`
}

func (s *PgStore) ListItemOrders(
	ctx context.Context,
	kingdomID string,
) ([]MarketItemOrder, error) {

	rows, err := s.db.Query(ctx, `
		SELECT o.id, i.id, o.seller_id, o.price,
		       i.tier, i.durability, i.bonus_pct
		FROM market_item_orders o
		JOIN user_items i ON i.id=o.item_id
		WHERE o.kingdom_id=$1
		ORDER BY o.created_at DESC
	`, kingdomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MarketItemOrder
	for rows.Next() {
		var m MarketItemOrder
		if err := rows.Scan(
			&m.OrderID, &m.ItemID, &m.SellerID, &m.Price,
			&m.Tier, &m.Dur, &m.Bonus,
		); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
