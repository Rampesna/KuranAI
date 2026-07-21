package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/talhacan/gateway/ai_gateway_services/types"
)

// generateRequest Gemini generateContent request (multi-turn: contents = history).
type generateRequest struct {
	Contents []contentPart `json:"contents"`
}

type contentPart struct {
	Role  string `json:"role,omitempty"` // user | model
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

// generateResponse Gemini generateContent response.
type generateResponse struct {
	Candidates []struct {
		Content struct {
			Parts []part `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// SendMessage Google Gemini generateContent API çağrısı yapar (messages = history + son user mesajı).
func SendMessage(ctx context.Context, g *Gateway, model string, messages []types.Message) (string, error) {
	contents := make([]contentPart, 0, len(messages))
	for _, m := range messages {
		role := "user"
		if m.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, contentPart{Role: role, Parts: []part{{Text: m.Content}}})
	}
	if len(contents) == 0 {
		contents = []contentPart{{Parts: []part{{Text: ""}}}}
	}
	body := generateRequest{Contents: contents}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("google request: %w", err)
	}
	path := g.BaseURL + "/models/" + url.PathEscape(model) + ":generateContent?key=" + url.QueryEscape(g.ApiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("google request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("google request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("google response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google api: %s %s", resp.Status, string(respBody))
	}
	var out generateResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", fmt.Errorf("google response: %w", err)
	}
	if len(out.Candidates) == 0 || len(out.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("google: no content in response")
	}
	return out.Candidates[0].Content.Parts[0].Text, nil
}
