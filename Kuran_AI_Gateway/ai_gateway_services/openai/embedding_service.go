package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	embeddingModel     = "text-embedding-3-large"
	embeddingBaseURL   = "https://api.openai.com/v1"
	embeddingMaxTokens = 8000
)

// EmbeddingService OpenAI embedding API servisi.
type EmbeddingService struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewEmbeddingService EmbeddingService oluşturur.
func NewEmbeddingService(apiKey string) *EmbeddingService {
	return &EmbeddingService{
		baseURL: embeddingBaseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// EmbeddingResponse OpenAI embedding API yanıtı.
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// EmbeddingRequest embedding isteği.
type EmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// Embed metni embed eder; 3072 boyutlu vektör döner.
func (s *EmbeddingService) Embed(ctx context.Context, text string) ([]float32, int, error) {
	if text == "" {
		return nil, 0, fmt.Errorf("text is empty")
	}
	if s.apiKey == "" {
		return nil, 0, fmt.Errorf("OpenAI API key not configured")
	}

	reqBody := EmbeddingRequest{
		Model: embeddingModel,
		Input: truncateForEmbedding(text, embeddingMaxTokens),
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/embeddings", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, 0, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		return nil, 0, fmt.Errorf("openai api %s: %s", resp.Status, string(buf[:n]))
	}

	var out EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, 0, fmt.Errorf("decode response: %w", err)
	}
	if len(out.Data) == 0 {
		return nil, 0, fmt.Errorf("no embedding in response")
	}

	return out.Data[0].Embedding, out.Usage.TotalTokens, nil
}

// truncateForEmbedding metni max token'a göre kırpar (yaklaşık: 4 char ≈ 1 token).
func truncateForEmbedding(text string, maxTokens int) string {
	runes := []rune(text)
	maxChars := maxTokens * 4
	if len(runes) <= maxChars {
		return text
	}
	return string(runes[:maxChars])
}
