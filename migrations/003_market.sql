CREATE TABLE IF NOT EXISTS market_orders (
  id UUID PRIMARY KEY,
  kingdom_id TEXT NOT NULL,
  seller_id TEXT NOT NULL,
  resource TEXT NOT NULL,
  quantity INT NOT NULL,
  price INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);
