package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// CourseOfferingRepo 开课信息数据访问
type CourseOfferingRepo struct{ db *gorm.DB }

func NewCourseOfferingRepo(db *gorm.DB) *CourseOfferingRepo { return &CourseOfferingRepo{db: db} }

func (r *CourseOfferingRepo) Create(ctx context.Context, m *model.CourseOffering) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *CourseOfferingRepo) GetByID(ctx context.Context, id int64) (*model.CourseOffering, error) {
	var m model.CourseOffering
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *CourseOfferingRepo) ListByCourse(ctx context.Context, courseID int64) ([]model.CourseOffering, error) {
	var list []model.CourseOffering
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).Order("semester ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *CourseOfferingRepo) GetByCourseSemester(ctx context.Context, courseID int64, semester string) (*model.CourseOffering, error) {
	var m model.CourseOffering
	if err := r.db.WithContext(ctx).Where("course_id = ? AND semester = ?", courseID, semester).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *CourseOfferingRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.CourseOffering{}, id).Error
}
