package public

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/models"
)

type createSupportRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=128"`
	Email   string `json:"email" binding:"required,email,max=256"`
	Subject string `json:"subject" binding:"max=64"`
	Message string `json:"message" binding:"required,min=10,max=5000"`
}

// SubmitSupportRequest web/uygulama üzerinden destek talebi oluşturur.
func SubmitSupportRequest(c *gin.Context) {
	var req createSupportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if sr := core.ValidationErrorResponse(err); sr != nil {
			c.JSON(sr.StatusCode, sr)
			return
		}
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}

	subject := strings.TrimSpace(req.Subject)
	if subject == "" {
		subject = "general"
	}

	supportReq := models.SupportRequest{
		Name:      strings.TrimSpace(req.Name),
		Email:     strings.TrimSpace(req.Email),
		Subject:   subject,
		Message:   strings.TrimSpace(req.Message),
		IPAddress: c.ClientIP(),
		Status:    "pending",
	}

	if err := database.DB.WithContext(c.Request.Context()).Create(&supportReq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, core.NewServiceResponse(false, "Failed to submit request", 500, nil))
		return
	}

	c.JSON(http.StatusCreated, core.NewServiceResponse(true, "Support request submitted", 201, nil))
}
