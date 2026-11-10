package model

import "time"

// ChatSession 智能体会话
type ChatSession struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;index:idx_session_user" json:"user_id"`
	Title     *string   `gorm:"type:varchar(128)" json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
