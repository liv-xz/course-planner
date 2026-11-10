package model

// Category 课程类别（parent_id 组层级，支持子模块）
type Category struct {
	ID              int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	PlanID          int64   `gorm:"not null;index:idx_cat_plan" json:"plan_id"`
	ParentID        *int64  `gorm:"index:idx_cat_parent" json:"parent_id"`
	Name            string  `gorm:"type:varchar(64);not null" json:"name"`
	RequiredCredits float64 `gorm:"type:decimal(5,1);not null;default:0" json:"required_credits"`
	Sort            int     `gorm:"not null;default:0" json:"sort"`
}
