CREATE TABLE IF NOT EXISTS leaderboard (
  kingdom_id TEXT NOT NULL,
  user_id    TEXT NOT NULL,
  score      BIGINT NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (kingdom_id, user_id)
);
