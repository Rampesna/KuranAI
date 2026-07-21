package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/talhacan/gateway/ai_gateway_services/types"
)

// chatRequest Anthropic messages API request.
type chatRequest struct {
	Model     string         `json:"model"`
	MaxTokens int            `json:"max_tokens"`
	Messages  []types.Message `json:"messages"`
}

// chatResponse Anthropic messages API response.
type chatResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

// SendMessage Anthropic messages API çağrısı yapar (messages = history + son user mesajı).
func SendMessage(ctx context.Context, g *Gateway, model string, messages []types.Message) (string, error) {
	if len(messages) == 0 {
		messages = []types.Message{{Role: "user", Content: ""}}
	}
	body := chatRequest{
		Model:     model,
		MaxTokens: 4096,
		Messages:  messages,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("anthropic request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.BaseURL+"/messages", bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("anthropic request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", g.ApiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	resp, err := g.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("anthropic response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic api: %s %s", resp.Status, string(respBody))
	}
	var out chatResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", fmt.Errorf("anthropic response: %w", err)
	}
	if len(out.Content) == 0 {
		return "", fmt.Errorf("anthropic: empty content")
	}
	return out.Content[0].Text, nil
}
