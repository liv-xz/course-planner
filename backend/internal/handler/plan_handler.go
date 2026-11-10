package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// PlanHandler 培养方案接口
type PlanHandler struct{ svc *service.PlanService }

func NewPlanHandler(svc *service.PlanService) *PlanHandler { return &PlanHandler{svc: svc} }

func (h *PlanHandler) List(c *gin.Context) {
	list, total, err := h.svc.List(c.Request.Context(), queryInt(c, "page", 1), queryInt(c, "size", 10))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.Page(c, list, total)
}

func (h *PlanHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	p, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, p)
}

func (h *PlanHandler) Create(c *gin.Context) {
	var req dto.PlanCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	p, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, p)
}

func (h *PlanHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.PlanUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	p, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, p)
}

func (h *PlanHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondServiceError(c, err)
		return
	}
	response.OKMessage(c, "删除成功")
}
