package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sjqu-course-planner/internal/logger"
)

// Recovery 捕获 panic，记录日志并返回统一错误
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		logger.Log.Error("panic_recovered", zap.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    50001,
			"message": "服务器内部错误",
		})
	})
}
