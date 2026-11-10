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
	ErrCategoryNotFound   = errors.New("课程类别不存在")
	ErrCategoryInUse      = errors.New("该类别下还有子类别、课程映射或审核数据，不能删除")
	ErrParentCategoryBad  = errors.New("父类别不存在或与当前类别不属于同一培养方案")
	ErrCategoryCycle      = errors.New("不能把类别移动到自己的子类别下")
	ErrCategoryHasCourses = errors.New("该类别下已有课程，不能作为父类别")
)

// CategoryService 课程类别业务
type CategoryService struct {
	categoryRepo *repository.CategoryRepo
	planRepo     *repository.TrainingPlanRepo
	ccRepo       *repository.CourseCategoryRepo
	auditRepo    *repository.CreditAuditRepo
}

func NewCategoryService(categoryRepo *repository.CategoryRepo, planRepo *repository.TrainingPlanRepo,
	ccRepo *repository.CourseCategoryRepo, auditRepo *repository.CreditAuditRepo) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo, planRepo: planRepo, ccRepo: ccRepo, auditRepo: auditRepo}
}

// ListTree 返回某培养方案下的类别树
func (s *CategoryService) ListTree(ctx context.Context, planID int64) ([]*dto.CategoryNode, error) {
	if _, err := s.planRepo.GetByID(ctx, planID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlanNotFound
		}
		return nil, err
	}
	list, err := s.categoryRepo.ListByPlan(ctx, planID)
	if err != nil {
		return nil, err
	}
	return buildCategoryTree(list), nil
}

func buildCategoryTree(list []model.Category) []*dto.CategoryNode {
	nodes := make(map[int64]*dto.CategoryNode, len(list))
	for _, c := range list {
		nodes[c.ID] = &dto.CategoryNode{
			ID: c.ID, PlanID: c.PlanID, ParentID: c.ParentID, Name: c.Name,
			RequiredCredits: c.RequiredCredits, Sort: c.Sort, Children: []*dto.CategoryNode{},
		}
	}
	roots := make([]*dto.CategoryNode, 0)
	for _, c := range list {
		n := nodes[c.ID]
		if c.ParentID != nil {
			if p, ok := nodes[*c.ParentID]; ok {
				p.Children = append(p.Children, n)
				continue
			}
		}
		roots = append(roots, n)
	}
	return roots
}

func (s *CategoryService) Create(ctx context.Context, req dto.CategoryCreateRequest) (*model.Category, error) {
	if _, err := s.planRepo.GetByID(ctx, req.PlanID); err != nil {
		return nil, ErrPlanNotFound
	}
	if req.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrParentCategoryBad
			}
			return nil, err
		}
		if parent.PlanID != req.PlanID {
			return nil, ErrParentCategoryBad
		}
		mappings, err := s.ccRepo.ListByCategory(ctx, parent.ID)
		if err != nil {
			return nil, err
		}
		if len(mappings) > 0 {
			return nil, ErrCategoryHasCourses
		}
	}
	c := &model.Category{
		PlanID: req.PlanID, ParentID: req.ParentID, Name: req.Name,
		RequiredCredits: req.RequiredCredits, Sort: req.Sort,
	}
	if err := s.categoryRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Update(ctx context.Context, id int64, req dto.CategoryUpdateRequest) (*model.Category, error) {
	c, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	if req.ParentID != nil {
		if *req.ParentID == id {
			return nil, ErrCategoryCycle
		}
		parent, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrParentCategoryBad
			}
			return nil, err
		}
		if parent.PlanID != c.PlanID {
			return nil, ErrParentCategoryBad
		}
		all, err := s.categoryRepo.ListByPlan(ctx, c.PlanID)
		if err != nil {
			return nil, err
		}
		byID := make(map[int64]*model.Category, len(all))
		for i := range all {
			byID[all[i].ID] = &all[i]
		}
		seen := map[int64]bool{}
		cur := &parent.ID
		for cur != nil {
			if *cur == id {
				return nil, ErrCategoryCycle
			}
			if seen[*cur] {
				return nil, ErrCategoryCycle
			}
			seen[*cur] = true
			node, ok := byID[*cur]
			if !ok {
				break
			}
			cur = node.ParentID
		}
		mappings, err := s.ccRepo.ListByCategory(ctx, parent.ID)
		if err != nil {
			return nil, err
		}
		if len(mappings) > 0 {
			return nil, ErrCategoryHasCourses
		}
	}
	c.ParentID = req.ParentID
	c.Name = req.Name
	c.RequiredCredits = req.RequiredCredits
	c.Sort = req.Sort
	if err := s.categoryRepo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

// Delete 有子类别/课程映射/审核记录时禁止删除
func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	if _, err := s.categoryRepo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}
	children, err := s.categoryRepo.ListByParent(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return ErrCategoryInUse
	}
	mappings, err := s.ccRepo.ListByCategory(ctx, id)
	if err != nil {
		return err
	}
	if len(mappings) > 0 {
		return ErrCategoryInUse
	}
	n, err := s.auditRepo.CountByCategory(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrCategoryInUse
	}
	return s.categoryRepo.Delete(ctx, id)
}
