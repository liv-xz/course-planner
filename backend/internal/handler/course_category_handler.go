package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// CourseCategoryHandler 课程-类别归属接口
type CourseCategoryHandler struct {
	svc *service.CourseCategoryService
}

func NewCourseCategoryHandler(svc *service.CourseCategoryService) *CourseCategoryHandler {
	return &CourseCategoryHandler{svc: svc}
}

func (h *CourseCategoryHandler) ListByCourse(c *gin.Context) {
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

func (h *CourseCategoryHandler) SetForPlan(c *gin.Context) {
	planID, ok := parseID(c)
	if !ok {
		return
	}
	courseID, ok := parseIDNamed(c, "courseId")
	if !ok {
		return
	}
	var req dto.SetCourseCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	if err := h.svc.SetForPlan(c.Request.Context(), planID, courseID, req.CategoryID); err != nil {
		respondServiceError(c, err)
		return
	}
	response.OKMessage(c, "保存成功")
}
