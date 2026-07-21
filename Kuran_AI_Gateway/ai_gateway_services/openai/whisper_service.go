package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	defaultBaseURL  = "https://api.openai.com/v1"
	defaultModel    = "whisper-1"
	defaultLanguage = ""
)

// Gateway OpenAI API gateway (Whisper için).
type Gateway struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// NewGateway API key ile OpenAI gateway oluşturur.
func NewGateway(apiKey string) *Gateway {
	return &Gateway{
		BaseURL: defaultBaseURL,
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// TranscriptionResponse Whisper API transcription yanıtı (response_format=json).
type TranscriptionResponse struct {
	Text string `json:"text"`
}

// Transcribe ses dosyasını metne çevirir.
// audioData: ham ses bytes (mp3, m4a, wav, webm, ogg, flac vb.)
// filename: dosya adı (örn. "recording.m4a") - extension'dan format çıkarılır.
func (g *Gateway) Transcribe(ctx context.Context, audioData []byte, filename string) (string, error) {
	if len(audioData) == 0 {
		return "", fmt.Errorf("audio data is empty")
	}
	if g.APIKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// file field
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write(audioData); err != nil {
		return "", fmt.Errorf("write audio: %w", err)
	}

	// model
	if err := w.WriteField("model", defaultModel); err != nil {
		return "", fmt.Errorf("write model: %w", err)
	}
	// response_format: text veya json - json daha yapılandırılmış
	if err := w.WriteField("response_format", "json"); err != nil {
		return "", fmt.Errorf("write response_format: %w", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("close multipart: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.BaseURL+"/audio/transcriptions", &buf)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.APIKey)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := g.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai api %s: %s", resp.Status, string(body))
	}

	var out TranscriptionResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return out.Text, nil
}
