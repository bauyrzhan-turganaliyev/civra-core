package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNotEnoughMaterials = errors.New("not enough materials")

func (s *PgStore) CraftTool(ctx context.Context, userID string, tier, bonusPct, maxDur, ironCost, woodCost int) (uuid.UUID, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	getQty := func(res string) (int, error) {
		var q int
		err := tx.QueryRow(ctx, `
			SELECT quantity FROM personal_inventory
			WHERE user_id=$1 AND resource=$2
			FOR UPDATE
		`, userID, res).Scan(&q)
		if err != nil {
			return 0, nil
		}
		return q, nil
	}

	ironQty, err := getQty("iron")
	if err != nil {
		return uuid.Nil, err
	}
	woodQty, err := getQty("wood")
	if err != nil {
		return uuid.Nil, err
	}

	if ironQty < ironCost || woodQty < woodCost {
		return uuid.Nil, ErrNotEnoughMaterials
	}

	_, err = tx.Exec(ctx, `
		UPDATE personal_inventory SET quantity = quantity - $3
		WHERE user_id=$1 AND resource=$2
	`, userID, "iron", ironCost)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE personal_inventory SET quantity = quantity - $3
		WHERE user_id=$1 AND resource=$2
	`, userID, "wood", woodCost)
	if err != nil {
		return uuid.Nil, err
	}

	id := uuid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO user_items (id, user_id, item_type, tier, durability, max_durability, bonus_pct)
		VALUES ($1,$2,'tool',$3,$4,$5,$6)
	`, id, userID, tier, maxDur, maxDur, bonusPct)
	if err != nil {
		return uuid.Nil, err
	}

	return id, tx.Commit(ctx)
}
