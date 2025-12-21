CREATE TABLE IF NOT EXISTS user_items (
  id UUID PRIMARY KEY,
  user_id TEXT NOT NULL,
  item_type TEXT NOT NULL,          -- "tool"
  tier INT NOT NULL,                -- 1..3
  durability INT NOT NULL,
  max_durability INT NOT NULL,
  bonus_pct INT NOT NULL,           -- 10, 20, 35
  equipped BOOLEAN NOT NULL DEFAULT FALSE,
  listed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_user_one_equipped
ON user_items(user_id)
WHERE equipped = TRUE;
