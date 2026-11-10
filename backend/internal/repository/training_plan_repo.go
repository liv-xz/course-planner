package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// TrainingPlanRepo 培养方案数据访问
type TrainingPlanRepo struct{ db *gorm.DB }

func NewTrainingPlanRepo(db *gorm.DB) *TrainingPlanRepo { return &TrainingPlanRepo{db: db} }

func (r *TrainingPlanRepo) Create(ctx context.Context, m *model.TrainingPlan) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *TrainingPlanRepo) GetByID(ctx context.Context, id int64) (*model.TrainingPlan, error) {
	var m model.TrainingPlan
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// GetByMajorGrade 按专业 + 年级取培养方案
func (r *TrainingPlanRepo) GetByMajorGrade(ctx context.Context, majorID int64, grade string) (*model.TrainingPlan, error) {
	var m model.TrainingPlan
	if err := r.db.WithContext(ctx).
		Where("major_id = ? AND grade = ?", majorID, grade).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *TrainingPlanRepo) List(ctx context.Context, page, size int) ([]model.TrainingPlan, int64, error) {
	page, size = normalizePage(page, size)
	var list []model.TrainingPlan
	var total int64
	q := r.db.WithContext(ctx).Model(&model.TrainingPlan{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *TrainingPlanRepo) Update(ctx context.Context, m *model.TrainingPlan) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *TrainingPlanRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.TrainingPlan{}, id).Error
}

func (r *TrainingPlanRepo) CountByMajor(ctx context.Context, majorID int64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.TrainingPlan{}).Where("major_id = ?", majorID).Count(&n).Error
	return n, err
}
