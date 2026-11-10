package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// SecurityHandler 安全日志接口
type SecurityHandler struct{ svc *service.SecurityService }

func NewSecurityHandler(svc *service.SecurityService) *SecurityHandler {
	return &SecurityHandler{svc: svc}
}

func (h *SecurityHandler) Logs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.svc.List(c.Request.Context(), page, size)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.Page(c, list, total)
}
