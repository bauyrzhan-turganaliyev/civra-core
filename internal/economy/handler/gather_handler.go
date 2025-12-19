package handler

import (
	"context"
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
	var req GatherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Err(w, 400, "invalid json")
		return
	}

	res, err := h.svc.Gather(
		context.Background(),
		req.UserID,
		req.KingdomID,
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
