package model

import "time"

// KnowledgeDoc 知识文档（RAG 元数据）
type KnowledgeDoc struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PlanID     *int64    `gorm:"index:idx_kdoc_plan" json:"plan_id"`
	Title      string    `gorm:"type:varchar(128);not null" json:"title"`
	Source     *string   `gorm:"type:varchar(255)" json:"source"`
	ChunkCount int       `gorm:"not null;default:0" json:"chunk_count"`
	CreatedAt  time.Time `json:"created_at"`
}
