package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// TTSService OpenAI Text-to-Speech servisi.
type TTSService struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewTTSService TTSService oluşturur.
func NewTTSService(apiKey string) *TTSService {
	return &TTSService{
		baseURL: "https://api.openai.com/v1",
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

// TTSSpeechRequest OpenAI TTS isteği.
// Voice: alloy, echo, fable, onyx, nova, shimmer
// Model: tts-1 (hızlı), tts-1-hd (daha kaliteli)
type TTSSpeechRequest struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	Voice          string `json:"voice"`
	ResponseFormat string `json:"response_format,omitempty"` // mp3, opus, aac, flac
}

// Speech metni sese çevirir, MP3 bytes döner.
func (s *TTSService) Speech(ctx context.Context, text, voice string) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}
	if text == "" {
		return nil, fmt.Errorf("input text is empty")
	}
	if voice == "" {
		voice = "ash"
	}

	reqBody := TTSSpeechRequest{
		Model:          "tts-1-hd",
		Input:          text,
		Voice:          voice,
		ResponseFormat: "mp3",
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/audio/speech", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai tts request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf, _ := io.ReadAll(resp.Body)
		log.Printf("[TTS] OpenAI API error: %s %s", resp.Status, string(buf))
		return nil, fmt.Errorf("openai tts api %s: %s", resp.Status, string(buf))
	}

	return io.ReadAll(resp.Body)
}
