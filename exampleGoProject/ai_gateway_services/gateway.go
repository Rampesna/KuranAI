package ai_gateway_services

import (
	"context"
	"fmt"
	"strings"

	"github.com/talhacan/gateway/ai_gateway_services/anthropic"
	"github.com/talhacan/gateway/ai_gateway_services/google"
	"github.com/talhacan/gateway/ai_gateway_services/moonshot"
	"github.com/talhacan/gateway/ai_gateway_services/openai"
	"github.com/talhacan/gateway/ai_gateway_services/types"
)

// ChatMessage tek mesaj (role + content); history için kullanılır.
type ChatMessage struct {
	Role    string `json:"role"`    // user | assistant
	Content string `json:"content"`
}

// ChatResponse chat cevabı (tüm sağlayıcılar için ortak).
type ChatResponse struct {
	Content string `json:"content"`
}

// ProviderGateway sağlayıcıya özel gateway; SendMessage ile chat çağrısı yapar (messages = history + son user mesajı).
type ProviderGateway interface {
	SendMessage(ctx context.Context, model string, messages []types.Message) (content string, err error)
}

// SendMessage provider'a göre ilgili gateway'i oluşturup chat isteği gönderir. messages boşsa tek user mesajı kabul edilir (prompt).
func SendMessage(ctx context.Context, provider, prompt, model, apiKey string, messages []types.Message) (*ChatResponse, error) {
	provider = strings.ToLower(strings.TrimSpace(provider))
	var gw ProviderGateway
	switch provider {
	case "openai":
		gw = openai.NewGateway(apiKey)
	case "anthropic":
		gw = anthropic.NewGateway(apiKey)
	case "google":
		gw = google.NewGateway(apiKey)
	case "moonshot":
		gw = moonshot.NewGateway(apiKey)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	msgs := messages
	if len(msgs) == 0 {
		msgs = []types.Message{{Role: "user", Content: prompt}}
	}
	content, err := gw.SendMessage(ctx, model, msgs)
	if err != nil {
		return nil, err
	}
	return &ChatResponse{Content: content}, nil
}
