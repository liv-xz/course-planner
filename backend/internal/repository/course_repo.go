package repository

import (
	"context"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// CourseRepo 课程库数据访问
type CourseRepo struct{ db *gorm.DB }

func NewCourseRepo(db *gorm.DB) *CourseRepo { return &CourseRepo{db: db} }

func (r *CourseRepo) Create(ctx context.Context, m *model.Course) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *CourseRepo) GetByID(ctx context.Context, id int64) (*model.Course, error) {
	var m model.Course
	if err := r.db.WithContext(ctx).Where("is_deleted = 0").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *CourseRepo) GetByCode(ctx context.Context, code string) (*model.Course, error) {
	var m model.Course
	if err := r.db.WithContext(ctx).Where("code = ? AND is_deleted = 0", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *CourseRepo) GetByIDIncludingDeleted(ctx context.Context, id int64) (*model.Course, error) {
	var m model.Course
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *CourseRepo) GetByCodeIncludingDeleted(ctx context.Context, code string) (*model.Course, error) {
	var m model.Course
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// ListByCategory 某类别下的课程（经 course_category 关联）
func (r *CourseRepo) ListByCategory(ctx context.Context, categoryID int64) ([]model.Course, error) {
	var list []model.Course
	err := r.db.WithContext(ctx).
		Joins("JOIN course_category cc ON cc.course_id = course.id").
		Where("cc.category_id = ? AND course.is_deleted = 0", categoryID).
		Find(&list).Error
	return list, err
}

func (r *CourseRepo) List(ctx context.Context, page, size int) ([]model.Course, int64, error) {
	return r.Search(ctx, page, size, "", 0)
}

// Search 分页 + 关键词（名称/代码）+ 类别筛选
func (r *CourseRepo) Search(ctx context.Context, page, size int, keyword string, categoryID int64) ([]model.Course, int64, error) {
	page, size = normalizePage(page, size)
	q := r.db.WithContext(ctx).Model(&model.Course{}).Where("is_deleted = 0")
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("(name LIKE ? OR code LIKE ?)", like, like)
	}
	if categoryID > 0 {
		q = q.Where("id IN (SELECT course_id FROM course_category WHERE category_id = ?)", categoryID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Course
	if err := q.Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *CourseRepo) Update(ctx context.Context, m *model.Course) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// Delete 逻辑删除
func (r *CourseRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.Course{}).Where("id = ?", id).Update("is_deleted", 1).Error
}
