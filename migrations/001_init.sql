CREATE TABLE kingdom_inventory (
  kingdom_id TEXT NOT NULL,
  resource   TEXT NOT NULL,
  quantity   INT  NOT NULL,
  PRIMARY KEY (kingdom_id, resource)
);

CREATE TABLE quota_progress (
  user_id    TEXT NOT NULL,
  day        DATE NOT NULL,
  resource   TEXT NOT NULL,
  progress   INT  NOT NULL,
  PRIMARY KEY (user_id, day, resource)
);
