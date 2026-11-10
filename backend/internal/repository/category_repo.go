package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// CategoryRepo 课程类别数据访问
type CategoryRepo struct{ db *gorm.DB }

func NewCategoryRepo(db *gorm.DB) *CategoryRepo { return &CategoryRepo{db: db} }

func (r *CategoryRepo) Create(ctx context.Context, m *model.Category) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	var m model.Category
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// ListByPlan 某培养方案下的全部类别（按 sort 排序）
func (r *CategoryRepo) CountByPlan(ctx context.Context, planID int64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.Category{}).Where("plan_id = ?", planID).Count(&n).Error
	return n, err
}

func (r *CategoryRepo) ListByPlan(ctx context.Context, planID int64) ([]model.Category, error) {
	var list []model.Category
	err := r.db.WithContext(ctx).
		Where("plan_id = ?", planID).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

// ListByParent 某类别的直接子模块
func (r *CategoryRepo) ListByParent(ctx context.Context, parentID int64) ([]model.Category, error) {
	var list []model.Category
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *CategoryRepo) List(ctx context.Context, page, size int) ([]model.Category, int64, error) {
	page, size = normalizePage(page, size)
	var list []model.Category
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Category{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *CategoryRepo) Update(ctx context.Context, m *model.Category) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}
