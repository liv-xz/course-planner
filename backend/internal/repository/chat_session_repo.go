package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// ChatSessionRepo 智能体会话数据访问
type ChatSessionRepo struct{ db *gorm.DB }

func NewChatSessionRepo(db *gorm.DB) *ChatSessionRepo { return &ChatSessionRepo{db: db} }

func (r *ChatSessionRepo) Create(ctx context.Context, m *model.ChatSession) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *ChatSessionRepo) GetByID(ctx context.Context, id int64) (*model.ChatSession, error) {
	var m model.ChatSession
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// ListByUser 某用户的会话列表
func (r *ChatSessionRepo) ListByUser(ctx context.Context, userID int64) ([]model.ChatSession, error) {
	var list []model.ChatSession
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("id DESC").Find(&list).Error
	return list, err
}

func (r *ChatSessionRepo) Update(ctx context.Context, m *model.ChatSession) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *ChatSessionRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.ChatSession{}, id).Error
}
