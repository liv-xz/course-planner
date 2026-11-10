package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/repository"
)

var (
	ErrCategoryNotLeaf   = errors.New("课程只能挂到叶子类别")
	ErrCategoryOtherPlan = errors.New("类别不属于该培养方案")
	ErrCourseCategoryDup = errors.New("该课程在培养方案下已有类别")
	ErrCourseDeleted     = errors.New("课程已删除")
)

// CourseCategoryService 课程-类别归属业务
type CourseCategoryService struct {
	ccRepo       *repository.CourseCategoryRepo
	courseRepo   *repository.CourseRepo
	categoryRepo *repository.CategoryRepo
	planRepo     *repository.TrainingPlanRepo
}

func NewCourseCategoryService(ccRepo *repository.CourseCategoryRepo, courseRepo *repository.CourseRepo,
	categoryRepo *repository.CategoryRepo, planRepo *repository.TrainingPlanRepo) *CourseCategoryService {
	return &CourseCategoryService{ccRepo: ccRepo, courseRepo: courseRepo, categoryRepo: categoryRepo, planRepo: planRepo}
}

func (s *CourseCategoryService) ListByCourse(ctx context.Context, courseID int64) ([]repository.CourseCategoryInfo, error) {
	if _, err := s.courseRepo.GetByIDIncludingDeleted(ctx, courseID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}
	return s.ccRepo.ListByCourseWithInfo(ctx, courseID)
}

func (s *CourseCategoryService) SetForPlan(ctx context.Context, planID, courseID int64, categoryID *int64) error {
	if _, err := s.planRepo.GetByID(ctx, planID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPlanNotFound
		}
		return err
	}
	course, err := s.courseRepo.GetByIDIncludingDeleted(ctx, courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCourseNotFound
		}
		return err
	}
	if course.IsDeleted == 1 {
		return ErrCourseDeleted
	}
	if categoryID == nil {
		return s.ccRepo.ReplaceForPlan(ctx, courseID, planID, nil)
	}
	c, err := s.categoryRepo.GetByID(ctx, *categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}
	if c.PlanID != planID {
		return ErrCategoryOtherPlan
	}
	children, err := s.categoryRepo.ListByParent(ctx, c.ID)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return ErrCategoryNotLeaf
	}
	return s.ccRepo.ReplaceForPlan(ctx, courseID, planID, categoryID)
}
