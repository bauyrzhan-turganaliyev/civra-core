CREATE TABLE IF NOT EXISTS personal_inventory (
  user_id  TEXT NOT NULL,
  resource TEXT NOT NULL,
  quantity INT  NOT NULL,
  PRIMARY KEY (user_id, resource)
);
