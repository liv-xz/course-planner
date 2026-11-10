package model

import "time"

// TrainingPlan 培养方案（按专业 + 年级）
type TrainingPlan struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MajorID      int64     `gorm:"not null;uniqueIndex:uk_plan,priority:1;index:idx_plan_major" json:"major_id"`
	Grade        string    `gorm:"type:varchar(16);not null;uniqueIndex:uk_plan,priority:2" json:"grade"`
	Name         string    `gorm:"type:varchar(128);not null" json:"name"`
	TotalCredits *float64  `gorm:"type:decimal(5,1)" json:"total_credits"`
	CreatedAt    time.Time `json:"created_at"`
}
