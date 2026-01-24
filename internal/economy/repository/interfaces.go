package repository

import (
	"context"
	"time"
)

type GatherStore interface {
	Gather(
		ctx context.Context,
		userID, kingdomID, resource string,
		quotaRequired int,
		now time.Time,
		amount int,
	) (toKingdom, toPersonal, progress int, toolBonusPct int, toolUsed bool, err error)
}

type InventoryQueryStore interface {
	GetKingdomInventory(ctx context.Context, kingdomID string) (map[string]int, error)
	GetPersonalInventory(ctx context.Context, userID string) (map[string]int, error)
}

type SetupStore interface {
	SetupDemo(ctx context.Context) error
}

type LeaderboardStore interface {
	AddScore(ctx context.Context, kingdomID, userID string, delta int) error
	TopLeaderboard(ctx context.Context, kingdomID string, limit int) ([]LeaderboardRow, error)
}

type LeaderboardRow struct {
	UserID string `json:"userId"`
	Score  int64  `json:"score"`
}
