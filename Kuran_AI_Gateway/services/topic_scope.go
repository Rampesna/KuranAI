package services

import (
	"context"
	"log"
	"strings"
	"unicode"

	"github.com/rampesna/kuran-ai-gateway/ai_gateway_services/openai"
	"github.com/rampesna/kuran-ai-gateway/models"
)

const (
	scopeIn  = "IN_SCOPE"
	scopeOut = "OUT_OF_SCOPE"

	scopeClassifierEffort = "low"

	scopeClassifierPrompt = `You are a strict topic classifier for Kuran AI.

Decide if the user's LATEST message is within product scope.

IN_SCOPE:
- Quran text, translation/meal, tafsir/exegesis, hadith
- Islamic belief, worship, ethics, prophets, companions, Arabic religious terms
- Greetings, thanks, and asking what the assistant can do
- Clarifying follow-ups about a previous in-scope religious answer

OUT_OF_SCOPE:
- Programming, SQL, databases, BI tools (Power BI, Qlik), software, IT, engineering
- Secular school/homework, sports, entertainment, politics, finance/trading
- General worldly advice unrelated to Quran/tafsir/hadith
- Medical/legal professional advice framed as secular expertise

Reply with exactly one token: IN_SCOPE or OUT_OF_SCOPE`
)

func parseScopeLabel(raw string) string {
	normalized := strings.ToUpper(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")
	switch {
	case strings.Contains(normalized, scopeOut):
		return scopeOut
	case strings.Contains(normalized, scopeIn):
		return scopeIn
	default:
		return ""
	}
}

func detectReplyLang(content string) string {
	var arabic, latin int
	for _, r := range content {
		switch {
		case unicode.In(r, unicode.Arabic):
			arabic++
		case unicode.IsLetter(r):
			latin++
		}
	}
	if arabic > latin && arabic > 0 {
		return "ar"
	}

	lower := strings.ToLower(content)
	turkishHints := []string{" mı", " mi", " mu", " mü", "için", "nedir", "nasıl", "mı?", "mi?", "mı ", "mi ", "ş", "ğ", "ı", "ü", "ö", "ç"}
	for _, h := range turkishHints {
		if strings.Contains(lower, h) {
			return "tr"
		}
	}

	englishHints := []string{" the ", " what ", " how ", " which ", " is ", " are ", " can ", " please "}
	padded := " " + lower + " "
	for _, h := range englishHints {
		if strings.Contains(padded, h) {
			return "en"
		}
	}
	return "tr"
}

func outOfScopeReply(userContent string) string {
	switch detectReplyLang(userContent) {
	case "en":
		return "This question is outside Kuran AI's scope. I only help with Quran translation, tafsir, and hadith. Please ask a religious question."
	case "ar":
		return "هذا السؤال خارج نطاق كران AI. يمكنني المساعدة فقط في ترجمة القرآن والتفسير والحديث. يرجى طرح سؤال ديني."
	default:
		return "Bu soru Kuran AI'nin kapsamı dışında. Ben yalnızca Kur'an meali, tefsir ve hadis konularında yardımcı olabilirim. Lütfen dini bir soru sorun."
	}
}

func buildScopeClassifierMessages(userContent string, history []models.ChatMessage) []openai.ChatMessage {
	msgs := make([]openai.ChatMessage, 0, 6)
	msgs = append(msgs, openai.ChatMessage{Role: "system", Content: scopeClassifierPrompt})

	start := 0
	if len(history) > 4 {
		start = len(history) - 4
	}
	for _, m := range history[start:] {
		role := m.Role
		if role != "user" && role != "assistant" {
			continue
		}
		msgs = append(msgs, openai.ChatMessage{Role: role, Content: m.Content})
	}
	msgs = append(msgs, openai.ChatMessage{Role: "user", Content: userContent})
	return msgs
}

// classifyTopicScope returns whether the latest user message is in product scope.
// On classifier failure it returns true (fail-open) so legitimate religious questions are not blocked.
func (s *ChatOrchestrationService) classifyTopicScope(ctx context.Context, userContent string, history []models.ChatMessage) bool {
	raw, err := s.chat.CompleteChat(ctx, buildScopeClassifierMessages(userContent, history), scopeClassifierEffort)
	if err != nil {
		log.Printf("[ChatOrchestration] scope classifier failed, allowing request: %v", err)
		return true
	}
	label := parseScopeLabel(raw)
	if label == "" {
		log.Printf("[ChatOrchestration] scope classifier ambiguous reply %q, allowing request", raw)
		return true
	}
	return label == scopeIn
}
