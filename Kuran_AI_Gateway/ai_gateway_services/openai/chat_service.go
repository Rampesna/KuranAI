package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	defaultChatModel           = "gpt-5.6-terra"
	defaultChatReasoningEffort = "medium"
	chatBaseURL                = "https://api.openai.com/v1"
)

// ChatService OpenAI chat completions servisi.
type ChatService struct {
	baseURL         string
	apiKey          string
	model           string
	reasoningEffort string
	client          *http.Client
}

// NewChatService ChatService oluşturur.
func NewChatService(apiKey, model, reasoningEffort string) *ChatService {
	model = strings.TrimSpace(model)
	if model == "" {
		model = defaultChatModel
	}
	reasoningEffort = strings.TrimSpace(reasoningEffort)
	if reasoningEffort == "" {
		reasoningEffort = defaultChatReasoningEffort
	}
	service := &ChatService{
		baseURL:         chatBaseURL,
		apiKey:          apiKey,
		model:           model,
		reasoningEffort: reasoningEffort,
		client:          &http.Client{Timeout: 120 * time.Second},
	}
	log.Printf("[Chat] OpenAI configured model=%s reasoning_effort=%s", service.model, service.reasoningEffort)
	return service
}

// ChatMessage chat mesajı.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest chat isteği.
type ChatRequest struct {
	Model           string        `json:"model"`
	Messages        []ChatMessage `json:"messages"`
	Stream          bool          `json:"stream"`
	ReasoningEffort string        `json:"reasoning_effort,omitempty"`
}

// ChatStreamDelta SSE stream delta.
type ChatStreamDelta struct {
	Content string `json:"content,omitempty"`
}

// ChatStreamChoice SSE stream choice.
type ChatStreamChoice struct {
	Delta        ChatStreamDelta `json:"delta"`
	FinishReason string          `json:"finish_reason,omitempty"`
}

// ChatStreamChunk SSE stream chunk.
type ChatStreamChunk struct {
	Choices []ChatStreamChoice `json:"choices,omitempty"`
}

// StreamChat chat completions stream eder; her chunk için callback çağrılır.
func (s *ChatService) StreamChat(ctx context.Context, messages []ChatMessage, onChunk func(string) error) error {
	if s.apiKey == "" {
		return fmt.Errorf("OpenAI API key not configured")
	}

	reqBody := ChatRequest{
		Model:           s.model,
		Messages:        messages,
		Stream:          true,
		ReasoningEffort: s.reasoningEffort,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf, _ := io.ReadAll(resp.Body)
		log.Printf("[Chat] OpenAI API error: %s %s", resp.Status, string(buf))
		return fmt.Errorf("openai api %s: %s", resp.Status, string(buf))
	}

	scanner := bufio.NewScanner(resp.Body)
	const maxLine = 512 * 1024
	scanner.Buffer(make([]byte, maxLine), maxLine)
	chunkCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk ChatStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			log.Printf("[Chat] stream parse skip (line len=%d): %v", len(data), err)
			continue
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		content := chunk.Choices[0].Delta.Content
		if content != "" && onChunk != nil {
			chunkCount++
			if err := onChunk(content); err != nil {
				return err
			}
		}
	}
	if chunkCount == 0 {
		log.Printf("[Chat] OpenAI stream completed with 0 content chunks")
	}
	return scanner.Err()
}
