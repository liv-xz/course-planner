package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StudentOnly 仅允许 student 角色访问（需在 JWTAuth 之后使用）
func StudentOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(CtxRole) != "student" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    40301,
				"message": "无权限，仅学生可访问",
			})
			return
		}
		c.Next()
	}
}
