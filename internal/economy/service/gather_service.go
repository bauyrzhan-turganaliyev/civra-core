package service

import (
	"context"
	"errors"
	"time"

	"civra-core/internal/economy/entity"
	"civra-core/internal/economy/repository"
)

var ErrInvalidResourceForProfession = errors.New("invalid resource for profession")

type GatherService struct {
	store repository.GatherStore
	now   func() time.Time
}

func NewGatherService(store repository.GatherStore) *GatherService {
	return &GatherService{
		store: store,
		now:   time.Now,
	}
}

type GatherResult struct {
	ToKingdomInventory int  `json:"toKingdomInventory"`
	ToPersonal         int  `json:"toPersonal"`
	QuotaDone          bool `json:"quotaDone"`
	QuotaProgress      int  `json:"quotaProgress"`
	QuotaRequired      int  `json:"quotaRequired"`
	ToolBonusPct       int  `json:"toolBonusPct"`
	ToolUsed           bool `json:"toolUsed"`
}

func quotaResourceFor(p entity.Profession) (entity.Resource, bool) {
	switch p {
	case entity.ProfFarmer:
		return entity.ResFood, true
	case entity.ProfMiner:
		return entity.ResIron, true
	case entity.ProfLumber:
		return entity.ResWood, true
	default:
		return "", false
	}
}

const quotaRequired = 1000 // MVP fixed

// Core rule:
// While quota not done -> gathered amount goes to Kingdom Inventory (and counts to quota).
// After quota done -> gathered amount goes to personal inventory.
func (s *GatherService) Gather(
	ctx context.Context,
	userID, kingdomID string,
	prof entity.Profession,
	res entity.Resource,
	amount int,
) (GatherResult, error) {
	if amount <= 0 {
		return GatherResult{}, errors.New("amount must be positive")
	}

	qRes, ok := quotaResourceFor(prof)
	if !ok {
		return GatherResult{}, errors.New("profession not supported in MVP gather")
	}
	if res != qRes {
		return GatherResult{}, ErrInvalidResourceForProfession
	}

	toKingdom, toPersonal, progress, toolBonusPct, toolUsed, err := s.store.Gather(
		ctx,
		userID,
		kingdomID,
		string(res),
		quotaRequired,
		s.now(),
		amount,
	)

	if err != nil {
		return GatherResult{}, err
	}

	return GatherResult{
		ToKingdomInventory: toKingdom,
		ToPersonal:         toPersonal,
		QuotaDone:          progress >= quotaRequired,
		QuotaProgress:      progress,
		QuotaRequired:      quotaRequired,
		ToolBonusPct:       toolBonusPct,
		ToolUsed:           toolUsed,
	}, nil
}
