package moonshot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/talhacan/gateway/ai_gateway_services/types"
)

// chatRequest OpenAI uyumlu chat completions request (Moonshot aynı formatı kullanır).
type chatRequest struct {
	Model    string           `json:"model"`
	Messages []types.Message `json:"messages"`
}

// chatResponse OpenAI uyumlu response.
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// SendMessage Moonshot chat completions API çağrısı yapar (messages = history + son user mesajı).
func SendMessage(ctx context.Context, g *Gateway, model string, messages []types.Message) (string, error) {
	if len(messages) == 0 {
		messages = []types.Message{{Role: "user", Content: ""}}
	}
	body := chatRequest{
		Model:    model,
		Messages: messages,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("moonshot request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.BaseURL+"/chat/completions", bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("moonshot request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.ApiKey)
	resp, err := g.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("moonshot request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("moonshot response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("moonshot api: %s %s", resp.Status, string(respBody))
	}
	var out chatResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", fmt.Errorf("moonshot response: %w", err)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("moonshot: no choices in response")
	}
	return out.Choices[0].Message.Content, nil
}
