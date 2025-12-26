package main

import (
	"encoding/base64"
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

type Session struct {
	UserID    string `json:"userId"`
	KingdomID string `json:"kingdomId"`
}

func main() {
	cfg := config.LoadGateway()

	kingdom := mustProxy(cfg.KingdomURL)
	economy := mustProxy(cfg.EconomyURL)
	market := mustProxy(cfg.MarketURL)

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
	var s Session
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		httpx.JSON(w, 400, map[string]string{"error": "invalid json"})
		return
	}
	if s.UserID == "" || s.KingdomID == "" {
		httpx.JSON(w, 400, map[string]string{"error": "missing userId or kingdomId"})
		return
	}

	raw, _ := json.Marshal(s)
	val := base64.StdEncoding.EncodeToString(raw)

	http.SetCookie(w, &http.Cookie{
		Name:     "civra_session",
		Value:    val,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	httpx.JSON(w, 200, map[string]any{"ok": true})
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("civra_session")
	if err != nil || c.Value == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return
	}

	raw, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		httpx.JSON(w, 401, map[string]string{"error": "bad session"})
		return
	}

	var s Session
	if err := json.Unmarshal(raw, &s); err != nil || s.UserID == "" || s.KingdomID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "bad session"})
		return
	}

	httpx.JSON(w, 200, s)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "civra_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	httpx.JSON(w, 200, map[string]any{"ok": true})
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
