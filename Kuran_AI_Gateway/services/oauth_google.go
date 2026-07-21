package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rampesna/kuran-ai-gateway/config"
)

// GoogleIDTokenPayload Google ID token payload (JWT claims).
type GoogleIDTokenPayload struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// Subject returns the Google user ID (sub claim).
func (p *GoogleIDTokenPayload) Subject() string { return p.Sub }

// ValidateGoogleIDToken Google ID token doğrular ve payload döner (tokeninfo API; secret gerekmez).
func ValidateGoogleIDToken(ctx context.Context, idToken string) (*GoogleIDTokenPayload, error) {
	cfg := config.Load()
	if cfg.OAuth.GoogleClientID == "" {
		return nil, fmt.Errorf("Google OAuth not configured")
	}
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token validation failed: %d", resp.StatusCode)
	}
	var payload struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Aud           string `json:"aud"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Aud != cfg.OAuth.GoogleClientID {
		return nil, fmt.Errorf("invalid audience")
	}
	return &GoogleIDTokenPayload{
		Sub:           payload.Sub,
		Email:         payload.Email,
		EmailVerified: payload.EmailVerified == "true",
		Name:          payload.Name,
		Picture:       payload.Picture,
	}, nil
}
