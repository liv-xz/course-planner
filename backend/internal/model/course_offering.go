package model

// CourseOffering 开课信息
type CourseOffering struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID int64  `gorm:"not null;index:idx_offering_course" json:"course_id"`
	Semester string `gorm:"type:varchar(16);not null" json:"semester"`
}
