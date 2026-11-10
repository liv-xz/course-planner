package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

var (
	ErrCourseNotFound   = errors.New("课程不存在")
	ErrCourseCodeExists = errors.New("课程代码已存在")
	ErrCourseInvalid    = errors.New("课程参数不合法")
)

var courseCodePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{2,32}$`)

// CourseService 课程库业务
type CourseService struct {
	courseRepo *repository.CourseRepo
	ccRepo     *repository.CourseCategoryRepo
}

func NewCourseService(courseRepo *repository.CourseRepo, ccRepo *repository.CourseCategoryRepo) *CourseService {
	return &CourseService{courseRepo: courseRepo, ccRepo: ccRepo}
}

func (s *CourseService) List(ctx context.Context, q dto.CourseListQuery) ([]model.Course, int64, error) {
	return s.courseRepo.Search(ctx, q.Page, q.Size, q.Keyword, q.CategoryID)
}

func (s *CourseService) Get(ctx context.Context, id int64) (*model.Course, error) {
	c, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}
	return c, nil
}

func validateCourse(code, name string, credits float64, nature int8) error {
	if !courseCodePattern.MatchString(strings.TrimSpace(code)) {
		return ErrCourseInvalid
	}
	return validateCourseFields(name, credits, nature)
}

func validateCourseFields(name string, credits float64, nature int8) error {
	if strings.TrimSpace(name) == "" || len([]rune(name)) > 128 || credits <= 0 || (nature != 1 && nature != 2) {
		return ErrCourseInvalid
	}
	return nil
}

func (s *CourseService) Create(ctx context.Context, req dto.CourseCreateRequest) (*model.Course, error) {
	if err := validateCourse(req.Code, req.Name, req.Credits, req.Nature); err != nil {
		return nil, err
	}
	if _, err := s.courseRepo.GetByCode(ctx, req.Code); err == nil {
		return nil, ErrCourseCodeExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if old, err := s.courseRepo.GetByCodeIncludingDeleted(ctx, req.Code); err == nil {
		old.Name = req.Name
		old.Credits = req.Credits
		old.Nature = req.Nature
		old.IsDeleted = 0
		if err := s.courseRepo.Update(ctx, old); err != nil {
			return nil, err
		}
		return old, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	c := &model.Course{Code: req.Code, Name: req.Name, Credits: req.Credits, Nature: req.Nature}
	if err := s.courseRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CourseService) Update(ctx context.Context, id int64, req dto.CourseUpdateRequest) (*model.Course, error) {
	if err := validateCourseFields(req.Name, req.Credits, req.Nature); err != nil {
		return nil, err
	}
	c, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Name = req.Name
	c.Credits = req.Credits
	c.Nature = req.Nature
	if err := s.courseRepo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// Delete 逻辑删除
func (s *CourseService) Delete(ctx context.Context, id int64) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.courseRepo.Delete(ctx, id)
}
