# Civra — Manual & Concurrency Test Cases

Gateway base URL: http://localhost:8080

---

## 0) Prerequisites

### Start the system
make up

### Check gateway
curl -s http://localhost:8080/healthz && echo

Expected: {"ok":true,"service":"gateway"}

### (Optional) Reset
make reset
make up

---

## 1) Daily quota — partial progress

curl -s -X POST http://localhost:8080/economy/gather \
 -H "Content-Type: application/json" \
 -d '{"userId":"u_q1","kingdomId":"k1","profession":"farmer","resource":"food","amount":600}' && echo

Expected:
- toKingdomInventory = 600
- toPersonal = 0
- quotaProgress = 600
- quotaDone = false

---

## 2) Daily quota — completion and split

curl -s -X POST http://localhost:8080/economy/gather \
 -H "Content-Type: application/json" \
 -d '{"userId":"u_q2","kingdomId":"k1","profession":"farmer","resource":"food","amount":960}' && echo

curl -s -X POST http://localhost:8080/economy/gather \
 -H "Content-Type: application/json" \
 -d '{"userId":"u_q2","kingdomId":"k1","profession":"farmer","resource":"food","amount":60}' && echo

Expected:
- toKingdomInventory = 40
- toPersonal = 20
- quotaDone = true

---

## 3) Tool bonus and durability

curl -s -X POST http://localhost:8080/economy/items/craft-tool \
 -H "Content-Type: application/json" \
 -d '{"userId":"u_tool","tier":1}' && echo

curl -s -X POST http://localhost:8080/economy/items/equip \
 -H "Content-Type: application/json" \
 -d '{"userId":"u_tool","itemId":"<ITEM_ID>"}' && echo

curl -s -X POST http://localhost:8080/economy/gather \
 -H "Content-Type: application/json" \
 -d '{"userId":"u_tool","kingdomId":"k1","profession":"farmer","resource":"food","amount":60}' && echo

Verify durability:
curl -s http://localhost:8080/economy/items?userId=u_tool && echo

---

## 4) Item market — sell item

curl -s -X POST http://localhost:8080/economy/market/items/sell \
 -H "Content-Type: application/json" \
 -d '{"kingdomId":"k1","sellerId":"u_tool","itemId":"<ITEM_ID>","price":100}' && echo

curl -s http://localhost:8080/economy/market/items/orders?kingdomId=k1 && echo

---

## 5) Concurrency test — item buy

for i in {1..10}; do
 curl -s -X POST http://localhost:8080/economy/market/items/buy \
  -H "Content-Type: application/json" \
  -d '{"orderId":"<ORDER_ID>","buyerId":"buyer_'$i'"}' &
done
wait
echo "done"

Expected: exactly one success.

---

## 6) Cancel item order

curl -s -X POST http://localhost:8080/economy/market/items/cancel \
 -H "Content-Type: application/json" \
 -d '{"orderId":"<ORDER_ID>","sellerId":"u_tool"}' && echo

Verify removal:
curl -s http://localhost:8080/economy/market/items/orders?kingdomId=k1 && echo
