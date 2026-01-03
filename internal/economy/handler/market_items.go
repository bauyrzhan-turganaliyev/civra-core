package handler

import (
	"encoding/json"
	"net/http"

	"civra-core/internal/economy/repository"
	"civra-core/pkg/httpx"

	"github.com/google/uuid"
)

type MarketItemsHandler struct {
	pg *repository.PgStore
}

func NewMarketItemsHandler(pg *repository.PgStore) *MarketItemsHandler {
	return &MarketItemsHandler{pg: pg}
}

// GET /market/items/orders?kingdomId=k1
func (h *MarketItemsHandler) Orders(w http.ResponseWriter, r *http.Request) {
	kingdomID := r.Header.Get("X-Kingdom-Id")
	if kingdomID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}

	orders, err := h.pg.ListItemOrders(r.Context(), kingdomID)
	if err != nil {
		httpx.JSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	httpx.JSON(w, 200, map[string]any{"orders": orders})
}

// POST /market/items/buy { "orderId":"uuid", "buyerId":"u2" }
func (h *MarketItemsHandler) Buy(w http.ResponseWriter, r *http.Request) {
	buyerID := r.Header.Get("X-User-Id")
	if buyerID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}

	var req struct {
		OrderID string `json:"orderId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}

	id, err := uuid.Parse(req.OrderID)
	if err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid orderId"})
		return
	}

	if err := h.pg.BuyItem(r.Context(), id, buyerID); err != nil {
		httpx.JSON(w, 409, map[string]string{"error": err.Error()})
		return
	}

	httpx.JSON(w, 200, map[string]any{"ok": true})
}

// POST /market/items/sell
func (h *MarketItemsHandler) Sell(w http.ResponseWriter, r *http.Request) {
	kingdomID := r.Header.Get("X-Kingdom-Id")
	sellerID := r.Header.Get("X-User-Id")
	if kingdomID == "" || sellerID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}

	var req struct {
		ItemID string `json:"itemId"`
		Price  int    `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}

	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid itemId"})
		return
	}

	if req.Price <= 0 {
		httpx.JSON(w, 400, map[string]string{"error": "price must be positive"})
		return
	}

	if err := h.pg.SellItem(
		r.Context(),
		kingdomID,
		sellerID,
		itemID,
		req.Price,
	); err != nil {
		httpx.JSON(w, 409, map[string]string{"error": err.Error()})
		return
	}

	httpx.JSON(w, 200, map[string]any{"ok": true})
}

// POST /market/items/cancel { "orderId":"uuid", "sellerId":"u1" }
func (h *MarketItemsHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	sellerID := r.Header.Get("X-User-Id")
	if sellerID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}

	var req struct {
		OrderID string `json:"orderId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}

	id, err := uuid.Parse(req.OrderID)
	if err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid orderId"})
		return
	}

	if err := h.pg.CancelItemSale(r.Context(), id, sellerID); err != nil {
		httpx.JSON(w, 409, map[string]string{"error": err.Error()})
		return
	}

	httpx.JSON(w, 200, map[string]any{"ok": true})
}
