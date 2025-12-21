package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrItemNotFoundOrListed = errors.New("item not found or listed")

func (s *PgStore) EquipItem(ctx context.Context, userID string, itemID uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `UPDATE user_items SET equipped=false WHERE user_id=$1 AND equipped=true`, userID)
	if err != nil {
		return err
	}

	tag, err := tx.Exec(ctx, `
		UPDATE user_items
		SET equipped=true
		WHERE id=$1 AND user_id=$2 AND listed=false
	`, itemID, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrItemNotFoundOrListed
	}

	return tx.Commit(ctx)
}
