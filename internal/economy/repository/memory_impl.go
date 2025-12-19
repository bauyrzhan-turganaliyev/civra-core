package repository

import "context"

// InventoryQueryStore
func (s *MemoryStore) GetKingdomInventory(ctx context.Context, kingdomID string) (map[string]int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	src := s.KingdomInventory[kingdomID]
	out := make(map[string]int, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out, nil
}

func (s *MemoryStore) GetPersonalInventory(ctx context.Context, userID string) (map[string]int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	src := s.PersonalInventory[userID]
	out := make(map[string]int, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out, nil
}

// SetupStore
func (s *MemoryStore) SetupDemo(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.KingdomInventory["k1"] == nil {
		s.KingdomInventory["k1"] = make(map[string]int)
	}
	if s.PersonalInventory["u1"] == nil {
		s.PersonalInventory["u1"] = make(map[string]int)
	}
	if s.PersonalInventory["u2"] == nil {
		s.PersonalInventory["u2"] = make(map[string]int)
	}
	if s.MainProfession == nil {
		s.MainProfession = make(map[string]string)
	}
	s.MainProfession["u1"] = "farmer"
	s.MainProfession["u2"] = "miner"

	return nil
}
