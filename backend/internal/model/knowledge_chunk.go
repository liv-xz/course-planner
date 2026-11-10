package model

import "time"

// KnowledgeChunk 知识切片（RAG 检索）
type KnowledgeChunk struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DocID      int64     `gorm:"not null;index:idx_chunk_doc" json:"doc_id"`
	ChunkIndex int       `gorm:"not null" json:"chunk_index"`
	Title      string    `gorm:"type:varchar(128);not null" json:"title"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Category   *string   `gorm:"type:varchar(64)" json:"category"`
	Source     *string   `gorm:"type:varchar(255)" json:"source"`
	Keywords   *string   `gorm:"type:varchar(255)" json:"keywords"`
	CreatedAt  time.Time `json:"created_at"`
}
