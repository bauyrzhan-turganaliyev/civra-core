package main

import (
	"log"
	"net/http"

	"civra-core/pkg/config"
	"civra-core/pkg/httpx"
)

func main() {
	cfg := config.LoadService()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, 200, map[string]any{"service": "economy", "ok": true})
	})

	addr := ":" + cfg.Port
	log.Printf("economy listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
