package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/ai_gateway_services/openai"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	"github.com/rampesna/kuran-ai-gateway/requests/chat"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// TTSHandler OpenAI TTS handler'ı.
type TTSHandler struct {
	tts      *openai.TTSService
	usageSvc *services.UsageService
}

// NewTTSHandler TTSHandler oluşturur.
func NewTTSHandler(tts *openai.TTSService, usageSvc *services.UsageService) *TTSHandler {
	return &TTSHandler{tts: tts, usageSvc: usageSvc}
}

// Speak metni sese çevirir, MP3 döner.
func (h *TTSHandler) Speak(c *gin.Context) {
	userIDVal, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	if err := h.usageSvc.CheckAndLogTTS(c.Request.Context(), userID); err != nil {
		if le, ok := services.IsLimitExceededError(err); ok {
			c.JSON(http.StatusPaymentRequired, core.NewServiceResponse(false, le.Message, 402, gin.H{
				"code":      le.Code,
				"limitType": le.LimitType,
				"used":      le.Used,
				"limit":     le.Limit,
				"packageId": le.PackageID,
			}))
			return
		}
		c.JSON(http.StatusInternalServerError, core.NewServiceResponse(false, "Kullanım kontrolü başarısız", 500, nil))
		return
	}

	var req chat.TTSRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
		return
	}

	// Max 4096 karakter (OpenAI TTS limiti)
	if len(req.Text) > 4096 {
		req.Text = req.Text[:4096]
	}

	voice := req.Voice
	if voice == "" {
		voice = "ash"
	}

	audioData, err := h.tts.Speech(c.Request.Context(), req.Text, voice)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "TTS failed: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, "audio/mpeg", audioData)
}
