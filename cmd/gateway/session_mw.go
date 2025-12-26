package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type sessionKeyType struct{}

var sessionKey = sessionKeyType{}

func withSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("civra_session")
		if err == nil {
			raw, err := base64.StdEncoding.DecodeString(c.Value)
			if err == nil {
				var s Session
				if json.Unmarshal(raw, &s) == nil {
					ctx := context.WithValue(r.Context(), sessionKey, s)

					// 👇 прокидываем дальше в микросервисы
					r = r.WithContext(ctx)
					r.Header.Set("X-User-Id", s.UserID)
					r.Header.Set("X-Kingdom-Id", s.KingdomID)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
