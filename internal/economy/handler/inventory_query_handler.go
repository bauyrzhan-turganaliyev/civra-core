package handler

import (
	"net/http"

	"civra-core/internal/economy/service"
	"civra-core/pkg/httpx"
)

type InventoryQueryHandler struct {
	svc *service.InventoryQueryService
}

func NewInventoryQueryHandler(svc *service.InventoryQueryService) *InventoryQueryHandler {
	return &InventoryQueryHandler{svc: svc}
}

// GET /kingdom-inventory?kingdomId=k1
func (h *InventoryQueryHandler) Kingdom(w http.ResponseWriter, r *http.Request) {
	kingdomID := r.Header.Get("X-Kingdom-Id")
	if kingdomID == "" {
		httpx.Err(w, 401, "no session")
		return
	}

	inv, err := h.svc.KingdomInventory(r.Context(), kingdomID)
	if err != nil {
		httpx.Err(w, 500, err.Error())
		return
	}

	httpx.JSON(w, 200, map[string]any{
		"kingdomId": kingdomID,
		"inventory": inv,
	})
}

// GET /personal-inventory?userId=u1
func (h *InventoryQueryHandler) Personal(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		httpx.Err(w, 401, "no session")
		return
	}

	inv, err := h.svc.PersonalInventory(r.Context(), userID)
	if err != nil {
		httpx.Err(w, 500, err.Error())
		return
	}

	httpx.JSON(w, 200, map[string]any{
		"userId":    userID,
		"inventory": inv,
	})
}
