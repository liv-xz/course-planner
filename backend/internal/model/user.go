package model

import "time"

// User 用户（学生 / 教务管理员）——只放通用字段，学生专属字段见 StudentProfile
type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"type:varchar(64);not null;uniqueIndex:uk_user_username" json:"username"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	RoleID       int64     `gorm:"not null;index:idx_user_role" json:"role_id"`
	Name         *string   `gorm:"type:varchar(32)" json:"name"`
	Status       int8      `gorm:"not null;default:1" json:"status"`
	IsDeleted    int8      `gorm:"not null;default:0" json:"is_deleted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
