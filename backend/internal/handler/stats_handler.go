package handler

import (
	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// StatsHandler 管理端统计接口
type StatsHandler struct{ svc *service.StatsService }

func NewStatsHandler(svc *service.StatsService) *StatsHandler { return &StatsHandler{svc: svc} }

func (h *StatsHandler) Overview(c *gin.Context) {
	data, err := h.svc.Overview(c.Request.Context())
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, data)
}

func (h *StatsHandler) CreditGap(c *gin.Context) {
	data, err := h.svc.CreditGap(c.Request.Context())
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, data)
}
