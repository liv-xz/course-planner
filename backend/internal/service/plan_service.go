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
	ErrPlanNotFound  = errors.New("培养方案不存在")
	ErrMajorNotFound = errors.New("专业不存在")
	ErrPlanInUse     = errors.New("该培养方案下还有类别或审核数据，不能删除")
	ErrPlanExists    = errors.New("该专业该年级的培养方案已存在")
)

// PlanService 培养方案业务
type PlanService struct {
	planRepo      *repository.TrainingPlanRepo
	categoryRepo  *repository.CategoryRepo
	auditRepo     *repository.CreditAuditRepo
	knowledgeRepo *repository.KnowledgeDocRepo
	majorRepo     *repository.MajorRepo
}

func NewPlanService(planRepo *repository.TrainingPlanRepo, categoryRepo *repository.CategoryRepo,
	auditRepo *repository.CreditAuditRepo, knowledgeRepo *repository.KnowledgeDocRepo,
	majorRepo *repository.MajorRepo) *PlanService {
	return &PlanService{
		planRepo: planRepo, categoryRepo: categoryRepo, auditRepo: auditRepo,
		knowledgeRepo: knowledgeRepo, majorRepo: majorRepo,
	}
}

func (s *PlanService) List(ctx context.Context, page, size int) ([]model.TrainingPlan, int64, error) {
	return s.planRepo.List(ctx, page, size)
}

func (s *PlanService) Get(ctx context.Context, id int64) (*model.TrainingPlan, error) {
	p, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlanNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *PlanService) Create(ctx context.Context, req dto.PlanCreateRequest) (*model.TrainingPlan, error) {
	if _, err := s.majorRepo.GetByID(ctx, req.MajorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMajorNotFound
		}
		return nil, err
	}
	if _, err := s.planRepo.GetByMajorGrade(ctx, req.MajorID, req.Grade); err == nil {
		return nil, ErrPlanExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	p := &model.TrainingPlan{
		MajorID:      req.MajorID,
		Grade:        req.Grade,
		Name:         req.Name,
		TotalCredits: req.TotalCredits,
	}
	if err := s.planRepo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PlanService) Update(ctx context.Context, id int64, req dto.PlanUpdateRequest) (*model.TrainingPlan, error) {
	p, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing, err := s.planRepo.GetByMajorGrade(ctx, p.MajorID, req.Grade); err == nil && existing.ID != id {
		return nil, ErrPlanExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	p.Grade = req.Grade
	p.Name = req.Name
	p.TotalCredits = req.TotalCredits
	if err := s.planRepo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Delete 有引用（类别/审核/知识文档）时禁止删除
func (s *PlanService) Delete(ctx context.Context, id int64) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	cats, err := s.categoryRepo.ListByPlan(ctx, id)
	if err != nil {
		return err
	}
	if len(cats) > 0 {
		return ErrPlanInUse
	}
	n, err := s.auditRepo.CountByPlan(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrPlanInUse
	}
	docs, err := s.knowledgeRepo.ListByPlan(ctx, id)
	if err != nil {
		return err
	}
	if len(docs) > 0 {
		return ErrPlanInUse
	}
	return s.planRepo.Delete(ctx, id)
}
