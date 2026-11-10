package dto

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string `json:"token"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

// ProfileResponse 当前用户信息
type ProfileResponse struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	StudentNo string `json:"student_no,omitempty"`
}
