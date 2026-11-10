package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// KnowledgeDocRepo 知识文档数据访问
type KnowledgeDocRepo struct{ db *gorm.DB }

func NewKnowledgeDocRepo(db *gorm.DB) *KnowledgeDocRepo { return &KnowledgeDocRepo{db: db} }

func (r *KnowledgeDocRepo) Create(ctx context.Context, m *model.KnowledgeDoc) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *KnowledgeDocRepo) GetByID(ctx context.Context, id int64) (*model.KnowledgeDoc, error) {
	var m model.KnowledgeDoc
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *KnowledgeDocRepo) ListByPlan(ctx context.Context, planID int64) ([]model.KnowledgeDoc, error) {
	var list []model.KnowledgeDoc
	err := r.db.WithContext(ctx).Where("plan_id = ?", planID).Find(&list).Error
	return list, err
}

func (r *KnowledgeDocRepo) List(ctx context.Context, page, size int) ([]model.KnowledgeDoc, int64, error) {
	page, size = normalizePage(page, size)
	var list []model.KnowledgeDoc
	var total int64
	q := r.db.WithContext(ctx).Model(&model.KnowledgeDoc{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *KnowledgeDocRepo) Update(ctx context.Context, m *model.KnowledgeDoc) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *KnowledgeDocRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.KnowledgeDoc{}, id).Error
}
