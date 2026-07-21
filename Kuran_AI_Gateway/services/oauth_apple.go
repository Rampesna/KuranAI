package services

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rampesna/kuran-ai-gateway/config"
)

// AppleIDTokenPayload Apple identity token payload.
type AppleIDTokenPayload struct {
	Sub            string `json:"sub"`
	Email          string `json:"email"`
	EmailVerified  string `json:"email_verified"`
	IsPrivateEmail string `json:"is_private_email"`
}

// ValidateAppleIDToken Apple identity token doğrular.
// Apple JWT'yi kendimiz parse edip doğruluyoruz (public key ile).
func ValidateAppleIDToken(ctx context.Context, identityToken string) (*AppleIDTokenPayload, error) {
	cfg := config.Load()
	if len(cfg.OAuth.AppleClientIDs) == 0 {
		return nil, fmt.Errorf("Apple OAuth not configured")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(identityToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	kid, _ := token.Header["kid"].(string)
	if kid == "" {
		return nil, fmt.Errorf("missing kid in token header")
	}

	// Apple public keys
	keysURL := "https://appleid.apple.com/auth/keys"
	req, _ := http.NewRequestWithContext(ctx, "GET", keysURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var keys struct {
		Keys []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			Use string `json:"use"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, err
	}
	var pubKey *rsa.PublicKey
	for _, k := range keys.Keys {
		if k.Kid == kid {
			nBytes, _ := base64.RawURLEncoding.DecodeString(k.N)
			eBytes, _ := base64.RawURLEncoding.DecodeString(k.E)
			n := new(big.Int).SetBytes(nBytes)
			e := int(binary.BigEndian.Uint32(append([]byte{0}, eBytes...)))
			pubKey = &rsa.PublicKey{N: n, E: e}
			break
		}
	}
	if pubKey == nil {
		return nil, fmt.Errorf("public key not found")
	}

	// Verify token
	parsed, err := jwt.ParseWithClaims(identityToken, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("unexpected signing algorithm")
		}
		return pubKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	if err != nil || !parsed.Valid {
		return nil, fmt.Errorf("invalid token signature")
	}
	c := parsed.Claims.(jwt.MapClaims)

	// Audience kontrolü (Apple client_id veya bundle id)
	audiences, err := c.GetAudience()
	if err != nil {
		return nil, fmt.Errorf("invalid audience")
	}
	iss, err := c.GetIssuer()
	if err != nil {
		return nil, fmt.Errorf("invalid issuer")
	}
	if iss != "https://appleid.apple.com" {
		return nil, fmt.Errorf("invalid issuer")
	}
	audienceAllowed := false
	for _, audience := range audiences {
		if containsString(cfg.OAuth.AppleClientIDs, audience) {
			audienceAllowed = true
			break
		}
	}
	if !audienceAllowed {
		return nil, fmt.Errorf("invalid audience")
	}

	exp, err := c.GetExpirationTime()
	if err != nil || exp == nil || time.Now().After(exp.Time) {
		return nil, fmt.Errorf("token expired")
	}

	sub, err := c.GetSubject()
	if err != nil || sub == "" {
		return nil, fmt.Errorf("missing sub")
	}

	email, _ := c["email"].(string)
	emailVerified, _ := c["email_verified"].(string)
	isPrivate, _ := c["is_private_email"].(string)

	return &AppleIDTokenPayload{
		Sub:            sub,
		Email:          email,
		EmailVerified:  emailVerified,
		IsPrivateEmail: strings.ToLower(isPrivate),
	}, nil
}
