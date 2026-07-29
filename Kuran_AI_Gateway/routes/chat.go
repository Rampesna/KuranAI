package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/ai_gateway_services/openai"
	"github.com/rampesna/kuran-ai-gateway/config"
	"github.com/rampesna/kuran-ai-gateway/handlers/chat"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// RegisterChatRoutes chat API rotalarını kaydeder.
func RegisterChatRoutes(rg *gin.RouterGroup, tokenService *services.TokenService, openAIConfig config.OpenAIConfig, usageSvc *services.UsageService) {
	whisper := openai.NewGateway(openAIConfig.APIKey)
	transcribeHandler := chat.NewTranscribeHandler(whisper)

	embeddingSvc := openai.NewEmbeddingService(openAIConfig.APIKey)
	chatSvc := openai.NewChatService(openAIConfig.APIKey, openAIConfig.ChatModel, openAIConfig.ReasoningEffort)
	sessionSvc := services.NewChatSessionService()
	messageSvc := services.NewChatMessageService()
	vectorSearch := services.NewVectorSearchService()
	orchestration := services.NewChatOrchestrationService(embeddingSvc, chatSvc, sessionSvc, messageSvc, vectorSearch)

	messageHandler := chat.NewMessageHandler(orchestration, usageSvc)
	sessionsHandler := chat.NewSessionsHandler(sessionSvc)
	messagesHandler := chat.NewMessagesHandler(sessionSvc, messageSvc)

	ttsSvc := openai.NewTTSService(openAIConfig.APIKey)
	ttsHandler := chat.NewTTSHandler(ttsSvc, usageSvc)

	chatProtected := rg.Group("/v1/chat")
	chatProtected.Use(middleware.AuthMiddleware(tokenService), middleware.RequireUser())
	{
		chatProtected.POST("/transcribe", transcribeHandler.Transcribe)
		chatProtected.POST("/message", messageHandler.SendMessage)
		chatProtected.POST("/tts", ttsHandler.Speak)
		chatProtected.GET("/sessions", sessionsHandler.ListSessions)
		chatProtected.GET("/sessions/:id/messages", messagesHandler.ListMessages)
	}
}
