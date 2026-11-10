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
	ErrMajorExists = errors.New("专业代码已存在")
	ErrMajorInUse  = errors.New("该专业下还有培养方案或学生，不能删除")
)

// MajorService 专业业务
type MajorService struct {
	repo     *repository.MajorRepo
	planRepo *repository.TrainingPlanRepo
	spRepo   *repository.StudentProfileRepo
}

func NewMajorService(repo *repository.MajorRepo, planRepo *repository.TrainingPlanRepo, spRepo *repository.StudentProfileRepo) *MajorService {
	return &MajorService{repo: repo, planRepo: planRepo, spRepo: spRepo}
}

func (s *MajorService) List(ctx context.Context, page, size int) ([]model.Major, int64, error) {
	return s.repo.List(ctx, page, size)
}

func (s *MajorService) Create(ctx context.Context, req dto.MajorCreateRequest) (*model.Major, error) {
	if _, err := s.repo.GetByCode(ctx, req.Code); err == nil {
		return nil, ErrMajorExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	m := &model.Major{Code: req.Code, Name: req.Name}
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MajorService) Update(ctx context.Context, id int64, req dto.MajorUpdateRequest) (*model.Major, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMajorNotFound
		}
		return nil, err
	}
	if old, err := s.repo.GetByCode(ctx, req.Code); err == nil && old.ID != id {
		return nil, ErrMajorExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	m.Code = req.Code
	m.Name = req.Name
	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MajorService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMajorNotFound
		}
		return err
	}
	n, err := s.planRepo.CountByMajor(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrMajorInUse
	}
	n, err = s.spRepo.CountByMajor(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrMajorInUse
	}
	return s.repo.Delete(ctx, id)
}
