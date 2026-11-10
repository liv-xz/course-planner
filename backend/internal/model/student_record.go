package model

// StudentRecord 学生修读记录
type StudentRecord struct {
	ID        int64    `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID int64    `gorm:"not null;index:idx_record_student" json:"student_id"`
	CourseID  int64    `gorm:"not null;index:idx_record_course" json:"course_id"`
	Score     *float64 `gorm:"type:decimal(5,1)" json:"score"`
	Semester  *string  `gorm:"type:varchar(16)" json:"semester"`
	Status    int8     `gorm:"not null;default:0" json:"status"` // 0在读 1通过 2未通过
	Credits   *float64 `gorm:"type:decimal(4,1)" json:"credits"`
}
