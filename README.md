# Civra – Distributed Economy Backend

Civra is a distributed backend system for a browser-based MMO game.
The project was developed as the final project for the course:

**Distributed Programming for Web, IoT and Mobile Systems (B032427 / B335)**  
University of Florence

---

## Project Overview

The system simulates a multi-player economy where:
- Players gather resources
- Daily quotas are enforced
- Resources are stored at kingdom and personal level
- Players trade resources on a local market

The main focus of the project is **data consistency under concurrency**.

---

## Architecture

- **Gateway Service** – HTTP entry point
- **Economy Service** – core business logic
- **PostgreSQL** – shared persistent state



The Economy Service follows a layered architecture:
- Handler layer (HTTP)
- Service layer (business rules)
- Repository layer (data access & transactions)

---

## Implemented Features

### Economy
- Daily quota enforcement per player
- Automatic split between Kingdom Inventory and Personal Inventory
- Persistent storage in PostgreSQL
- Transaction-safe updates

### Market
- Sell orders (resources removed immediately)
- Buy orders (single winner under concurrency)
- Market order listing
- Fully transactional execution

---

## Concurrency Handling

The system handles concurrent access using:
- PostgreSQL transactions
- Row-level locking (`SELECT ... FOR UPDATE`)
- Atomic UPSERT operations (`ON CONFLICT DO NOTHING`)

Race conditions were intentionally tested using parallel HTTP requests.

---

## Technologies Used

- Go
- PostgreSQL
- Docker & Docker Compose
- HTTP / JSON

---

## How to Run

### Requirements
- Docker
- Docker Compose

### Start the system
```bash
docker compose up --build
