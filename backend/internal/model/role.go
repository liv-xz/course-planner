package model

// Role 角色
type Role struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(32);not null;uniqueIndex:uk_role_name" json:"name"`
	Description string `gorm:"type:varchar(128)" json:"description"`
}
