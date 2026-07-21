package chat

// TTSRequest sesli okuma isteği.
type TTSRequest struct {
	Text  string `json:"text" binding:"required"`
	Voice string `json:"voice"` // alloy, echo, fable, onyx, nova, shimmer; boşsa nova
}
