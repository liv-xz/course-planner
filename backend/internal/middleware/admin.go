package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminOnly 仅允许 admin 角色访问（需在 JWTAuth 之后使用）
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(CtxRole) != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    40301,
				"message": "无权限，仅管理员可操作",
			})
			return
		}
		c.Next()
	}
}
