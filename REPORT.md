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

## 4. Course topics covered (revised)
 - Distributed systems architecture  
    The project is implemented as a set of independent services (gateway, economy, market, kingdom) communicating over the network. Each service has a clear responsibility and runs in its own container, reflecting the core principles of distributed systems.
 - RESTful web services
    All functionality is exposed through REST-style HTTP APIs using standard methods (GET, POST) and JSON payloads. Resources such as inventories, items, and market orders are modeled as web resources, following REST design principles.
 - Communication mechanisms
    Services communicate synchronously using HTTP request/response interactions. The gateway acts as a middleware component and a single entry point, abstracting internal service communication from the client.
 - Session management and middleware
    User sessions are managed at the gateway level using HTTP-only cookies. The gateway validates sessions and propagates user identity to downstream services, keeping backend services stateless with respect to authentication.
 - Data access and persistence
    Persistent state is stored in PostgreSQL and accessed through a repository layer. Database transactions are used to ensure consistency when modifying inventories, quotas, and market data, as discussed in the Accessing Data lecture.
 - Concurrency and synchronization
    Concurrent operations (e.g., gathering resources or buying from the market) are handled using transactional database access and row-level locking. This guarantees correct synchronization when multiple users access shared data.
 - Broadcast and coordination concepts
    Kingdom inventories and market orders represent shared state visible to all players in a kingdom, conceptually corresponding to broadcast-style dissemination of information within a distributed system.
 - Leader and coordination role
    The gateway plays the role of a logical coordinator by centralizing authentication and request routing. While explicit leader election algorithms are not implemented, the design follows the coordinator-based model discussed in the course.
 - RPC-style service interaction
    Internal service communication follows an RPC-style model over REST, with clearly defined service interfaces. The architecture is compatible with gRPC-style communication as presented in the lectures.
 - Pervasive systems principles
    The system continuously reacts to user actions and context (daily quotas, equipped tools, durability), demonstrating autonomous and context-aware behavior typical of pervasive and distributed systems.

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
