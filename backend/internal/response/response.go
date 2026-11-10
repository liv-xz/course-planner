package response

import "github.com/gin-gonic/gin"

// OK 成功返回 data
func OK(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": data})
}

// OKMessage 成功但无数据（如删除）
func OKMessage(c *gin.Context, message string) {
	c.JSON(200, gin.H{"code": 0, "message": message})
}

// Page 分页返回
func Page(c *gin.Context, list interface{}, total int64) {
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    gin.H{"list": list, "total": total},
	})
}

// Error 错误返回
func Error(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, gin.H{"code": code, "message": message})
}
