package dto

// ChatRequest 智能体对话请求
type ChatRequest struct {
	Message   string `json:"message" binding:"required,max=2000"`
	SessionID int64  `json:"session_id"`
}
