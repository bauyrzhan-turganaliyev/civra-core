CREATE TABLE IF NOT EXISTS market_item_orders (
  id UUID PRIMARY KEY,
  kingdom_id TEXT NOT NULL,
  seller_id TEXT NOT NULL,
  item_id UUID NOT NULL UNIQUE,
  price INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);
