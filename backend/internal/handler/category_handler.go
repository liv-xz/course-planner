package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// CategoryHandler 课程类别接口
type CategoryHandler struct{ svc *service.CategoryService }

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// ListTree GET /plans/:id/categories
func (h *CategoryHandler) ListTree(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	tree, err := h.svc.ListTree(c.Request.Context(), id)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, tree)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	cat, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, cat)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	cat, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, cat)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
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
