package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// KnowledgeChunkRepo 知识切片数据访问
type KnowledgeChunkRepo struct{ db *gorm.DB }

func NewKnowledgeChunkRepo(db *gorm.DB) *KnowledgeChunkRepo { return &KnowledgeChunkRepo{db: db} }

func (r *KnowledgeChunkRepo) ListAll(ctx context.Context) ([]model.KnowledgeChunk, error) {
	var list []model.KnowledgeChunk
	err := r.db.WithContext(ctx).Order("doc_id ASC, chunk_index ASC").Find(&list).Error
	return list, err
}

func (r *KnowledgeChunkRepo) ListByDoc(ctx context.Context, docID int64) ([]model.KnowledgeChunk, error) {
	var list []model.KnowledgeChunk
	err := r.db.WithContext(ctx).Where("doc_id = ?", docID).Order("chunk_index ASC").Find(&list).Error
	return list, err
}
