package model

// Semester 学期字典（日历学期）
type Semester struct {
	ID           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code         string `gorm:"type:varchar(16);not null;uniqueIndex:uk_semester_code" json:"code"`
	AcademicYear string `gorm:"type:varchar(9);not null" json:"academic_year"`
	TermNo       int8   `gorm:"not null" json:"term_no"`
	SequenceNo   int    `gorm:"not null;index:idx_semester_seq" json:"sequence_no"`
	DisplayName  string `gorm:"type:varchar(32);not null" json:"display_name"`
}
