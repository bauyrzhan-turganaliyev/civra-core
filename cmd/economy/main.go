package main

import (
	"context"
	"log"
	"net/http"

	"civra-core/internal/economy/handler"
	"civra-core/internal/economy/repository"
	"civra-core/internal/economy/service"
	"civra-core/pkg/config"
	"civra-core/pkg/db"
	"civra-core/pkg/httpx"
)

func main() {
	cfg := config.LoadService()

	dbPool, err := db.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	mem := repository.NewMemoryStore()
	pg := repository.NewPgStore(dbPool)

	gatherSvc := service.NewGatherService(pg)
	gatherH := handler.NewGatherHandler(gatherSvc)

	querySvc := service.NewInventoryQueryService(pg)
	queryH := handler.NewInventoryQueryHandler(querySvc)

	setupSvc := service.NewSetupService(mem)
	setupH := handler.NewSetupHandler(setupSvc)

	marketH := handler.NewMarketHandler(pg)

	// ✅ register routes
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, 200, map[string]any{"service": "economy", "ok": true})
	})

	mux.HandleFunc("/gather", gatherH.Handle)
	mux.HandleFunc("/kingdom-inventory", queryH.Kingdom)
	mux.HandleFunc("/personal-inventory", queryH.Personal)
	mux.HandleFunc("/setup/demo", setupH.Demo)
	mux.HandleFunc("/market/sell", marketH.Sell)
	mux.HandleFunc("/market/buy", marketH.Buy)
	mux.HandleFunc("/market/orders", marketH.Orders)

	// ✅ then start server
	addr := ":" + cfg.Port
	log.Printf("economy listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
