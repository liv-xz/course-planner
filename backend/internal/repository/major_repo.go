package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// MajorRepo 专业数据访问
type MajorRepo struct{ db *gorm.DB }

func NewMajorRepo(db *gorm.DB) *MajorRepo { return &MajorRepo{db: db} }

func (r *MajorRepo) Create(ctx context.Context, m *model.Major) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MajorRepo) GetByID(ctx context.Context, id int64) (*model.Major, error) {
	var m model.Major
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MajorRepo) GetByCode(ctx context.Context, code string) (*model.Major, error) {
	var m model.Major
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MajorRepo) List(ctx context.Context, page, size int) ([]model.Major, int64, error) {
	page, size = normalizePage(page, size)
	var list []model.Major
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Major{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *MajorRepo) Update(ctx context.Context, m *model.Major) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *MajorRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Major{}, id).Error
}
