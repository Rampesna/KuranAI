package google

import (
	"context"
	"net/http"
	"time"

	"github.com/talhacan/gateway/ai_gateway_services/types"
)

const defaultBaseURL = "https://generativelanguage.googleapis.com/v1beta"

// Gateway Google (Gemini) API için gateway; base URL, client ve API key tutar.
type Gateway struct {
	BaseURL string
	ApiKey  string
	Client  *http.Client
}

// NewGateway API key ile Google gateway oluşturur.
func NewGateway(apiKey string) *Gateway {
	return &Gateway{
		BaseURL: defaultBaseURL,
		ApiKey:  apiKey,
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// SendMessage generateContent API ile chat isteği gönderir (messages = history + son user mesajı).
func (g *Gateway) SendMessage(ctx context.Context, model string, messages []types.Message) (string, error) {
	return SendMessage(ctx, g, model, messages)
}
