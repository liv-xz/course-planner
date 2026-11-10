package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

var (
	ErrSemesterNotFound = errors.New("学期不存在")
	ErrSemesterExists   = errors.New("学期代码已存在")
	ErrSemesterInvalid  = errors.New("学期格式不正确，应为 2023-2024-1")
	ErrSemesterInUse    = errors.New("该学期已有开课或修读记录，不能删除")
)

var semesterCodePattern = regexp.MustCompile(`^(\d{4})-(\d{4})-([12])$`)

// SemesterService 学期字典业务
type SemesterService struct {
	repo *repository.SemesterRepo
}

func NewSemesterService(repo *repository.SemesterRepo) *SemesterService {
	return &SemesterService{repo: repo}
}

func (s *SemesterService) List(ctx context.Context, page, size int) ([]model.Semester, int64, error) {
	return s.repo.List(ctx, page, size)
}

func (s *SemesterService) All(ctx context.Context) ([]model.Semester, error) {
	list, _, err := s.repo.List(ctx, 1, 100)
	return list, err
}

func validateSemester(code, academicYear string, termNo int8, sequenceNo int) error {
	m := semesterCodePattern.FindStringSubmatch(code)
	if m == nil {
		return ErrSemesterInvalid
	}
	start, _ := strconv.Atoi(m[1])
	end, _ := strconv.Atoi(m[2])
	if end != start+1 || academicYear != fmt.Sprintf("%d-%d", start, end) {
		return ErrSemesterInvalid
	}
	if termNo != 1 && termNo != 2 {
		return ErrSemesterInvalid
	}
	if sequenceNo <= 0 {
		return ErrSemesterInvalid
	}
	return nil
}

func (s *SemesterService) Create(ctx context.Context, req dto.SemesterCreateRequest) (*model.Semester, error) {
	if err := validateSemester(req.Code, req.AcademicYear, req.TermNo, req.SequenceNo); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetByCode(ctx, req.Code); err == nil {
		return nil, ErrSemesterExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	m := &model.Semester{
		Code: req.Code, AcademicYear: req.AcademicYear, TermNo: req.TermNo,
		SequenceNo: req.SequenceNo, DisplayName: req.DisplayName,
	}
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *SemesterService) Update(ctx context.Context, id int64, req dto.SemesterUpdateRequest) (*model.Semester, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSemesterNotFound
		}
		return nil, err
	}
	if err := validateSemester(req.Code, req.AcademicYear, req.TermNo, req.SequenceNo); err != nil {
		return nil, err
	}
	if old, err := s.repo.GetByCode(ctx, req.Code); err == nil && old.ID != id {
		return nil, ErrSemesterExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if req.Code != m.Code {
		offerings, err := s.repo.CountOfferings(ctx, m.Code)
		if err != nil {
			return nil, err
		}
		records, err := s.repo.CountRecords(ctx, m.Code)
		if err != nil {
			return nil, err
		}
		if offerings > 0 || records > 0 {
			return nil, ErrSemesterInUse
		}
	}
	m.Code = req.Code
	m.AcademicYear = req.AcademicYear
	m.TermNo = req.TermNo
	m.SequenceNo = req.SequenceNo
	m.DisplayName = req.DisplayName
	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *SemesterService) Delete(ctx context.Context, id int64) error {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSemesterNotFound
		}
		return err
	}
	offerings, err := s.repo.CountOfferings(ctx, m.Code)
	if err != nil {
		return err
	}
	records, err := s.repo.CountRecords(ctx, m.Code)
	if err != nil {
		return err
	}
	if offerings > 0 || records > 0 {
		return ErrSemesterInUse
	}
	return s.repo.Delete(ctx, id)
}
