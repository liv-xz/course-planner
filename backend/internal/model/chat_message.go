package model

import "time"

// ChatMessage 会话消息
type ChatMessage struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID int64     `gorm:"not null;index:idx_msg_session" json:"session_id"`
	Role      string    `gorm:"type:varchar(16);not null" json:"role"` // user/assistant/tool
	Content   *string   `gorm:"type:text" json:"content"`
	ToolCalls *string   `gorm:"type:json" json:"tool_calls"`
	CreatedAt time.Time `json:"created_at"`
}
