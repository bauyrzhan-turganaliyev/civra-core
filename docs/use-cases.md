
---

# 📄 `docs/use-cases.md`

```md
# Use Cases

## Use Case 1 – Daily Resource Gathering

**Actor:** Player  
**Description:**  
A player gathers resources according to their profession.  
Until the daily quota is reached, resources are sent to the Kingdom Inventory.  
After the quota is completed, resources are stored in the Personal Inventory.

**Steps:**
1. Player sends a gather request
2. System checks daily quota
3. Resources are split between:
   - Kingdom Inventory
   - Personal Inventory
4. Quota progress is updated atomically

**Result:**  
Quota is enforced consistently even under concurrent requests.

---

## Use Case 2 – Selling Resources on Local Market

**Actor:** Player (Seller)

**Steps:**
1. Player creates a sell order
2. Resources are immediately removed from Personal Inventory
3. A market order is created

**Guarantee:**  
Resources cannot be sold twice.

---

## Use Case 3 – Buying from Market

**Actor:** Player (Buyer)

**Steps:**
1. Buyer selects a market order
2. Order row is locked
3. Resources are transferred to buyer
4. Order is removed

**Guarantee:**  
Only one buyer can complete the order.
