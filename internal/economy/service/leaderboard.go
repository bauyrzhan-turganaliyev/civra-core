package service

import (
	"civra-core/internal/economy/repository"
	"context"
)

type LeaderboardService struct {
	store *repository.PgStore
}

func NewLeaderboardService(store *repository.PgStore) *LeaderboardService {
	return &LeaderboardService{store: store}
}

func (s *LeaderboardService) Get(
	ctx context.Context,
	kingdomID string,
) ([]repository.LeaderboardRow, error) {
	return s.store.GetLeaderboard(ctx, kingdomID, 20) // top-20
}
