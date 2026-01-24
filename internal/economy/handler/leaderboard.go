package handler

import (
	"civra-core/internal/economy/service"
	"civra-core/pkg/httpx"
	"net/http"
)

type LeaderboardHandler struct {
	svc *service.LeaderboardService
}

func NewLeaderboardHandler(svc *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{svc: svc}
}

func (h *LeaderboardHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.JSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}

	kingdomID := r.Header.Get("X-Kingdom-Id")
	if kingdomID == "" {
		httpx.JSON(w, 400, map[string]string{"error": "missing kingdom"})
		return
	}

	rows, err := h.svc.Get(r.Context(), kingdomID)
	if err != nil {
		httpx.JSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	httpx.JSON(w, 200, map[string]any{
		"kingdomId": kingdomID,
		"leaders":   rows,
	})
}
