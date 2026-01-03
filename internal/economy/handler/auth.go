package handler

import (
	"net/http"

	"civra-core/pkg/httpx"
)

func getIdentity(w http.ResponseWriter, r *http.Request) (userID, kingdomID string, ok bool) {
	userID = r.Header.Get("X-User-Id")
	kingdomID = r.Header.Get("X-Kingdom-Id")
	if userID == "" || kingdomID == "" {
		httpx.JSON(w, 401, map[string]string{"error": "no session"})
		return "", "", false
	}
	return userID, kingdomID, true
}
