package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// MajorHandler 专业接口
type MajorHandler struct{ svc *service.MajorService }

func NewMajorHandler(svc *service.MajorService) *MajorHandler { return &MajorHandler{svc: svc} }

func (h *MajorHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "100"))
	list, total, err := h.svc.List(c.Request.Context(), page, size)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.Page(c, list, total)
}

func (h *MajorHandler) Create(c *gin.Context) {
	var req dto.MajorCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	m, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, m)
}

func (h *MajorHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.MajorUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, m)
}

func (h *MajorHandler) Delete(c *gin.Context) {
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
