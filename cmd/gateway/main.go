package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"civra-core/pkg/auth"
	"civra-core/pkg/config"
	"civra-core/pkg/httpx"
)

type LoginReq struct {
	UserID    string `json:"userId"`
	KingdomID string `json:"kingdomId"`
	Role      string `json:"role,omitempty"`
}

func main() {
	cfg := config.LoadGateway()
	jwe, err := auth.NewJWE([]byte(cfg.AuthSecret))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, 200, map[string]any{"service": "gateway", "ok": true})
	})

	// auth
	mux.HandleFunc("/auth/login", handleLogin(jwe))
	mux.HandleFunc("/auth/me", handleMe(jwe))
	mux.HandleFunc("/auth/logout", handleLogout)

	// proxies
	kingdom := mustProxy(cfg.KingdomURL)
	economy := mustProxy(cfg.EconomyURL)
	market := mustProxy(cfg.MarketURL)

	mux.Handle("/kingdom/", http.StripPrefix("/kingdom", kingdom))
	mux.Handle("/economy/", http.StripPrefix("/economy", economy))
	mux.Handle("/market/", http.StripPrefix("/market", market))

	// UI
	mux.Handle("/", http.FileServer(http.Dir("./ui")))

	// wrap middleware(s)
	handler := logRequests(withAuth(jwe, mux))

	addr := ":" + cfg.Port
	log.Printf("gateway listening on %s", addr)
	log.Printf("proxy: kingdom=%s economy=%s market=%s", cfg.KingdomURL, cfg.EconomyURL, cfg.MarketURL)

	log.Fatal(http.ListenAndServe(addr, handler))
}

func handleLogin(jwe *auth.JWE) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httpx.JSON(w, 405, map[string]string{"error": "method not allowed"})
			return
		}

		var req LoginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.JSON(w, 400, map[string]string{"error": "invalid json"})
			return
		}
		if req.UserID == "" || req.KingdomID == "" {
			httpx.JSON(w, 400, map[string]string{"error": "missing userId or kingdomId"})
			return
		}

		claims := auth.Claims{
			UserID:    req.UserID,
			KingdomID: req.KingdomID,
			Role:      req.Role,
			Exp:       time.Now().Add(24 * time.Hour),
		}

		token, err := jwe.Encrypt(claims)
		if err != nil {
			httpx.JSON(w, 500, map[string]string{"error": "token encrypt failed"})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "civra_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Expires:  claims.Exp,
		})

		httpx.JSON(w, 200, map[string]any{"ok": true})
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "civra_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	httpx.JSON(w, 200, map[string]any{"ok": true})
}

func handleMe(jwe *auth.JWE) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("civra_token")
		if err != nil || c.Value == "" {
			httpx.JSON(w, 401, map[string]string{"error": "unauthorized"})
			return
		}
		claims, err := jwe.Decrypt(c.Value)
		if err != nil {
			httpx.JSON(w, 401, map[string]string{"error": "unauthorized"})
			return
		}
		httpx.JSON(w, 200, claims)
	}
}

func withAuth(jwe *auth.JWE, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// public
		if path == "/" ||
			strings.HasPrefix(path, "/auth/") ||
			strings.HasPrefix(path, "/style.css") ||
			strings.HasPrefix(path, "/app.js") ||
			strings.HasPrefix(path, "/favicon") {
			next.ServeHTTP(w, r)
			return
		}

		// protect only API routes (proxy)
		if !(strings.HasPrefix(path, "/economy/") || strings.HasPrefix(path, "/market/") || strings.HasPrefix(path, "/kingdom/")) {
			next.ServeHTTP(w, r)
			return
		}

		c, err := r.Cookie("civra_token")
		if err != nil || c.Value == "" {
			httpx.JSON(w, 401, map[string]string{"error": "unauthorized"})
			return
		}

		claims, err := jwe.Decrypt(c.Value)
		if err != nil {
			httpx.JSON(w, 401, map[string]string{"error": "unauthorized"})
			return
		}

		r.Header.Set("X-User-Id", claims.UserID)
		r.Header.Set("X-Kingdom-Id", claims.KingdomID)
		if claims.Role != "" {
			r.Header.Set("X-Role", claims.Role)
		}

		next.ServeHTTP(w, r)
	})
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
		if !strings.HasPrefix(r.URL.Path, "/healthz") {
			log.Printf("%s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}
