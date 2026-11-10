package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// PrerequisiteHandler 先修关系接口
type PrerequisiteHandler struct{ svc *service.PrerequisiteService }

func NewPrerequisiteHandler(svc *service.PrerequisiteService) *PrerequisiteHandler {
	return &PrerequisiteHandler{svc: svc}
}

func (h *PrerequisiteHandler) List(c *gin.Context) {
	courseID, ok := parseID(c)
	if !ok {
		return
	}
	list, err := h.svc.ListByCourse(c.Request.Context(), courseID)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, list)
}

func (h *PrerequisiteHandler) Create(c *gin.Context) {
	courseID, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.PrerequisiteCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	p, err := h.svc.Create(c.Request.Context(), courseID, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, p)
}

func (h *PrerequisiteHandler) Delete(c *gin.Context) {
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
