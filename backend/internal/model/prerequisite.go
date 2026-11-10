package model

// Prerequisite 先修关系（课程自关联 N:M 中间表）
type Prerequisite struct {
	ID             int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID       int64 `gorm:"not null;uniqueIndex:uk_prereq,priority:1;index:idx_prereq_course" json:"course_id"`
	PrereqCourseID int64 `gorm:"not null;uniqueIndex:uk_prereq,priority:2;index:idx_prereq_pre" json:"prereq_course_id"`
}
