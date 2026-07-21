package chat

// SendMessageRequest mesaj gönderme isteği.
// sourceType: meal (kuran), tefsir, hadis (iOS kuran gönderirse meal'e map edilir).
type SendMessageRequest struct {
	SessionID  *uint  `json:"sessionId"`
	SourceType string `json:"sourceType" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// ListSessionsRequest oturum listesi isteği.
type ListSessionsRequest struct {
	PageIndex int `form:"pageIndex"`
	PageSize  int `form:"pageSize"`
	UserID    uint
}

// SetUserID core.UserIDSetter implementasyonu.
func (r *ListSessionsRequest) SetUserID(id uint) { r.UserID = id }

// ListMessagesRequest mesaj listesi isteği.
type ListMessagesRequest struct {
	PageIndex int  `form:"pageIndex"`
	PageSize  int  `form:"pageSize"`
	SessionID uint `uri:"id" binding:"required"`
	UserID    uint
}

// SetUserID core.UserIDSetter implementasyonu.
func (r *ListMessagesRequest) SetUserID(id uint) { r.UserID = id }
