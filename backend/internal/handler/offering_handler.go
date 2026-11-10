package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// OfferingHandler 开课信息接口
type OfferingHandler struct{ svc *service.OfferingService }

func NewOfferingHandler(svc *service.OfferingService) *OfferingHandler {
	return &OfferingHandler{svc: svc}
}

func (h *OfferingHandler) List(c *gin.Context) {
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

func (h *OfferingHandler) Create(c *gin.Context) {
	courseID, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.OfferingCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}
	o, err := h.svc.Create(c.Request.Context(), courseID, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, o)
}

func (h *OfferingHandler) Delete(c *gin.Context) {
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
