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
	ErrPrerequisiteSelf     = errors.New("不能将课程自身设为先修")
	ErrPrerequisiteExists   = errors.New("先修关系已存在")
	ErrPrerequisiteCycle    = errors.New("先修关系存在循环，不能添加")
	ErrPrerequisiteNotFound = errors.New("先修关系不存在")
)

// PrerequisiteService 先修关系业务
type PrerequisiteService struct {
	prereqRepo *repository.PrerequisiteRepo
	courseRepo *repository.CourseRepo
}

func NewPrerequisiteService(prereqRepo *repository.PrerequisiteRepo, courseRepo *repository.CourseRepo) *PrerequisiteService {
	return &PrerequisiteService{prereqRepo: prereqRepo, courseRepo: courseRepo}
}

func (s *PrerequisiteService) ListByCourse(ctx context.Context, courseID int64) ([]model.Prerequisite, error) {
	if _, err := s.courseRepo.GetByIDIncludingDeleted(ctx, courseID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}
	return s.prereqRepo.ListByCourse(ctx, courseID)
}

// Create 新增 courseID ← prereqCourseID，含自关联校验与 DFS 循环检测
func (s *PrerequisiteService) Create(ctx context.Context, courseID int64, req dto.PrerequisiteCreateRequest) (*model.Prerequisite, error) {
	if courseID == req.PrereqCourseID {
		return nil, ErrPrerequisiteSelf
	}
	if _, err := s.courseRepo.GetByID(ctx, courseID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}
	if _, err := s.courseRepo.GetByID(ctx, req.PrereqCourseID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}

	existing, err := s.prereqRepo.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}
	for _, p := range existing {
		if p.PrereqCourseID == req.PrereqCourseID {
			return nil, ErrPrerequisiteExists
		}
	}

	// 若 prereqCourseID 能沿先修关系到达 courseID，则新增会成环
	cycle, err := s.canReach(ctx, req.PrereqCourseID, courseID)
	if err != nil {
		return nil, err
	}
	if cycle {
		return nil, ErrPrerequisiteCycle
	}

	p := &model.Prerequisite{CourseID: courseID, PrereqCourseID: req.PrereqCourseID}
	if err := s.prereqRepo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// canReach 从 start 出发沿"先修"关系 DFS，判断能否到达 target
func (s *PrerequisiteService) canReach(ctx context.Context, start, target int64) (bool, error) {
	visited := map[int64]bool{}
	var dfs func(cur int64) (bool, error)
	dfs = func(cur int64) (bool, error) {
		if cur == target {
			return true, nil
		}
		if visited[cur] {
			return false, nil
		}
		visited[cur] = true
		list, err := s.prereqRepo.ListByCourse(ctx, cur)
		if err != nil {
			return false, err
		}
		for _, p := range list {
			ok, err := dfs(p.PrereqCourseID)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	}
	return dfs(start)
}

func (s *PrerequisiteService) Delete(ctx context.Context, id int64) error {
	if _, err := s.prereqRepo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPrerequisiteNotFound
		}
		return err
	}
	return s.prereqRepo.Delete(ctx, id)
}
