package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// CourseCategoryRepo 课程-类别映射数据访问
type CourseCategoryRepo struct{ db *gorm.DB }

func NewCourseCategoryRepo(db *gorm.DB) *CourseCategoryRepo { return &CourseCategoryRepo{db: db} }

func (r *CourseCategoryRepo) Create(ctx context.Context, m *model.CourseCategory) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *CourseCategoryRepo) ListByCourse(ctx context.Context, courseID int64) ([]model.CourseCategory, error) {
	var list []model.CourseCategory
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).Find(&list).Error
	return list, err
}

func (r *CourseCategoryRepo) ListByCategory(ctx context.Context, categoryID int64) ([]model.CourseCategory, error) {
	var list []model.CourseCategory
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&list).Error
	return list, err
}

func (r *CourseCategoryRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.CourseCategory{}, id).Error
}

// DeleteByCourseAndCategory 解除某课程与某类别的映射
func (r *CourseCategoryRepo) DeleteByCourseAndCategory(ctx context.Context, courseID, categoryID int64) error {
	return r.db.WithContext(ctx).
		Where("course_id = ? AND category_id = ?", courseID, categoryID).
		Delete(&model.CourseCategory{}).Error
}

type CourseCategoryInfo struct {
	PlanID          int64   `json:"plan_id"`
	PlanName        string  `json:"plan_name"`
	CategoryID      int64   `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	RequiredCredits float64 `json:"required_credits"`
}

func (r *CourseCategoryRepo) ListByCourseWithInfo(ctx context.Context, courseID int64) ([]CourseCategoryInfo, error) {
	var list []CourseCategoryInfo
	err := r.db.WithContext(ctx).Table("course_category cc").
		Select("c.plan_id, p.name AS plan_name, c.id AS category_id, c.name AS category_name, c.required_credits").
		Joins("JOIN category c ON c.id = cc.category_id").
		Joins("JOIN training_plan p ON p.id = c.plan_id").
		Where("cc.course_id = ?", courseID).
		Order("c.plan_id ASC, c.sort ASC, c.id ASC").
		Scan(&list).Error
	return list, err
}

func (r *CourseCategoryRepo) DeleteByCourseAndPlan(ctx context.Context, courseID, planID int64) error {
	sub := r.db.Model(&model.Category{}).Select("id").Where("plan_id = ?", planID)
	return r.db.WithContext(ctx).
		Where("course_id = ? AND category_id IN (?)", courseID, sub).
		Delete(&model.CourseCategory{}).Error
}

func (r *CourseCategoryRepo) ReplaceForPlan(ctx context.Context, courseID, planID int64, categoryID *int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sub := tx.Model(&model.Category{}).Select("id").Where("plan_id = ?", planID)
		if err := tx.Where("course_id = ? AND category_id IN (?)", courseID, sub).
			Delete(&model.CourseCategory{}).Error; err != nil {
			return err
		}
		if categoryID == nil {
			return nil
		}
		return tx.Create(&model.CourseCategory{CourseID: courseID, CategoryID: *categoryID}).Error
	})
}
