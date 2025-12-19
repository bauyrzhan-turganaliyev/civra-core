package service

import (
	"civra-core/internal/economy/repository"
	"context"
)

type SetupService struct {
	store repository.SetupStore
}

func NewSetupService(store repository.SetupStore) *SetupService {
	return &SetupService{store: store}
}

func (s *SetupService) Demo(ctx context.Context) error {
	return s.store.SetupDemo(ctx)
}
