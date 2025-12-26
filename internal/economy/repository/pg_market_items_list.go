package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MarketItemOrder struct {
	OrderID   uuid.UUID `json:"orderId"`
	KingdomID string    `json:"kingdomId"`
	SellerID  string    `json:"sellerId"`
	ItemID    uuid.UUID `json:"itemId"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`

	// snapshot from item
	Tier       int `json:"tier"`
	Durability int `json:"durability"`
	MaxDur     int `json:"maxDurability"`
	BonusPct   int `json:"bonusPct"`
}

func (s *PgStore) ListItemOrders(ctx context.Context, kingdomID string) ([]MarketItemOrder, error) {
	rows, err := s.db.Query(ctx, `
		SELECT o.id, o.kingdom_id, o.seller_id, o.item_id, o.price, o.created_at,
		       i.tier, i.durability, i.max_durability, i.bonus_pct
		FROM market_item_orders o
		JOIN user_items i ON i.id = o.item_id
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
			&m.OrderID, &m.KingdomID, &m.SellerID, &m.ItemID, &m.Price, &m.CreatedAt,
			&m.Tier, &m.Durability, &m.MaxDur, &m.BonusPct,
		); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
