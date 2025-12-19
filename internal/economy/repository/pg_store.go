package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgStore struct {
	db *pgxpool.Pool
}

func NewPgStore(db *pgxpool.Pool) *PgStore {
	return &PgStore{db: db}
}

func (s *PgStore) Gather(
	ctx context.Context,
	userID, kingdomID, resource string,
	quotaRequired int,
	now time.Time,
	amount int,
) (toKingdom, toPersonal, progress int, err error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// ✅ лучше передавать date как time.Time (Postgres сам возьмёт дату)
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// 0) ensure quota row exists (race-safe)
	_, err = tx.Exec(ctx,
		`INSERT INTO quota_progress (user_id, day, resource, progress)
     VALUES ($1,$2,$3,0)
     ON CONFLICT (user_id, day, resource) DO NOTHING`,
		userID, day, resource,
	)
	if err != nil {
		return
	}

	// 1) now lock and read
	err = tx.QueryRow(ctx,
		`SELECT progress FROM quota_progress
     WHERE user_id=$1 AND day=$2 AND resource=$3
     FOR UPDATE`,
		userID, day, resource,
	).Scan(&progress)
	if err != nil {
		return
	}

	remaining := quotaRequired - progress
	if remaining < 0 {
		remaining = 0
	}

	if remaining > 0 {
		if amount <= remaining {
			toKingdom = amount
		} else {
			toKingdom = remaining
			toPersonal = amount - remaining
		}
	} else {
		toPersonal = amount
	}

	if toPersonal > 0 {
		_, err = tx.Exec(ctx,
			`INSERT INTO personal_inventory (user_id, resource, quantity)
		 VALUES ($1,$2,$3)
		 ON CONFLICT (user_id, resource)
		 DO UPDATE SET quantity = personal_inventory.quantity + EXCLUDED.quantity`,
			userID, resource, toPersonal,
		)
		if err != nil {
			return
		}
	}

	// 2) обновляем quota + kingdom inventory
	if toKingdom > 0 {
		progress += toKingdom

		_, err = tx.Exec(ctx,
			`UPDATE quota_progress
             SET progress=$4
             WHERE user_id=$1 AND day=$2 AND resource=$3`,
			userID, day, resource, progress,
		)
		if err != nil {
			return
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO kingdom_inventory (kingdom_id, resource, quantity)
             VALUES ($1,$2,$3)
             ON CONFLICT (kingdom_id, resource)
             DO UPDATE SET quantity = kingdom_inventory.quantity + EXCLUDED.quantity`,
			kingdomID, resource, toKingdom,
		)
		if err != nil {
			return
		}
	}

	err = tx.Commit(ctx)
	return
}

func (s *PgStore) BuyOrder(
	ctx context.Context,
	orderID uuid.UUID,
	buyerID string,
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
		 WHERE id=$1
		 FOR UPDATE`,
		orderID,
	).Scan(&resource, &quantity)
	if err != nil {
		return err
	}

	// 2) give resource to buyer
	_, err = tx.Exec(ctx,
		`INSERT INTO personal_inventory (user_id, resource, quantity)
		 VALUES ($1,$2,$3)
		 ON CONFLICT (user_id, resource)
		 DO UPDATE SET quantity = personal_inventory.quantity + EXCLUDED.quantity`,
		buyerID, resource, quantity,
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
