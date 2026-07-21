package verse

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	verseRequest "github.com/rampesna/kuran-ai-gateway/requests/verse"
	"github.com/rampesna/kuran-ai-gateway/services"
)

// VerseHandler günün ayeti handler.
type VerseHandler struct {
	verseService *services.DailyVerseService
}

// NewVerseHandler yeni handler oluşturur.
func NewVerseHandler(verseService *services.DailyVerseService) *VerseHandler {
	return &VerseHandler{verseService: verseService}
}

// GetToday bugünün ayetini döndürür (public endpoint).
// GET /api/verse-of-day?lang=tr
func (h *VerseHandler) GetToday(c *gin.Context) {
	lang := c.DefaultQuery("lang", "")
	sr := h.verseService.GetToday(c.Request.Context(), lang)
	c.JSON(sr.StatusCode, sr)
}

// List tüm ayetleri listeler (manager endpoint).
// GET /api/manager/verse-of-day
func (h *VerseHandler) List(c *gin.Context) {
	sr := h.verseService.List(c.Request.Context())
	c.JSON(sr.StatusCode, sr)
}

// GetByID tek ayet detayı (manager endpoint).
// GET /api/manager/verse-of-day/:id
func (h *VerseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid ID", 400, nil))
		return
	}
	sr := h.verseService.GetByID(c.Request.Context(), uint(id))
	c.JSON(sr.StatusCode, sr)
}

// Create yeni ayet oluşturur (manager endpoint).
// POST /api/manager/verse-of-day
func (h *VerseHandler) Create(c *gin.Context) {
	var req verseRequest.CreateDailyVerseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if sr := core.ValidationErrorResponse(err); sr != nil {
			c.JSON(sr.StatusCode, sr)
			return
		}
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	sr := h.verseService.Create(c.Request.Context(), &req)
	c.JSON(sr.StatusCode, sr)
}

// Update ayet günceller (manager endpoint).
// PUT /api/manager/verse-of-day/:id
func (h *VerseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid ID", 400, nil))
		return
	}
	var req verseRequest.UpdateDailyVerseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if sr := core.ValidationErrorResponse(err); sr != nil {
			c.JSON(sr.StatusCode, sr)
			return
		}
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid request", 400, nil))
		return
	}
	sr := h.verseService.Update(c.Request.Context(), uint(id), &req)
	c.JSON(sr.StatusCode, sr)
}

// Delete ayet siler (manager endpoint).
// DELETE /api/manager/verse-of-day/:id
func (h *VerseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid ID", 400, nil))
		return
	}
	sr := h.verseService.Delete(c.Request.Context(), uint(id))
	c.JSON(sr.StatusCode, sr)
}
