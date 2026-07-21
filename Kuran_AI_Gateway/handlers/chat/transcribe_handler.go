package chat

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/ai_gateway_services/openai"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
)

// TranscribeHandler ses kaydını metne çevirme handler'ı.
type TranscribeHandler struct {
	whisper *openai.Gateway
}

// NewTranscribeHandler TranscribeHandler oluşturur.
func NewTranscribeHandler(whisper *openai.Gateway) *TranscribeHandler {
	return &TranscribeHandler{whisper: whisper}
}

// Transcribe multipart form ile gelen ses dosyasını Whisper API'ye gönderip metin döner.
func (h *TranscribeHandler) Transcribe(c *gin.Context) {
	_, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}

	file, err := c.FormFile("audio")
	if err != nil || file == nil {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "audio file required", 400, nil))
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "failed to read file", 400, nil))
		return
	}
	defer f.Close()

	// Max 25MB (Whisper limiti)
	if file.Size > 25*1024*1024 {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "file too large (max 25MB)", 400, nil))
		return
	}

	// Desteklenen formatlar: flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, webm
	ext := strings.ToLower(file.Filename)
	if idx := strings.LastIndex(ext, "."); idx >= 0 {
		ext = ext[idx:]
	}
	allowed := map[string]bool{
		".flac": true, ".mp3": true, ".mp4": true, ".mpeg": true,
		".mpga": true, ".m4a": true, ".ogg": true, ".wav": true, ".webm": true,
	}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "unsupported audio format", 400, nil))
		return
	}

	data := make([]byte, file.Size)
	if _, err := f.Read(data); err != nil {
		c.JSON(http.StatusInternalServerError, core.NewServiceResponse(false, "failed to read audio", 500, nil))
		return
	}

	text, err := h.whisper.Transcribe(c.Request.Context(), data, file.Filename)
	if err != nil {
		c.JSON(http.StatusBadGateway, core.NewServiceResponse(false, "transcription failed: "+err.Error(), 502, nil))
		return
	}

	c.JSON(http.StatusOK, core.NewServiceResponse(true, "ok", 200, gin.H{"text": strings.TrimSpace(text)}))
}
