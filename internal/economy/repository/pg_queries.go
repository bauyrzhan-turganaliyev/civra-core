package repository

import "context"

func (s *PgStore) GetKingdomInventory(ctx context.Context, kingdomID string) (map[string]int, error) {
	rows, err := s.db.Query(ctx,
		`SELECT resource, quantity FROM kingdom_inventory WHERE kingdom_id=$1`,
		kingdomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string]int{}
	for rows.Next() {
		var resource string
		var qty int
		if err := rows.Scan(&resource, &qty); err != nil {
			return nil, err
		}
		out[resource] = qty
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
