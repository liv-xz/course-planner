package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/middleware"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// ChatHandler 智能体对话接口
type ChatHandler struct{ svc *service.AgentService }

func NewChatHandler(svc *service.AgentService) *ChatHandler { return &ChatHandler{svc: svc} }

func (h *ChatHandler) Chat(c *gin.Context) {
	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}
	userID := c.GetInt64(middleware.CtxUserID)
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	flusher, _ := c.Writer.(http.Flusher)
	emit := func(event string, data any) {
		b, _ := json.Marshal(data)
		fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, b)
		if flusher != nil {
			flusher.Flush()
		}
	}
	if err := h.svc.ChatStream(c.Request.Context(), userID, req.SessionID, req.Message, emit); err != nil {
		emit("error", map[string]any{"message": err.Error()})
	}
}

func (h *ChatHandler) Sessions(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)
	list, err := h.svc.ListSessions(c.Request.Context(), userID)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, list)
}

func (h *ChatHandler) Messages(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)
	id, ok := parseID(c)
	if !ok {
		return
	}
	list, err := h.svc.ListMessages(c.Request.Context(), userID, id)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, list)
}

func (h *ChatHandler) DeleteSession(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteSession(c.Request.Context(), userID, id); err != nil {
		respondServiceError(c, err)
		return
	}
	response.OKMessage(c, "删除成功")
}

func (h *ChatHandler) Status(c *gin.Context) {
	response.OK(c, gin.H{
		"enabled": true,
		"mode":    h.svc.Mode(),
	})
}
