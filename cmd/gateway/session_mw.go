package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

type Session struct {
	UserID    string    `json:"userId"`
	KingdomID string    `json:"kingdomId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type sessionKeyType struct{}

var sessionKey = sessionKeyType{}

const sessionCookieName = "civra_session"

var sessionKeyBytes = mustKeyFromB64Env("CIVRA_SESSION_KEY_B64")

func mustKeyFromB64Env(env string) []byte {
	v := os.Getenv(env)
	if v == "" {
		panic(env + " is empty")
	}
	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		panic(env + " invalid base64: " + err.Error())
	}
	if len(b) != 32 {
		panic(env + " must decode to 32 bytes (AES-256-GCM)")
	}
	return b
}

func GetSession(r *http.Request) (Session, bool) {
	s, ok := r.Context().Value(sessionKey).(Session)
	return s, ok
}

func SetSessionCookie(w http.ResponseWriter, sess Session) error {
	if sess.ExpiresAt.IsZero() {
		sess.ExpiresAt = time.Now().Add(24 * time.Hour)
	}
	token, err := encryptSession(sess)
	if err != nil {
		return err
	}

	maxAge := int(time.Until(sess.ExpiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		Expires:  sess.ExpiresAt,
		HttpOnly: true,
		Secure:   os.Getenv("CIVRA_COOKIE_SECURE") != "false",
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		HttpOnly: true,
		Secure:   os.Getenv("CIVRA_COOKIE_SECURE") != "false",
		SameSite: http.SameSiteLaxMode,
	})
}

func withSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookieName)
		if err == nil && c.Value != "" {
			sess, derr := decryptSession(c.Value)
			if derr == nil && time.Now().Before(sess.ExpiresAt) {
				ctx := context.WithValue(r.Context(), sessionKey, sess)
				r = r.WithContext(ctx)
			} else {
				ClearSessionCookie(w)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func encryptSession(sess Session) (string, error) {
	plain, err := json.Marshal(sess)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(sessionKeyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, plain, nil)

	out := append(nonce, ciphertext...)
	return base64.RawURLEncoding.EncodeToString(out), nil
}

func decryptSession(token string) (Session, error) {
	var sess Session

	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return sess, err
	}

	block, err := aes.NewCipher(sessionKeyBytes)
	if err != nil {
		return sess, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return sess, err
	}

	if len(raw) < gcm.NonceSize() {
		return sess, errors.New("token too short")
	}
	nonce := raw[:gcm.NonceSize()]
	ciphertext := raw[gcm.NonceSize():]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return sess, err
	}

	if err := json.Unmarshal(plain, &sess); err != nil {
		return sess, err
	}
	return sess, nil
}
