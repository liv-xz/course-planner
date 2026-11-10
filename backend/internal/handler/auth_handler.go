package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/middleware"
	"sjqu-course-planner/internal/service"
)

// AuthHandler 认证接口
type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 40001, "message": "参数错误"})
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 40101, "message": err.Error()})
		return
	}

	name := ""
	if res.User.Name != nil {
		name = *res.User.Name
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": dto.LoginResponse{
			Token:    res.Token,
			UserID:   res.User.ID,
			Username: res.User.Username,
			Name:     name,
			Role:     res.Role,
		},
	})
}

// Profile GET /api/v1/auth/profile（需登录）
func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)
	user, profile, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 40401, "message": "用户不存在"})
		return
	}

	name := ""
	if user.Name != nil {
		name = *user.Name
	}
	studentNo := ""
	if profile != nil {
		studentNo = profile.StudentNo
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dto.ProfileResponse{
			UserID:    user.ID,
			Username:  user.Username,
			Name:      name,
			Role:      c.GetString(middleware.CtxRole),
			StudentNo: studentNo,
		},
	})
}
