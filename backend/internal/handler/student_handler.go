package handler

import (
	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/middleware"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// StudentHandler 学生端接口
type StudentHandler struct {
	svc *service.StudentService
}

func NewStudentHandler(svc *service.StudentService) *StudentHandler {
	return &StudentHandler{svc: svc}
}

// Records GET /api/v1/records
func (h *StudentHandler) Records(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)
	list, err := h.svc.MyRecords(c.Request.Context(), userID)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, list)
}

// Audit GET /api/v1/audit/me
func (h *StudentHandler) Audit(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)
	res, err := h.svc.MyAudit(c.Request.Context(), userID)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, res)
}

// Recommendations GET /api/v1/recommendations?semester=
func (h *StudentHandler) Recommendations(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)
	list, err := h.svc.MyRecommendations(c.Request.Context(), userID, c.Query("semester"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, list)
}
