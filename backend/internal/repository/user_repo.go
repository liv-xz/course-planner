package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// UserRepo 用户数据访问
type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, m *model.User) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("is_deleted = 0").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("username = ? AND is_deleted = 0", username).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *UserRepo) List(ctx context.Context, page, size int) ([]model.User, int64, error) {
	page, size = normalizePage(page, size)
	var list []model.User
	var total int64
	q := r.db.WithContext(ctx).Model(&model.User{}).Where("is_deleted = 0")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *UserRepo) Update(ctx context.Context, m *model.User) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// Delete 逻辑删除
func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("is_deleted", 1).Error
}
