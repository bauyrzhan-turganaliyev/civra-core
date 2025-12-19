package service

import (
	"civra-core/internal/economy/repository"
	"context"
)

type InventoryQueryService struct {
	store repository.InventoryQueryStore
}

func NewInventoryQueryService(store repository.InventoryQueryStore) *InventoryQueryService {
	return &InventoryQueryService{store: store}
}
func (s *InventoryQueryService) KingdomInventory(ctx context.Context, kingdomID string) (map[string]int, error) {
	return s.store.GetKingdomInventory(ctx, kingdomID)
}

func (s *InventoryQueryService) PersonalInventory(ctx context.Context, userID string) (map[string]int, error) {
	return s.store.GetPersonalInventory(ctx, userID)
}
