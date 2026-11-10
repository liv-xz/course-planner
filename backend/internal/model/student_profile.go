package model

// StudentProfile 学生档案（学生专属字段，1:1 关联 user）
type StudentProfile struct {
	UserID    int64  `gorm:"primaryKey" json:"user_id"`
	MajorID   int64  `gorm:"not null;index:idx_sp_major" json:"major_id"`
	StudentNo string `gorm:"type:varchar(32);not null;uniqueIndex:uk_sp_student_no" json:"student_no"`
	Grade     string `gorm:"type:varchar(16);not null" json:"grade"`
}
