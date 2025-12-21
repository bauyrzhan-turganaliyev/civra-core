package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"civra-core/internal/economy/repository"
	"civra-core/internal/economy/service"
	"civra-core/pkg/httpx"

	"github.com/google/uuid"
)

type ItemsHandler struct {
	pg *repository.PgStore
}

func NewItemsHandler(pg *repository.PgStore) *ItemsHandler { return &ItemsHandler{pg: pg} }

// GET /items?userId=u1
func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		httpx.JSON(w, 400, map[string]string{"error": "missing userId"})
		return
	}
	items, err := h.pg.ListUserItems(r.Context(), userID)
	if err != nil {
		httpx.JSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	httpx.JSON(w, 200, map[string]any{"userId": userID, "items": items})
}

// POST /items/craft-tool  { "userId":"u1", "tier":1 }
func (h *ItemsHandler) CraftTool(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"userId"`
		Tier   int    `json:"tier"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}
	if req.UserID == "" || req.Tier == 0 {
		httpx.JSON(w, 400, map[string]string{"error": "userId and tier are required"})
		return
	}

	spec, ok := service.ToolSpecForUI(req.Tier) // см. ниже как сделать без экспорта лишнего
	if !ok {
		httpx.JSON(w, 400, map[string]string{"error": "invalid tier"})
		return
	}

	id, err := h.pg.CraftTool(r.Context(), req.UserID, spec.Tier, spec.BonusPct, spec.MaxDur, spec.IronCost, spec.WoodCost)
	if err != nil {
		if err == repository.ErrNotEnoughMaterials {
			httpx.JSON(w, 409, map[string]string{"error": err.Error()})
			return
		}
		httpx.JSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	httpx.JSON(w, 200, map[string]any{"itemId": id.String()})
}

// POST /items/equip { "userId":"u1", "itemId":"uuid" }
func (h *ItemsHandler) Equip(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"userId"`
		ItemID string `json:"itemId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}

	id, err := uuid.Parse(req.ItemID)
	if err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid itemId"})
		return
	}

	if err := h.pg.EquipItem(r.Context(), req.UserID, id); err != nil {
		httpx.JSON(w, 404, map[string]string{"error": err.Error()})
		return
	}
	httpx.JSON(w, 200, map[string]any{"ok": true})
}

// tiny helper for query ints (если захочешь позже)
func qInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
