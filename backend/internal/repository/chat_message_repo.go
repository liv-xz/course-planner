package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// ChatMessageRepo 会话消息数据访问
type ChatMessageRepo struct{ db *gorm.DB }

func NewChatMessageRepo(db *gorm.DB) *ChatMessageRepo { return &ChatMessageRepo{db: db} }

func (r *ChatMessageRepo) Create(ctx context.Context, m *model.ChatMessage) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ListBySession 某会话的全部消息（按时间正序）
func (r *ChatMessageRepo) ListBySession(ctx context.Context, sessionID int64) ([]model.ChatMessage, error) {
	var list []model.ChatMessage
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).Order("id ASC").Find(&list).Error
	return list, err
}

func (r *ChatMessageRepo) DeleteBySession(ctx context.Context, sessionID int64) error {
	return r.db.WithContext(ctx).Where("session_id = ?", sessionID).Delete(&model.ChatMessage{}).Error
}
