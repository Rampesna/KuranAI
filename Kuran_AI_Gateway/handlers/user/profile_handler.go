package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/middleware"
	userRequest "github.com/rampesna/kuran-ai-gateway/requests/user"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// ProfileHandler profil güncelleme handler.
type ProfileHandler struct {
	userService *services.UserService
}

// NewProfileHandler ProfileHandler oluşturur.
func NewProfileHandler(userService *services.UserService) *ProfileHandler {
	return &ProfileHandler{userService: userService}
}

// UpdateUsername kullanıcı adını günceller (uniqueness kontrolü).
func (h *ProfileHandler) UpdateUsername(c *gin.Context) {
	userIDVal, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	var req userRequest.UpdateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if sr := core.ValidationErrorResponse(err); sr != nil {
			c.JSON(sr.StatusCode, sr)
			return
		}
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	sr := h.userService.UpdateUsername(c.Request.Context(), userID, req.Username)
	c.JSON(sr.StatusCode, sr)
}

// UploadAvatar profil resmi yükler.
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	userIDVal, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	file, err := c.FormFile("avatar")
	if err != nil || file == nil {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "avatar file required", 400, nil))
		return
	}
	sr := h.userService.UploadAvatar(c.Request.Context(), userID, file)
	c.JSON(sr.StatusCode, sr)
}

// RemoveAvatar profil resmini siler.
func (h *ProfileHandler) RemoveAvatar(c *gin.Context) {
	userIDVal, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	sr := h.userService.RemoveAvatar(c.Request.Context(), userID)
	c.JSON(sr.StatusCode, sr)
}

// UpdatePassword şifreyi günceller.
func (h *ProfileHandler) UpdatePassword(c *gin.Context) {
	userIDVal, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	var req userRequest.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if sr := core.ValidationErrorResponse(err); sr != nil {
			c.JSON(sr.StatusCode, sr)
			return
		}
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	sr := h.userService.UpdatePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	c.JSON(sr.StatusCode, sr)
}

// DeleteAccount hesabı kalıcı olarak siler (Apple Guideline 5.1.1(v)).
func (h *ProfileHandler) DeleteAccount(c *gin.Context) {
	userIDVal, ok := c.Get(string(middleware.ContextKeyUserID))
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, core.NewServiceResponse(false, "Unauthorized", 401, nil))
		return
	}
	sr := h.userService.DeleteAccount(c.Request.Context(), userID)
	c.JSON(sr.StatusCode, sr)
}
