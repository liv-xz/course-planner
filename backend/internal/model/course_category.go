package model

// CourseCategory 课程-类别映射（N:M 中间表）
type CourseCategory struct {
	ID         int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID   int64 `gorm:"not null;uniqueIndex:uk_cc,priority:1;index:idx_cc_course" json:"course_id"`
	CategoryID int64 `gorm:"not null;uniqueIndex:uk_cc,priority:2;index:idx_cc_category" json:"category_id"`
}
