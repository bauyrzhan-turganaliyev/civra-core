package repository

import (
	"context"
)

func (s *PgStore) GetLeaderboard(
	ctx context.Context,
	kingdomID string,
	limit int,
) ([]LeaderboardRow, error) {

	rows, err := s.db.Query(ctx, `
		SELECT user_id, score
		FROM leaderboard
		WHERE kingdom_id = $1
		ORDER BY score DESC
		LIMIT $2
	`, kingdomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []LeaderboardRow
	for rows.Next() {
		var r LeaderboardRow
		if err := rows.Scan(&r.UserID, &r.Score); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}
