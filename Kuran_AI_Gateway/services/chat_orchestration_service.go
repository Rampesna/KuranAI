package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rampesna/kuran-ai-gateway/ai_gateway_services/openai"
	"github.com/rampesna/kuran-ai-gateway/models"
)

// ChatOrchestrationService chat akışını orkestre eder.
type ChatOrchestrationService struct {
	embedding     *openai.EmbeddingService
	chat          *openai.ChatService
	sessionSvc    *ChatSessionService
	messageSvc    *ChatMessageService
	vectorSearch  *VectorSearchService
}

// NewChatOrchestrationService ChatOrchestrationService oluşturur.
func NewChatOrchestrationService(
	embedding *openai.EmbeddingService,
	chat *openai.ChatService,
	sessionSvc *ChatSessionService,
	messageSvc *ChatMessageService,
	vectorSearch *VectorSearchService,
) *ChatOrchestrationService {
	return &ChatOrchestrationService{
		embedding:    embedding,
		chat:         chat,
		sessionSvc:   sessionSvc,
		messageSvc:   messageSvc,
		vectorSearch: vectorSearch,
	}
}

// ProcessMessage kullanıcı mesajını işler: embed → arama → chat history → LLM stream.
// onSessionReady session+user mesaj oluşturulduktan hemen sonra çağrılır; SSE ile "started" göndermek için.
func (s *ChatOrchestrationService) ProcessMessage(ctx context.Context, authType string, authID uint, sessionID *uint, sourceType, content string, onSessionReady func(uint) error, onChunk func(string) error) (uint, uint, error) {
	if content == "" {
		return 0, 0, fmt.Errorf("content is empty")
	}

	embedding, _, err := s.embedding.Embed(ctx, content)
	if err != nil {
		return 0, 0, fmt.Errorf("embed: %w", err)
	}

	similarDocs, err := s.vectorSearch.SearchBySourceType(ctx, sourceType, embedding, 10)
	if err != nil {
		return 0, 0, fmt.Errorf("vector search: %w", err)
	}

	var session *models.ChatSession
	if sessionID != nil && *sessionID > 0 {
		s, err := s.sessionSvc.GetByID(ctx, *sessionID, authType, authID)
		if err == nil {
			session = s
		}
	}
	if session == nil {
		newSession, err := s.sessionSvc.Create(ctx, authType, authID, content)
		if err != nil {
			return 0, 0, fmt.Errorf("create session: %w", err)
		}
		session = newSession
	}

	messages, _, err := s.messageSvc.ListBySessionID(ctx, session.ID, 0, 50)
	if err != nil {
		return 0, 0, fmt.Errorf("list messages: %w", err)
	}

	chatMessages := buildChatMessages(content, messages, similarDocs, sourceType)

	userMsg, err := s.messageSvc.Create(ctx, session.ID, "user", content, sourceType, embedding)
	if err != nil {
		return 0, 0, fmt.Errorf("create user message: %w", err)
	}
	if onSessionReady != nil {
		if err := onSessionReady(session.ID); err != nil {
			return 0, 0, fmt.Errorf("onSessionReady: %w", err)
		}
	}

	var fullResponse strings.Builder
	err = s.chat.StreamChat(ctx, chatMessages, func(chunk string) error {
		fullResponse.WriteString(chunk)
		if onChunk != nil {
			return onChunk(chunk)
		}
		return nil
	})
	if err != nil {
		return 0, 0, fmt.Errorf("chat stream: %w", err)
	}

	respText := fullResponse.String()
	assistantMsg, err := s.messageSvc.Create(ctx, session.ID, "assistant", respText, sourceType, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("create assistant message: %w", err)
	}

	// Assistant embedding arka planda; cevap hemen kullanıcıya gitti
	go func() {
		bgCtx := context.Background()
		emb, _, err := s.embedding.Embed(bgCtx, respText)
		if err == nil && len(emb) > 0 {
			_ = s.messageSvc.UpdateEmbedding(bgCtx, assistantMsg.ID, emb)
		}
	}()

	return session.ID, userMsg.ID, nil
}

func buildChatMessages(userContent string, history []models.ChatMessage, similarDocs []map[string]interface{}, sourceType string) []openai.ChatMessage {
	ragContext := buildRAGContext(similarDocs, sourceType)
	systemPrompt := buildSystemPrompt(sourceType, ragContext)

	msgs := make([]openai.ChatMessage, 0, len(history)+3)
	msgs = append(msgs, openai.ChatMessage{Role: "system", Content: systemPrompt})
	for _, m := range history {
		msgs = append(msgs, openai.ChatMessage{Role: m.Role, Content: m.Content})
	}
	msgs = append(msgs, openai.ChatMessage{Role: "user", Content: userContent})
	return msgs
}

func buildRAGContext(docs []map[string]interface{}, sourceType string) string {
	if len(docs) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, d := range docs {
		data, ok := d["data"]
		if !ok {
			continue
		}
		text := formatDocData(data)
		if text != "" {
			sb.WriteString(fmt.Sprintf("[%d] %s\n\n", i+1, text))
		}
	}
	return sb.String()
}

func formatDocData(data interface{}) string {
	switch v := data.(type) {
	case string:
		return v
	case []byte:
		var m map[string]interface{}
		if json.Unmarshal(v, &m) == nil {
			return extractTextFromData(m)
		}
		return string(v)
	case map[string]interface{}:
		return extractTextFromData(v)
	default:
		b, _ := json.Marshal(data)
		var m map[string]interface{}
		if json.Unmarshal(b, &m) == nil {
			return extractTextFromData(m)
		}
		return string(b)
	}
}

func extractTextFromData(m map[string]interface{}) string {
	var parts []string
	for _, key := range []string{"meal", "tefsir", "title", "hakkinda", "nuzul", "konu", "fazilet"} {
		if val, ok := m[key]; ok && val != nil {
			if s, ok := val.(string); ok && s != "" {
				parts = append(parts, s)
			}
		}
	}
	if len(parts) > 0 {
		return strings.Join(parts, "\n\n")
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func buildSystemPrompt(sourceType, ragContext string) string {
	base := "Sen Kur'an, tefsir ve hadis konularında yardımcı bir asistansın. Sorulara cevap verirken kaynaklara dayan."
	switch sourceType {
	case "meal":
		base += " Özellikle meal (çeviri) ve bağlam üzerine odaklan."
	case "tefsir":
		base += " Klasik ve modern tefsir kaynaklarından özetler sun."
	case "hadis":
		base += " Sahih hadis referanslarıyla cevap ver."
	}
	if ragContext != "" {
		base += "\n\nAşağıdaki bağlam bilgisini kullanarak soruyu cevapla:\n\n" + ragContext
	}
	return base
}
