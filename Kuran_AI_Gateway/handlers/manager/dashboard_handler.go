package manager

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/services"
)

type DashboardHandler struct {
	service *services.ManagerDashboardService
}

func NewDashboardHandler(service *services.ManagerDashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) Overview(c *gin.Context) {
	days := queryInt(c, "days", 14)
	data, err := h.service.Overview(c.Request.Context(), days)
	h.respond(c, data, err)
}

func (h *DashboardHandler) Users(c *gin.Context) {
	data, err := h.service.ListUsers(
		c.Request.Context(),
		queryInt(c, "page", 1),
		queryInt(c, "pageSize", 20),
		c.Query("search"),
		c.Query("platform"),
		c.Query("status"),
	)
	h.respond(c, data, err)
}

func (h *DashboardHandler) UserDetail(c *gin.Context) {
	userID, ok := uriID(c, "id")
	if !ok {
		return
	}
	data, err := h.service.UserDetail(c.Request.Context(), userID)
	h.respond(c, data, err)
}

func (h *DashboardHandler) Conversations(c *gin.Context) {
	userID, ok := uriID(c, "id")
	if !ok {
		return
	}
	data, err := h.service.Conversations(
		c.Request.Context(),
		userID,
		queryInt(c, "page", 1),
		queryInt(c, "pageSize", 20),
	)
	h.respond(c, data, err)
}

func (h *DashboardHandler) ConversationMessages(c *gin.Context) {
	sessionID, ok := uriID(c, "id")
	if !ok {
		return
	}
	data, err := h.service.ConversationMessages(c.Request.Context(), sessionID)
	h.respond(c, data, err)
}

func (h *DashboardHandler) Purchases(c *gin.Context) {
	data, err := h.service.Purchases(
		c.Request.Context(),
		queryInt(c, "page", 1),
		queryInt(c, "pageSize", 20),
		c.Query("search"),
		c.Query("provider"),
		c.Query("status"),
	)
	h.respond(c, data, err)
}

func (h *DashboardHandler) Questions(c *gin.Context) {
	data, err := h.service.Questions(
		c.Request.Context(),
		queryInt(c, "page", 1),
		queryInt(c, "pageSize", 20),
		c.Query("search"),
		c.Query("sourceType"),
	)
	h.respond(c, data, err)
}

func (h *DashboardHandler) respond(c *gin.Context, data interface{}, err error) {
	if err != nil {
		status := http.StatusInternalServerError
		message := "Dashboard data could not be loaded"
		if err.Error() == "user not found" || err.Error() == "conversation not found" {
			status = http.StatusNotFound
			message = err.Error()
		}
		c.JSON(status, core.NewServiceResponse(false, message, status, nil))
		return
	}
	c.JSON(http.StatusOK, core.NewServiceResponse(true, "OK", http.StatusOK, data))
}

func queryInt(c *gin.Context, key string, fallback int) int {
	value, err := strconv.Atoi(c.Query(key))
	if err != nil {
		return fallback
	}
	return value
}

func uriID(c *gin.Context, key string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, core.NewServiceResponse(false, "Invalid id", http.StatusBadRequest, nil))
		return 0, false
	}
	return uint(value), true
}
