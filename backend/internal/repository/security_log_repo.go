package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// SecurityLogRepo 安全拦截日志数据访问
type SecurityLogRepo struct{ db *gorm.DB }

func NewSecurityLogRepo(db *gorm.DB) *SecurityLogRepo { return &SecurityLogRepo{db: db} }

func (r *SecurityLogRepo) Create(ctx context.Context, m *model.SecurityLog) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SecurityLogRepo) ListByUser(ctx context.Context, userID int64) ([]model.SecurityLog, error) {
	var list []model.SecurityLog
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("id DESC").Find(&list).Error
	return list, err
}

func (r *SecurityLogRepo) List(ctx context.Context, page, size int) ([]model.SecurityLog, int64, error) {
	page, size = normalizePage(page, size)
	var list []model.SecurityLog
	var total int64
	q := r.db.WithContext(ctx).Model(&model.SecurityLog{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *SecurityLogRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.SecurityLog{}, id).Error
}
