package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"civra-core/pkg/config"
	"civra-core/pkg/httpx"
)

func main() {
	cfg := config.LoadGateway()

	kingdom := mustProxy(cfg.KingdomURL)
	economy := mustProxy(cfg.EconomyURL)
	market := mustProxy(cfg.MarketURL)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, 200, map[string]any{"service": "gateway", "ok": true})
	})

	// Proxy routes (через один вход)
	mux.Handle("/kingdom/", http.StripPrefix("/kingdom", kingdom))
	mux.Handle("/economy/", http.StripPrefix("/economy", economy))
	mux.Handle("/market/", http.StripPrefix("/market", market))

	// Serve demo UI (static)
	mux.Handle("/", http.FileServer(http.Dir("./ui")))

	addr := ":" + cfg.Port
	log.Printf("gateway listening on %s", addr)
	log.Printf("proxy: kingdom=%s economy=%s market=%s", cfg.KingdomURL, cfg.EconomyURL, cfg.MarketURL)
	log.Fatal(http.ListenAndServe(addr, logRequests(mux)))
}

func mustProxy(raw string) *httputil.ReverseProxy {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if !strings.HasPrefix(path, "/healthz") {
			log.Printf("%s %s", r.Method, path)
		}
		next.ServeHTTP(w, r)
	})
}
