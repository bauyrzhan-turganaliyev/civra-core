# Concurrency and Consistency

## Problem Statement
The system must handle multiple concurrent requests modifying shared state:
- Daily quotas
- Inventories
- Market orders

Naive implementations lead to race conditions and inconsistent data.

---

## Solution Strategy

### Database Transactions
All critical operations are executed inside PostgreSQL transactions.

### Row-Level Locking
`SELECT ... FOR UPDATE` is used to:
- Lock quota rows
- Lock market orders
- Lock personal inventory rows

This ensures serialized access to shared data.

---

## Handling Race Conditions

### Quota Initialization
To avoid concurrent insert conflicts:
```sql
INSERT INTO quota_progress (...)
ON CONFLICT DO NOTHING
