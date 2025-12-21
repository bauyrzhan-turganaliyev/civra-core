package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserItem struct {
	ID         uuid.UUID `json:"id"`
	UserID     string    `json:"userId"`
	ItemType   string    `json:"itemType"`
	Tier       int       `json:"tier"`
	Durability int       `json:"durability"`
	MaxDur     int       `json:"maxDurability"`
	BonusPct   int       `json:"bonusPct"`
	Equipped   bool      `json:"equipped"`
	Listed     bool      `json:"listed"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (s *PgStore) ListUserItems(ctx context.Context, userID string) ([]UserItem, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, user_id, item_type, tier, durability, max_durability, bonus_pct, equipped, listed, created_at
		FROM user_items
		WHERE user_id=$1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []UserItem
	for rows.Next() {
		var it UserItem
		if err := rows.Scan(
			&it.ID, &it.UserID, &it.ItemType, &it.Tier,
			&it.Durability, &it.MaxDur, &it.BonusPct,
			&it.Equipped, &it.Listed, &it.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}
