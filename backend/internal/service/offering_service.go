package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

var (
	ErrOfferingNotFound = errors.New("开课信息不存在")
	ErrOfferingExists   = errors.New("该课程在该学期已存在开课记录")
)

// OfferingService 开课信息业务
type OfferingService struct {
	offeringRepo *repository.CourseOfferingRepo
	courseRepo   *repository.CourseRepo
	semesterRepo *repository.SemesterRepo
}

func NewOfferingService(offeringRepo *repository.CourseOfferingRepo, courseRepo *repository.CourseRepo,
	semesterRepo *repository.SemesterRepo) *OfferingService {
	return &OfferingService{offeringRepo: offeringRepo, courseRepo: courseRepo, semesterRepo: semesterRepo}
}

func (s *OfferingService) ListByCourse(ctx context.Context, courseID int64) ([]model.CourseOffering, error) {
	if _, err := s.courseRepo.GetByIDIncludingDeleted(ctx, courseID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}
	return s.offeringRepo.ListByCourse(ctx, courseID)
}

func (s *OfferingService) Create(ctx context.Context, courseID int64, req dto.OfferingCreateRequest) (*model.CourseOffering, error) {
	if _, err := s.courseRepo.GetByID(ctx, courseID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}
	if _, err := s.semesterRepo.GetByCode(ctx, req.Semester); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSemesterNotFound
		}
		return nil, err
	}
	if _, err := s.offeringRepo.GetByCourseSemester(ctx, courseID, req.Semester); err == nil {
		return nil, ErrOfferingExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	o := &model.CourseOffering{CourseID: courseID, Semester: req.Semester}
	if err := s.offeringRepo.Create(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OfferingService) Delete(ctx context.Context, id int64) error {
	if _, err := s.offeringRepo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOfferingNotFound
		}
		return err
	}
	return s.offeringRepo.Delete(ctx, id)
}
