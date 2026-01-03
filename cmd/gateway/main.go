package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"civra-core/pkg/config"
	"civra-core/pkg/httpx"
)

func main() {
	cfg := config.LoadGateway()

	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, 200, map[string]any{"service": "gateway", "ok": true})
	})

	// auth
	mux.HandleFunc("/auth/login", handleLogin)
	mux.HandleFunc("/auth/me", handleMe)
	mux.HandleFunc("/auth/logout", handleLogout)

	// proxies
	kingdom := mustProxyWithSessionHeaders(cfg.KingdomURL)
	economy := mustProxyWithSessionHeaders(cfg.EconomyURL)
	market := mustProxyWithSessionHeaders(cfg.MarketURL)

	mux.Handle("/kingdom/", http.StripPrefix("/kingdom", kingdom))
	mux.Handle("/economy/", http.StripPrefix("/economy", economy))
	mux.Handle("/market/", http.StripPrefix("/market", market))

	// UI
	mux.Handle("/", http.FileServer(http.Dir("./ui")))

	// wrap middleware(s)
	handler := logRequests(withSession(mux))

	addr := ":" + cfg.Port
	log.Printf("gateway listening on %s", addr)
	log.Printf("proxy: kingdom=%s economy=%s market=%s", cfg.KingdomURL, cfg.EconomyURL, cfg.MarketURL)

	log.Fatal(http.ListenAndServe(addr, handler))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string `json:"userId"`
		KingdomID string `json:"kingdomId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}
	if req.UserID == "" || req.KingdomID == "" {
		httpx.JSON(w, 400, map[string]string{"error": "missing userId or kingdomId"})
		return
	}

	sess := Session{
		UserID:    req.UserID,
		KingdomID: req.KingdomID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := SetSessionCookie(w, sess); err != nil {
		httpx.JSON(w, 500, map[string]string{"error": "failed to set session"})
		return
	}

	httpx.JSON(w, 200, map[string]any{"ok": true})
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	sess, ok := GetSession(r)
	if !ok {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}
	httpx.JSON(w, 200, sess)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	ClearSessionCookie(w)
	httpx.JSON(w, 200, map[string]any{"ok": true})
}

func mustProxyWithSessionHeaders(raw string) *httputil.ReverseProxy {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}

	p := httputil.NewSingleHostReverseProxy(u)
	baseDirector := p.Director

	p.Director = func(r *http.Request) {
		baseDirector(r)

		if sess, ok := GetSession(r); ok {
			r.Header.Set("X-User-Id", sess.UserID)
			r.Header.Set("X-Kingdom-Id", sess.KingdomID)
		} else {
			r.Header.Del("X-User-Id")
			r.Header.Del("X-Kingdom-Id")
		}
	}

	return p
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
