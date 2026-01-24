package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jose "github.com/go-jose/go-jose/v4"
)

type Claims struct {
	UserID    string    `json:"userId"`
	KingdomID string    `json:"kingdomId"`
	Role      string    `json:"role,omitempty"`
	Exp       time.Time `json:"exp"`
}

var ErrExpired = errors.New("token expired")

type JWE struct {
	enc jose.Encrypter
	key []byte
}

func NewJWE(secret []byte) (*JWE, error) {
	if len(secret) != 32 {
		return nil, fmt.Errorf("JWE secret must be 32 bytes for A256GCM, got %d", len(secret))
	}
	jwk := jose.JSONWebKey{
		Key:       secret,
		Algorithm: string(jose.DIRECT),
		Use:       "enc",
	}

	enc, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.DIRECT,
			Key:       jwk,
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &JWE{enc: enc, key: secret}, nil
}

func (j *JWE) Encrypt(c Claims) (string, error) {
	raw, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	obj, err := j.enc.Encrypt(raw)
	if err != nil {
		return "", err
	}

	return obj.CompactSerialize()
}

func (j *JWE) Decrypt(token string) (Claims, error) {
	obj, err := jose.ParseEncrypted(
		token,
		[]jose.KeyAlgorithm{jose.DIRECT},
		[]jose.ContentEncryption{jose.A256GCM},
	)
	if err != nil {
		return Claims{}, err
	}

	raw, err := obj.Decrypt(j.key)
	if err != nil {
		return Claims{}, err
	}

	var c Claims
	if err := json.Unmarshal(raw, &c); err != nil {
		return Claims{}, err
	}

	if !c.Exp.IsZero() && time.Now().After(c.Exp) {
		return Claims{}, ErrExpired
	}

	return c, nil
}
