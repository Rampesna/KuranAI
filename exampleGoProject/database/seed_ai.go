package database

import (
	"log"

	"github.com/talhacan/gateway/models"
	"gorm.io/gorm"
)

// runSeedAI 29 Ocak 2026 itibarıyla OpenAI, Google, Anthropic, Moonshot sağlayıcılarını ve API modellerini seed eder (idempotent).
func runSeedAI() error {
	ctx := DB.Session(&gorm.Session{})
	providers := []struct {
		name string
		slug string
	}{
		{"OpenAI", "openai"},
		{"Google", "google"},
		{"Anthropic", "anthropic"},
		{"Moonshot", "moonshot"},
	}
	providerIDs := make(map[string]uint)
	for _, p := range providers {
		var prov models.AiProvider
		err := ctx.Where("slug = ?", p.slug).FirstOrCreate(&prov, models.AiProvider{Name: p.name, Slug: p.slug}).Error
		if err != nil {
			return err
		}
		providerIDs[p.slug] = prov.ID
	}

	// OpenAI modelleri (29 Ocak 2026)
	openaiModels := []struct{ name, apiID string }{
		{"GPT-5.2", "gpt-5.2"},
		{"GPT-5.2 Pro", "gpt-5.2-pro"},
		{"GPT-5 Mini", "gpt-5-mini"},
		{"GPT-4.1", "gpt-4.1"},
		{"GPT-4.1 Mini", "gpt-4.1-mini"},
		{"GPT-4.1 Nano", "gpt-4.1-nano"},
		{"o4-mini", "o4-mini"},
		{"GPT Realtime", "gpt-realtime"},
		{"GPT Realtime Mini", "gpt-realtime-mini"},
		{"Sora 2", "sora-2"},
		{"Sora 2 Pro", "sora-2-pro"},
	}
	for _, m := range openaiModels {
		if err := firstOrCreateModel(ctx, providerIDs["openai"], m.name, m.apiID); err != nil {
			return err
		}
	}

	// Google (Gemini) modelleri (29 Ocak 2026)
	googleModels := []struct{ name, apiID string }{
		{"Gemini 3 Pro Preview", "gemini-3-pro-preview"},
		{"Gemini 3 Pro Image Preview", "gemini-3-pro-image-preview"},
		{"Gemini 3 Flash Preview", "gemini-3-flash-preview"},
		{"Gemini 2.5 Pro", "gemini-2.5-pro"},
		{"Gemini 2.5 Flash", "gemini-2.5-flash"},
		{"Gemini 2.5 Flash-Lite", "gemini-2.5-flash-lite"},
	}
	for _, m := range googleModels {
		if err := firstOrCreateModel(ctx, providerIDs["google"], m.name, m.apiID); err != nil {
			return err
		}
	}

	// Anthropic (Claude) modelleri (29 Ocak 2026)
	anthropicModels := []struct{ name, apiID string }{
		{"Claude 3.5 Sonnet", "claude-3-5-sonnet-20241022"},
		{"Claude 3 Opus", "claude-3-opus-20240229"},
		{"Claude 3.5 Haiku", "claude-3-5-haiku-20241022"},
		{"Claude 3 Haiku", "claude-3-haiku-20240307"},
		{"Claude Sonnet 4", "claude-sonnet-4-20250514"},
		{"Claude Opus 4", "claude-opus-4-20250514"},
	}
	for _, m := range anthropicModels {
		if err := firstOrCreateModel(ctx, providerIDs["anthropic"], m.name, m.apiID); err != nil {
			return err
		}
	}

	// Moonshot (Kimi) modelleri (29 Ocak 2026)
	moonshotModels := []struct{ name, apiID string }{
		{"Kimi K2.5", "kimi-k2.5"},
		{"Kimi K2 Thinking", "kimi-k2-thinking"},
		{"Kimi K2", "kimi-k2"},
	}
	for _, m := range moonshotModels {
		if err := firstOrCreateModel(ctx, providerIDs["moonshot"], m.name, m.apiID); err != nil {
			return err
		}
	}

	log.Println("seed AI: providers and models ready")
	return nil
}

func firstOrCreateModel(ctx *gorm.DB, providerID uint, name, apiID string) error {
	var row models.AiModel
	return ctx.Where("provider_id = ? AND api_id = ?", providerID, apiID).
		FirstOrCreate(&row, models.AiModel{ProviderID: providerID, Name: name, ApiID: apiID}).Error
}
