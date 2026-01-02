# Civra — Distributed Programming Final Project Report

## 1. Project idea
Civra is a browser-based MMO prototype where many players cooperate inside a kingdom. Players gather resources, fulfill a daily quota for the kingdom inventory, craft RPG-like tools with durability, and trade resources and items through a local market. A gateway provides a single entry point and manages lightweight sessions with HTTP-only cookies.

## 2. Architecture and design
### 2.1 High-level architecture
- **Gateway**: reverse proxy + static UI + session cookie endpoints
- **Economy service**: gathering quota, inventories, tools (items), item market endpoints
- **Market service**: resource market (sell/buy/cancel orders) *(if you kept it separate)*
- **Kingdom service**: basic kingdom endpoints / placeholder for future expansion *(if present)*
- **PostgreSQL**: persistence and transactional concurrency control

### 2.2 Main components
- `cmd/gateway`: routing, `/auth/login`, `/auth/me`, `/auth/logout`, reverse proxy
- `cmd/economy`: economy HTTP API
- `internal/economy/service`: business rules (quota, tool bonus)
- `internal/economy/repository`: PostgreSQL transactions, row locks (`FOR UPDATE`), upserts
- `ui/`: demo frontend served by gateway

### 2.3 Data model (main tables)
- `kingdom_inventory(kingdom_id, resource, quantity)`
- `personal_inventory(user_id, resource, quantity)`
- `quota_progress(user_id, day, resource, progress)`
- `user_items(id, user_id, item_type, tier, durability, max_durability, bonus_pct, equipped, listed, created_at)`
- `market_orders(...)` for resource market
- `market_item_orders(...)` for item market

## 3. Technologies and tools
- Go (standard library `net/http`, `httputil.ReverseProxy`)
- PostgreSQL + pgx
- Docker + Docker Compose
- Basic HTML/CSS/JS frontend
- Makefile for one-command run

## 4. Course topics covered
- Distributed services and separation of responsibilities (gateway vs services)
- REST-style APIs over HTTP
- Reverse proxy and single entry point
- Persistence and database transactions
- Concurrency control (`FOR UPDATE`, upsert, unique constraints)
- Containerization and reproducible deployment with Docker Compose
- Session management using HTTP-only cookies (gateway-managed)

## 5. Use case (end-to-end)
### 5.1 Daily quota & gathering
A player gathers the profession resource. Until the daily quota is completed, gathered amount is deposited into **Kingdom Inventory** and counted toward quota. After quota completion, additional gathered resources go to the player’s personal inventory.

### 5.2 Tools (RPG items)
A player crafts a tool (tier-based). If equipped, the tool increases gathering output by a bonus percentage and loses durability each gathering action. When durability reaches 0, the tool breaks.

### 5.3 Trading
Players can:
- sell/buy/cancel resource orders (local market)
- sell/buy/cancel tool items through item market orders
All buy operations are transactional, so in concurrent purchases only one buyer succeeds.

## 6. How to run
See `README.md`.

## 7. Test cases
See `TESTS.md`

## 8. Notes and future work
- roles/permissions can be expanded (king decisions, builder projects)
- real auth can replace demo sessions
- additional item rarities, recipes, and balance rules
