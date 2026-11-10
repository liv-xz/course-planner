package model

// Course 课程库
type Course struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string  `gorm:"type:varchar(32);not null;uniqueIndex:uk_course_code" json:"code"`
	Name      string  `gorm:"type:varchar(128);not null" json:"name"`
	Credits   float64 `gorm:"type:decimal(4,1);not null" json:"credits"`
	Nature    int8    `gorm:"not null;default:1" json:"nature"` // 1必修 2选修
	IsDeleted int8    `gorm:"not null;default:0" json:"is_deleted"`
}
