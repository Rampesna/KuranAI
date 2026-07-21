package user

// SendMessageRequest chat sendMessage için request (prompt, provider, model, apiKey; isteğe bağlı sessionId).
type SendMessageRequest struct {
	Prompt   string `json:"prompt" binding:"required"`
	Provider string `json:"provider" binding:"required"` // openai, google, anthropic, moonshot (slug)
	Model    string `json:"model" binding:"required"`   // api model id (apiId)
	ApiKey    string `json:"apiKey" binding:"required"` // sağlayıcı API key
	SessionID *uint  `json:"sessionId,omitempty"`        // isteğe bağlı; yoksa yeni session oluşturulur
	UserID    *uint  // core.WithValidatedJSONAndUser tarafından set edilir
}

// SetUserID core.UserIDSetter implementasyonu.
func (r *SendMessageRequest) SetUserID(id uint) { r.UserID = &id }
