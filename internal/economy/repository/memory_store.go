package repository

import (
	"context"
	"sync"
	"time"
)

type MemoryStore struct {
	mu sync.Mutex

	// KingdomInventory[kingdomID][resource] = qty
	KingdomInventory map[string]map[string]int

	// PersonalInventory[userID][resource] = qty
	PersonalInventory map[string]map[string]int

	// QuotaProgress[userID][date][resource] = progress
	QuotaProgress map[string]map[string]map[string]int

	// MainProfession[userID] = "farmer" / "miner" / "lumber" ...
	MainProfession map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		KingdomInventory:  make(map[string]map[string]int),
		PersonalInventory: make(map[string]map[string]int),
		QuotaProgress:     make(map[string]map[string]map[string]int),
		MainProfession:    make(map[string]string),
	}
}

func (s *MemoryStore) Gather(
	ctx context.Context,
	userID, kingdomID, resource string,
	quotaRequired int,
	now time.Time,
	amount int,
) (toKingdom, toPersonal, progress int) {

	s.mu.Lock()
	defer s.mu.Unlock()

	date := now.Format("2006-01-02")

	if s.KingdomInventory[kingdomID] == nil {
		s.KingdomInventory[kingdomID] = make(map[string]int)
	}
	if s.PersonalInventory[userID] == nil {
		s.PersonalInventory[userID] = make(map[string]int)
	}
	if s.QuotaProgress[userID] == nil {
		s.QuotaProgress[userID] = make(map[string]map[string]int)
	}
	if s.QuotaProgress[userID][date] == nil {
		s.QuotaProgress[userID][date] = make(map[string]int)
	}

	progress = s.QuotaProgress[userID][date][resource]
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

	if toKingdom > 0 {
		s.KingdomInventory[kingdomID][resource] += toKingdom
		s.QuotaProgress[userID][date][resource] += toKingdom
		progress += toKingdom
	}
	if toPersonal > 0 {
		s.PersonalInventory[userID][resource] += toPersonal
	}

	return
}

func todayKey(now time.Time) string {
	return now.Format("2006-01-02")
}
