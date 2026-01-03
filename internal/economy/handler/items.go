package handler

import (
	"encoding/json"
	"net/http"

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
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}

	items, err := h.pg.ListUserItems(r.Context(), userID)
	if err != nil {
		httpx.JSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	httpx.JSON(w, 200, map[string]any{
		"userId": userID,
		"items":  items,
	})
}

// POST /items/craft-tool  { "userId":"u1", "tier":1 }
func (h *ItemsHandler) CraftTool(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}

	var req struct {
		Tier int `json:"tier"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}
	if req.Tier == 0 {
		httpx.JSON(w, 400, map[string]string{"error": "tier is required"})
		return
	}

	spec, ok := service.ToolSpecForUI(req.Tier)
	if !ok {
		httpx.JSON(w, 400, map[string]string{"error": "invalid tier"})
		return
	}

	id, err := h.pg.CraftTool(
		r.Context(),
		userID,
		spec.Tier,
		spec.BonusPct,
		spec.MaxDur,
		spec.IronCost,
		spec.WoodCost,
	)
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
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}

	var req struct {
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

	if err := h.pg.EquipItem(r.Context(), userID, id); err != nil {
		httpx.JSON(w, 404, map[string]string{"error": err.Error()})
		return
	}

	httpx.JSON(w, 200, map[string]any{"ok": true})
}
