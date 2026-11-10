package model

// Major 专业
type Major struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code string `gorm:"type:varchar(32);not null;uniqueIndex:uk_major_code" json:"code"`
	Name string `gorm:"type:varchar(64);not null" json:"name"`
}
