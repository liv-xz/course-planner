package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sjqu-course-planner/internal/model"
)

// CreditAuditRepo 学分审核结果数据访问
type CreditAuditRepo struct{ db *gorm.DB }

func NewCreditAuditRepo(db *gorm.DB) *CreditAuditRepo { return &CreditAuditRepo{db: db} }

func (r *CreditAuditRepo) Create(ctx context.Context, m *model.CreditAudit) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// Upsert 按 (student_id, plan_id, category_id) 冲突时更新
func (r *CreditAuditRepo) Upsert(ctx context.Context, m *model.CreditAudit) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "student_id"}, {Name: "plan_id"}, {Name: "category_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"completed_credits", "gap_credits", "audited_at"}),
	}).Create(m).Error
}

func (r *CreditAuditRepo) ListByStudent(ctx context.Context, studentID int64) ([]model.CreditAudit, error) {
	var list []model.CreditAudit
	err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Find(&list).Error
	return list, err
}

func (r *CreditAuditRepo) DeleteByStudent(ctx context.Context, studentID int64) error {
	return r.db.WithContext(ctx).Where("student_id = ?", studentID).Delete(&model.CreditAudit{}).Error
}

// CountByPlan 某培养方案的审核记录数
func (r *CreditAuditRepo) CountByPlan(ctx context.Context, planID int64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.CreditAudit{}).Where("plan_id = ?", planID).Count(&n).Error
	return n, err
}

// CountByCategory 某类别的审核记录数
func (r *CreditAuditRepo) CountByCategory(ctx context.Context, categoryID int64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.CreditAudit{}).Where("category_id = ?", categoryID).Count(&n).Error
	return n, err
}
