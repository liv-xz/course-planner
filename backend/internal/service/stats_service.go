package service

import (
	"context"
	"sort"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/repository"
)

// StatsService 管理端统计
type StatsService struct {
	profileRepo  *repository.StudentProfileRepo
	planRepo     *repository.TrainingPlanRepo
	courseRepo   *repository.CourseRepo
	majorRepo    *repository.MajorRepo
	semesterRepo *repository.SemesterRepo
	categoryRepo *repository.CategoryRepo
	studentSvc   *StudentService
}

func NewStatsService(profileRepo *repository.StudentProfileRepo, planRepo *repository.TrainingPlanRepo,
	courseRepo *repository.CourseRepo, majorRepo *repository.MajorRepo,
	semesterRepo *repository.SemesterRepo, categoryRepo *repository.CategoryRepo,
	studentSvc *StudentService) *StatsService {
	return &StatsService{
		profileRepo: profileRepo, planRepo: planRepo, courseRepo: courseRepo,
		majorRepo: majorRepo, semesterRepo: semesterRepo, categoryRepo: categoryRepo,
		studentSvc: studentSvc,
	}
}

func (s *StatsService) Overview(ctx context.Context) (*dto.OverviewStats, error) {
	_, students, err := s.profileRepo.ListWithUser(ctx, 1, 1, "")
	if err != nil {
		return nil, err
	}
	_, plans, err := s.planRepo.List(ctx, 1, 1)
	if err != nil {
		return nil, err
	}
	_, courses, err := s.courseRepo.List(ctx, 1, 1)
	if err != nil {
		return nil, err
	}
	_, majors, err := s.majorRepo.List(ctx, 1, 1)
	if err != nil {
		return nil, err
	}
	_, semesters, err := s.semesterRepo.List(ctx, 1, 1)
	if err != nil {
		return nil, err
	}
	planList, _, err := s.planRepo.List(ctx, 1, 100)
	if err != nil {
		return nil, err
	}
	var categories int64
	for _, p := range planList {
		n, err := s.categoryRepo.CountByPlan(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		categories += n
	}
	return &dto.OverviewStats{
		Students: students, Plans: plans, Courses: courses, Majors: majors,
		Semesters: semesters, Categories: categories,
	}, nil
}

func (s *StatsService) CreditGap(ctx context.Context) ([]dto.CategoryGapStat, error) {
	rows, err := s.allStudents(ctx)
	if err != nil {
		return nil, err
	}
	type key struct {
		planID int64
		catID  int64
	}
	agg := map[key]*dto.CategoryGapStat{}
	for _, row := range rows {
		audit, err := s.studentSvc.MyAudit(ctx, row.UserID)
		if err != nil {
			continue
		}
		walkAudit(audit.Categories, func(c *dto.AuditCategory) {
			if len(c.Children) > 0 || c.RequiredCredits <= 0 {
				return
			}
			k := key{planID: audit.PlanID, catID: c.ID}
			item, ok := agg[k]
			if !ok {
				item = &dto.CategoryGapStat{
					PlanID: audit.PlanID, PlanName: audit.PlanName,
					CategoryID: c.ID, CategoryName: c.Name,
					RequiredCredits: c.RequiredCredits,
				}
				agg[k] = item
			}
			item.StudentCount++
			if c.GapCredits > 0 {
				item.GapStudentCount++
				item.TotalGap += c.GapCredits
			}
		})
	}
	out := make([]dto.CategoryGapStat, 0, len(agg))
	for _, item := range agg {
		if item.StudentCount > 0 {
			item.AverageGap = round1(item.TotalGap / float64(item.StudentCount))
		}
		out = append(out, *item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].PlanID != out[j].PlanID {
			return out[i].PlanID < out[j].PlanID
		}
		return out[i].CategoryName < out[j].CategoryName
	})
	return out, nil
}

func (s *StatsService) allStudents(ctx context.Context) ([]repository.StudentRow, error) {
	var out []repository.StudentRow
	page := 1
	for {
		list, total, err := s.profileRepo.ListWithUser(ctx, page, 100, "")
		if err != nil {
			return nil, err
		}
		out = append(out, list...)
		if len(out) >= int(total) || len(list) == 0 {
			break
		}
		page++
	}
	return out, nil
}

func walkAudit(nodes []*dto.AuditCategory, fn func(*dto.AuditCategory)) {
	for _, n := range nodes {
		fn(n)
		walkAudit(n.Children, fn)
	}
}

func round1(v float64) float64 {
	return float64(int(v*10+0.5)) / 10
}
