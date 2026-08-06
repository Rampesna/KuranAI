package services

import (
	"strings"
	"testing"

	"github.com/rampesna/kuran-ai-gateway/models"
)

func TestParseScopeLabel(t *testing.T) {
	cases := map[string]string{
		"IN_SCOPE":                 scopeIn,
		"OUT_OF_SCOPE":             scopeOut,
		"  out-of-scope\n":         scopeOut,
		"The answer is IN_SCOPE.":  scopeIn,
		"Label: OUT OF SCOPE":      scopeOut,
		"maybe":                    "",
	}
	for raw, want := range cases {
		if got := parseScopeLabel(raw); got != want {
			t.Fatalf("parseScopeLabel(%q)=%q want %q", raw, got, want)
		}
	}
}

func TestOutOfScopeReplyLanguage(t *testing.T) {
	tr := outOfScopeReply("SQL tarafında raporlama için Power BI mı Qlik mi?")
	if !strings.Contains(tr, "kapsamı dışında") {
		t.Fatalf("expected Turkish refusal, got %q", tr)
	}

	en := outOfScopeReply("Which is better for SQL reporting, Power BI or Qlik?")
	if !strings.Contains(en, "outside Kuran AI") {
		t.Fatalf("expected English refusal, got %q", en)
	}

	ar := outOfScopeReply("ما هو أفضل برنامج تقارير؟")
	if !strings.Contains(ar, "خارج نطاق") {
		t.Fatalf("expected Arabic refusal, got %q", ar)
	}
}

func TestBuildSystemPromptEnforcesHardScope(t *testing.T) {
	prompt := buildSystemPrompt("meal", "")
	for _, required := range []string{
		"HARD SCOPE",
		"Power BI",
		"Do not partially answer off-topic",
		"kapsamı dışında",
	} {
		if !strings.Contains(prompt, required) {
			t.Fatalf("prompt missing %q", required)
		}
	}
}

func TestBuildScopeClassifierMessagesIncludesLatestUser(t *testing.T) {
	history := []models.ChatMessage{
		{Role: "user", Content: "Fatiha ne demek?"},
		{Role: "assistant", Content: "Açılış suresidir."},
	}
	msgs := buildScopeClassifierMessages("Power BI mı Qlik mi?", history)
	if len(msgs) < 2 {
		t.Fatalf("expected classifier messages, got %d", len(msgs))
	}
	last := msgs[len(msgs)-1]
	if last.Role != "user" || last.Content != "Power BI mı Qlik mi?" {
		t.Fatalf("unexpected last classifier message: %+v", last)
	}
	if !strings.Contains(msgs[0].Content, "strict topic classifier") {
		t.Fatalf("missing classifier system prompt")
	}
}
