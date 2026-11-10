package model

import "time"

// CreditAudit 学分审核结果
type CreditAudit struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID        int64     `gorm:"not null;uniqueIndex:uk_audit,priority:1;index:idx_audit_student" json:"student_id"`
	PlanID           int64     `gorm:"not null;uniqueIndex:uk_audit,priority:2" json:"plan_id"`
	CategoryID       int64     `gorm:"not null;uniqueIndex:uk_audit,priority:3" json:"category_id"`
	CompletedCredits float64   `gorm:"type:decimal(5,1);not null;default:0" json:"completed_credits"`
	GapCredits       float64   `gorm:"type:decimal(5,1);not null;default:0" json:"gap_credits"`
	AuditedAt        time.Time `gorm:"autoCreateTime" json:"audited_at"`
}
