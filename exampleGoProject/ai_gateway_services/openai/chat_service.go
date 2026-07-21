package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/talhacan/gateway/ai_gateway_services/types"
)

// chatRequest OpenAI chat completions request body.
type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatResponse OpenAI chat completions response.
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// SendMessage OpenAI chat completions API çağrısı yapar (messages = history + son user mesajı).
func SendMessage(ctx context.Context, g *Gateway, model string, messages []types.Message) (string, error) {
	msgs := make([]message, 0, len(messages))
	for _, m := range messages {
		msgs = append(msgs, message{Role: m.Role, Content: m.Content})
	}
	if len(msgs) == 0 {
		msgs = []message{{Role: "user", Content: ""}}
	}
	body := chatRequest{
		Model:    model,
		Messages: msgs,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.BaseURL+"/chat/completions", bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.ApiKey)
	resp, err := g.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("openai response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai api: %s %s", resp.Status, string(respBody))
	}
	var out chatResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", fmt.Errorf("openai response: %w", err)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("openai: no choices in response")
	}
	return out.Choices[0].Message.Content, nil
}
