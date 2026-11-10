package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// SemesterRepo 学期字典数据访问
type SemesterRepo struct{ db *gorm.DB }

func NewSemesterRepo(db *gorm.DB) *SemesterRepo { return &SemesterRepo{db: db} }

func (r *SemesterRepo) Create(ctx context.Context, m *model.Semester) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SemesterRepo) GetByID(ctx context.Context, id int64) (*model.Semester, error) {
	var m model.Semester
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SemesterRepo) GetByCode(ctx context.Context, code string) (*model.Semester, error) {
	var m model.Semester
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SemesterRepo) List(ctx context.Context, page, size int) ([]model.Semester, int64, error) {
	page, size = normalizePage(page, size)
	var list []model.Semester
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Semester{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("sequence_no ASC, id ASC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *SemesterRepo) Update(ctx context.Context, m *model.Semester) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *SemesterRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Semester{}, id).Error
}

func (r *SemesterRepo) CountOfferings(ctx context.Context, code string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.CourseOffering{}).Where("semester = ?", code).Count(&n).Error
	return n, err
}

func (r *SemesterRepo) CountRecords(ctx context.Context, code string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.StudentRecord{}).Where("semester = ?", code).Count(&n).Error
	return n, err
}
