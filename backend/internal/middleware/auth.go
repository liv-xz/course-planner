package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/utils"
)

// Context 中存放用户信息的 key
const (
	CtxUserID = "user_id"
	CtxRole   = "role"
)

// JWTAuth JWT 鉴权中间件：校验 Authorization: Bearer <token>
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    40101,
				"message": "未登录或缺少 Token",
			})
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		claims, err := utils.ParseToken(tokenStr, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    40101,
				"message": "Token 无效或已过期",
			})
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}
