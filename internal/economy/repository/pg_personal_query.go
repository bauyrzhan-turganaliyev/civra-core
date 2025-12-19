package repository

import "context"

func (s *PgStore) GetPersonalInventory(ctx context.Context, userID string) (map[string]int, error) {
	rows, err := s.db.Query(ctx,
		`SELECT resource, quantity FROM personal_inventory WHERE user_id=$1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string]int{}
	for rows.Next() {
		var res string
		var qty int
		if err := rows.Scan(&res, &qty); err != nil {
			return nil, err
		}
		out[res] = qty
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
