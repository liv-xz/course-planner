package model

import "time"

// SecurityLog 安全拦截日志
type SecurityLog struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    *int64    `gorm:"index:idx_seclog_user" json:"user_id"`
	Input     *string   `gorm:"type:text" json:"input"`
	Reason    *string   `gorm:"type:varchar(255)" json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}
