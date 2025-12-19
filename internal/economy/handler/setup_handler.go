package handler

import (
	"net/http"

	"civra-core/internal/economy/service"
	"civra-core/pkg/httpx"
)

type SetupHandler struct {
	svc *service.SetupService
}

func NewSetupHandler(svc *service.SetupService) *SetupHandler {
	return &SetupHandler{svc: svc}
}

// POST /setup/demo
func (h *SetupHandler) Demo(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Demo(r.Context()); err != nil {
		httpx.Err(w, 500, err.Error())
		return
	}
	httpx.JSON(w, 200, map[string]any{
		"ok":       true,
		"kingdoms": []string{"k1"},
		"users": []map[string]string{
			{"userId": "u1", "mainProfession": "farmer"},
			{"userId": "u2", "mainProfession": "miner"},
		},
	})
}
