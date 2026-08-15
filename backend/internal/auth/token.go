package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Claims struct {
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
}

func Issue(email, secret string) (string, error) {
	payload, err := json.Marshal(Claims{Email: email, ExpiresAt: time.Now().Add(24 * time.Hour).Unix()})
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	return encoded + "." + signature(encoded, secret), nil
}
func Verify(token, secret string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 || !hmac.Equal([]byte(signature(parts[0], secret)), []byte(parts[1])) {
		return Claims{}, errors.New("invalid token")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, errors.New("invalid token")
	}
	var claims Claims
	if json.Unmarshal(raw, &claims) != nil || claims.ExpiresAt < time.Now().Unix() {
		return Claims{}, errors.New("expired token")
	}
	return claims, nil
}
func signature(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
