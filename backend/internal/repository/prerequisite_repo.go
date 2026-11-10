package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// PrerequisiteRepo 先修关系数据访问
type PrerequisiteRepo struct{ db *gorm.DB }

func NewPrerequisiteRepo(db *gorm.DB) *PrerequisiteRepo { return &PrerequisiteRepo{db: db} }

func (r *PrerequisiteRepo) Create(ctx context.Context, m *model.Prerequisite) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PrerequisiteRepo) GetByID(ctx context.Context, id int64) (*model.Prerequisite, error) {
	var m model.Prerequisite
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// ListByCourse 某门课的所有先修课 id
func (r *PrerequisiteRepo) ListByCourse(ctx context.Context, courseID int64) ([]model.Prerequisite, error) {
	var list []model.Prerequisite
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).Find(&list).Error
	return list, err
}

// ListByPrereq 某门课被哪些课当作先修
func (r *PrerequisiteRepo) ListByPrereq(ctx context.Context, prereqCourseID int64) ([]model.Prerequisite, error) {
	var list []model.Prerequisite
	err := r.db.WithContext(ctx).Where("prereq_course_id = ?", prereqCourseID).Find(&list).Error
	return list, err
}

func (r *PrerequisiteRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Prerequisite{}, id).Error
}
