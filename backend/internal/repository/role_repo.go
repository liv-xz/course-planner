package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// RoleRepo 角色数据访问
type RoleRepo struct{ db *gorm.DB }

func NewRoleRepo(db *gorm.DB) *RoleRepo { return &RoleRepo{db: db} }

func (r *RoleRepo) Create(ctx context.Context, m *model.Role) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *RoleRepo) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	var m model.Role
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *RoleRepo) GetByName(ctx context.Context, name string) (*model.Role, error) {
	var m model.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *RoleRepo) List(ctx context.Context, page, size int) ([]model.Role, int64, error) {
	page, size = normalizePage(page, size)
	var list []model.Role
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Role{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *RoleRepo) Update(ctx context.Context, m *model.Role) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *RoleRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}
