package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// SemesterHandler 学期字典接口
type SemesterHandler struct{ svc *service.SemesterService }

func NewSemesterHandler(svc *service.SemesterService) *SemesterHandler {
	return &SemesterHandler{svc: svc}
}

func (h *SemesterHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "100"))
	list, total, err := h.svc.List(c.Request.Context(), page, size)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.Page(c, list, total)
}

func (h *SemesterHandler) Create(c *gin.Context) {
	var req dto.SemesterCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, 40001, "参数错误: "+err.Error())
		return
	}
	m, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, m)
}

func (h *SemesterHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.SemesterUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, 40001, "参数错误: "+err.Error())
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, m)
}

func (h *SemesterHandler) Delete(c *gin.Context) {
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
