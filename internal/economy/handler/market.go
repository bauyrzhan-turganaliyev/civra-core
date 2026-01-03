package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"civra-core/internal/economy/repository"
	"civra-core/pkg/httpx"

	"github.com/google/uuid"
)

type MarketHandler struct {
	store *repository.PgStore
}

func NewMarketHandler(store *repository.PgStore) *MarketHandler {
	return &MarketHandler{store: store}
}

type SellReq struct {
	KingdomID string `json:"kingdomId"`
	SellerID  string `json:"sellerId"`
	Resource  string `json:"resource"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}

type CancelReq struct {
	OrderID  string `json:"orderId"`
	SellerID string `json:"sellerId"`
}

func (h *MarketHandler) Sell(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	kingdomID := r.Header.Get("X-Kingdom-Id")
	if userID == "" || kingdomID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}

	var req struct {
		Resource string `json:"resource"`
		Quantity int    `json:"quantity"`
		Price    int    `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}

	id, err := h.store.CreateSellOrder(
		r.Context(),
		kingdomID,
		userID,
		req.Resource,
		req.Quantity,
		req.Price,
	)
	if err != nil {
		if err == repository.ErrNotEnoughResource {
			httpx.JSON(w, 409, map[string]string{"error": err.Error()})
			return
		}
		httpx.JSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	httpx.JSON(w, 200, map[string]string{"orderId": id.String()})
}

type BuyReq struct {
	OrderID string `json:"orderId"`
	BuyerID string `json:"buyerId"`
}

func (h *MarketHandler) Buy(w http.ResponseWriter, r *http.Request) {
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

	if err := h.store.BuyOrder(r.Context(), id, buyerID); err != nil {
		httpx.JSON(w, 404, map[string]string{"error": "order not found"})
		return
	}

	httpx.JSON(w, 200, map[string]any{"ok": true})
}

func (h *MarketHandler) Orders(w http.ResponseWriter, r *http.Request) {
	kingdomID := r.Header.Get("X-Kingdom-Id")
	if kingdomID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}

	limit := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	orders, err := h.store.ListMarketOrders(r.Context(), kingdomID, limit)
	if err != nil {
		httpx.JSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	httpx.JSON(w, 200, map[string]any{
		"kingdomId": kingdomID,
		"count":     len(orders),
		"orders":    orders,
	})
}

func (h *MarketHandler) Cancel(w http.ResponseWriter, r *http.Request) {
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

	if err := h.store.CancelSellOrder(r.Context(), id, sellerID); err != nil {
		httpx.JSON(w, 404, map[string]string{"error": "order not found or not owner"})
		return
	}

	httpx.JSON(w, 200, map[string]any{"ok": true})
}
