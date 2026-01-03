package handler

import (
	"encoding/json"
	"net/http"

	"civra-core/internal/economy/entity"
	"civra-core/internal/economy/service"
	"civra-core/pkg/httpx"
)

type GatherHandler struct {
	svc *service.GatherService
}

func NewGatherHandler(svc *service.GatherService) *GatherHandler {
	return &GatherHandler{svc: svc}
}

type GatherRequest struct {
	UserID     string `json:"userId"`
	KingdomID  string `json:"kingdomId"`
	Profession string `json:"profession"`
	Resource   string `json:"resource"`
	Amount     int    `json:"amount"`
}

func (h *GatherHandler) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	kingdomID := r.Header.Get("X-Kingdom-Id")
	if userID == "" || kingdomID == "" {
		httpx.Err(w, 401, "no session")
		return
	}

	var req struct {
		Profession string `json:"profession"`
		Resource   string `json:"resource"`
		Amount     int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Err(w, 400, "invalid json")
		return
	}

	res, err := h.svc.Gather(
		r.Context(),
		userID,
		kingdomID,
		entity.Profession(req.Profession),
		entity.Resource(req.Resource),
		req.Amount,
	)
	if err != nil {
		httpx.Err(w, 400, err.Error())
		return
	}

	httpx.JSON(w, 200, res)
}
