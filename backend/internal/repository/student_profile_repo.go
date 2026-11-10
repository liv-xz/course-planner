package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// StudentProfileRepo 学生档案数据访问
type StudentProfileRepo struct{ db *gorm.DB }

func NewStudentProfileRepo(db *gorm.DB) *StudentProfileRepo { return &StudentProfileRepo{db: db} }

func (r *StudentProfileRepo) Create(ctx context.Context, m *model.StudentProfile) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *StudentProfileRepo) GetByUserID(ctx context.Context, userID int64) (*model.StudentProfile, error) {
	var m model.StudentProfile
	if err := r.db.WithContext(ctx).First(&m, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *StudentProfileRepo) GetByStudentNo(ctx context.Context, studentNo string) (*model.StudentProfile, error) {
	var m model.StudentProfile
	if err := r.db.WithContext(ctx).Where("student_no = ?", studentNo).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *StudentProfileRepo) Update(ctx context.Context, m *model.StudentProfile) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *StudentProfileRepo) Delete(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Delete(&model.StudentProfile{}, "user_id = ?", userID).Error
}

func (r *StudentProfileRepo) CountByMajor(ctx context.Context, majorID int64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.StudentProfile{}).Where("major_id = ?", majorID).Count(&n).Error
	return n, err
}

// StudentRow 管理端学生列表行
type StudentRow struct {
	UserID    int64   `json:"user_id"`
	Username  string  `json:"username"`
	Name      *string `json:"name"`
	StudentNo string  `json:"student_no"`
	MajorID   int64   `json:"major_id"`
	MajorName string  `json:"major_name"`
	Grade     string  `json:"grade"`
	Status    int8    `json:"status"`
}

func (r *StudentProfileRepo) GetRowByUserID(ctx context.Context, userID int64) (*StudentRow, error) {
	var row StudentRow
	err := r.db.WithContext(ctx).Table("student_profile sp").
		Select("sp.user_id, u.username, u.name, sp.student_no, sp.major_id, m.name AS major_name, sp.grade, u.status").
		Joins("JOIN `user` u ON u.id = sp.user_id").
		Joins("JOIN major m ON m.id = sp.major_id").
		Where("sp.user_id = ? AND u.is_deleted = 0", userID).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.UserID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (r *StudentProfileRepo) ListWithUser(ctx context.Context, page, size int, keyword string) ([]StudentRow, int64, error) {
	page, size = normalizePage(page, size)
	base := r.db.WithContext(ctx).Table("student_profile sp").
		Joins("JOIN `user` u ON u.id = sp.user_id").
		Joins("JOIN major m ON m.id = sp.major_id").
		Where("u.is_deleted = 0")
	if keyword != "" {
		like := "%" + keyword + "%"
		base = base.Where("(u.username LIKE ? OR u.name LIKE ? OR sp.student_no LIKE ?)", like, like, like)
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []StudentRow
	err := base.Select("sp.user_id, u.username, u.name, sp.student_no, sp.major_id, m.name AS major_name, sp.grade, u.status").
		Order("sp.user_id ASC").Offset((page - 1) * size).Limit(size).Scan(&list).Error
	return list, total, err
}
