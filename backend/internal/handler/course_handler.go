package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// CourseHandler 课程库接口
type CourseHandler struct{ svc *service.CourseService }

func NewCourseHandler(svc *service.CourseService) *CourseHandler { return &CourseHandler{svc: svc} }

func (h *CourseHandler) List(c *gin.Context) {
	var q dto.CourseListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	list, total, err := h.svc.List(c.Request.Context(), q)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.Page(c, list, total)
}

func (h *CourseHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	course, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, course)
}

func (h *CourseHandler) Create(c *gin.Context) {
	var req dto.CourseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	course, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, course)
}

func (h *CourseHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.CourseUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	course, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, course)
}

func (h *CourseHandler) Delete(c *gin.Context) {
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
