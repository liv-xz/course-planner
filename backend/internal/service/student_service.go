package service

import (
	"context"
	"errors"
	"sort"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

var ErrStudentProfileNotFound = errors.New("学生档案不存在")

// StudentService 学生端业务
type StudentService struct {
	profileRepo  *repository.StudentProfileRepo
	planRepo     *repository.TrainingPlanRepo
	categoryRepo *repository.CategoryRepo
	courseRepo   *repository.CourseRepo
	ccRepo       *repository.CourseCategoryRepo
	prereqRepo   *repository.PrerequisiteRepo
	offeringRepo *repository.CourseOfferingRepo
	recordRepo   *repository.StudentRecordRepo
	semesterRepo *repository.SemesterRepo
}

func NewStudentService(profileRepo *repository.StudentProfileRepo, planRepo *repository.TrainingPlanRepo,
	categoryRepo *repository.CategoryRepo, courseRepo *repository.CourseRepo,
	ccRepo *repository.CourseCategoryRepo, prereqRepo *repository.PrerequisiteRepo,
	offeringRepo *repository.CourseOfferingRepo, recordRepo *repository.StudentRecordRepo,
	semesterRepo *repository.SemesterRepo) *StudentService {
	return &StudentService{
		profileRepo: profileRepo, planRepo: planRepo, categoryRepo: categoryRepo,
		courseRepo: courseRepo, ccRepo: ccRepo, prereqRepo: prereqRepo,
		offeringRepo: offeringRepo, recordRepo: recordRepo, semesterRepo: semesterRepo,
	}
}

// MyRecords 我的修读记录
func (s *StudentService) MyRecords(ctx context.Context, userID int64) ([]dto.RecordItem, error) {
	records, err := s.recordRepo.ListByStudent(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.RecordItem, 0, len(records))
	for _, r := range records {
		item := dto.RecordItem{
			ID: r.ID, CourseID: r.CourseID, Score: r.Score, Semester: r.Semester, Status: r.Status,
		}
		credits := 0.0
		c, err := s.courseRepo.GetByIDIncludingDeleted(ctx, r.CourseID)
		if err != nil {
			return nil, err
		}
		item.CourseCode = c.Code
		item.CourseName = c.Name
		credits = c.Credits
		if r.Credits != nil {
			credits = *r.Credits
		}
		item.Credits = credits
		out = append(out, item)
	}
	return out, nil
}

// MyAudit 我的学分审核（实时计算，不落库）
func (s *StudentService) MyAudit(ctx context.Context, userID int64) (*dto.AuditResponse, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, ErrStudentProfileNotFound
	}
	plan, err := s.planRepo.GetByMajorGrade(ctx, profile.MajorID, profile.Grade)
	if err != nil {
		return nil, ErrPlanNotFound
	}
	categories, err := s.categoryRepo.ListByPlan(ctx, plan.ID)
	if err != nil {
		return nil, err
	}
	records, err := s.recordRepo.ListByStudent(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 同一门课只保留一条有效记录：通过记录中取最高分，同分取最新
	effective := map[int64]model.StudentRecord{}
	for _, r := range records {
		if r.Status != 1 {
			continue
		}
		old, ok := effective[r.CourseID]
		if !ok || recordBetter(r, old) {
			effective[r.CourseID] = r
		}
	}

	direct := map[int64]float64{}
	for _, r := range effective {
		credits := 0.0
		if r.Credits != nil {
			credits = *r.Credits
		}
		c, err := s.courseRepo.GetByIDIncludingDeleted(ctx, r.CourseID)
		if err != nil {
			return nil, err
		}
		if r.Credits == nil {
			credits = c.Credits
		}
		ccs, err := s.ccRepo.ListByCourse(ctx, r.CourseID)
		if err != nil {
			return nil, err
		}
		for _, cc := range ccs {
			direct[cc.CategoryID] += credits
		}
	}

	return &dto.AuditResponse{
		PlanID:     plan.ID,
		PlanName:   plan.Name,
		Grade:      profile.Grade,
		Categories: buildAuditTree(categories, direct),
	}, nil
}

func recordBetter(a, b model.StudentRecord) bool {
	if a.Score == nil && b.Score == nil {
		return a.ID > b.ID
	}
	if a.Score == nil {
		return false
	}
	if b.Score == nil {
		return true
	}
	if *a.Score != *b.Score {
		return *a.Score > *b.Score
	}
	return a.ID > b.ID
}

func buildAuditTree(categories []model.Category, direct map[int64]float64) []*dto.AuditCategory {
	nodes := make(map[int64]*dto.AuditCategory, len(categories))
	for _, c := range categories {
		nodes[c.ID] = &dto.AuditCategory{
			ID: c.ID, Name: c.Name, RequiredCredits: c.RequiredCredits,
			Children: []*dto.AuditCategory{},
		}
	}
	roots := make([]*dto.AuditCategory, 0)
	for _, c := range categories {
		n := nodes[c.ID]
		if c.ParentID != nil {
			if p, ok := nodes[*c.ParentID]; ok {
				p.Children = append(p.Children, n)
				continue
			}
		}
		roots = append(roots, n)
	}

	var calc func(n *dto.AuditCategory) float64
	calc = func(n *dto.AuditCategory) float64 {
		total := direct[n.ID]
		for _, ch := range n.Children {
			total += calc(ch)
		}
		n.CompletedCredits = total
		gap := n.RequiredCredits - total
		if gap < 0 {
			gap = 0
		}
		n.GapCredits = gap
		return total
	}
	for _, r := range roots {
		calc(r)
	}
	return roots
}

// MyRecommendations 选课推荐（可选 semester 过滤开课）
func (s *StudentService) MyRecommendations(ctx context.Context, userID int64, semester string) ([]dto.RecommendationItem, error) {
	if semester != "" {
		if _, err := s.semesterRepo.GetByCode(ctx, semester); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrSemesterNotFound
			}
			return nil, err
		}
	}
	audit, err := s.MyAudit(ctx, userID)
	if err != nil {
		return nil, err
	}

	records, err := s.recordRepo.ListByStudent(ctx, userID)
	if err != nil {
		return nil, err
	}
	passedSet := map[int64]bool{}
	takenSet := map[int64]bool{}
	for _, r := range records {
		switch r.Status {
		case 1:
			passedSet[r.CourseID] = true
		case 0:
			takenSet[r.CourseID] = true
		}
	}

	leaves := make([]*dto.AuditCategory, 0)
	collectGapLeaves(audit.Categories, &leaves)

	out := make([]dto.RecommendationItem, 0)
	seen := map[int64]bool{}
	for _, cat := range leaves {
		courses, err := s.courseRepo.ListByCategory(ctx, cat.ID)
		if err != nil {
			return nil, err
		}
		for _, c := range courses {
			if seen[c.ID] || passedSet[c.ID] || takenSet[c.ID] {
				continue
			}
			ok, err := s.prereqMet(ctx, c.ID, passedSet)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
			offerings, err := s.offeringRepo.ListByCourse(ctx, c.ID)
			if err != nil {
				return nil, err
			}
			if len(offerings) == 0 {
				continue
			}
			semesterSet := map[string]bool{}
			offered := false
			for _, o := range offerings {
				if !semesterSet[o.Semester] {
					semesterSet[o.Semester] = true
				}
				if semester != "" && o.Semester == semester {
					offered = true
				}
			}
			if semester != "" && !offered {
				continue
			}
			semesters := make([]string, 0, len(semesterSet))
			for code := range semesterSet {
				semesters = append(semesters, code)
			}
			sort.Strings(semesters)
			seen[c.ID] = true
			out = append(out, dto.RecommendationItem{
				CourseID: c.ID, CourseCode: c.Code, CourseName: c.Name, Credits: c.Credits,
				CategoryID: cat.ID, CategoryName: cat.Name, Semesters: semesters,
				Reason: "补足「" + cat.Name + "」的学分缺口",
			})
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Credits != out[j].Credits {
			return out[i].Credits > out[j].Credits
		}
		return out[i].CourseCode < out[j].CourseCode
	})
	return out, nil
}

func collectGapLeaves(nodes []*dto.AuditCategory, out *[]*dto.AuditCategory) {
	for _, n := range nodes {
		if len(n.Children) > 0 {
			collectGapLeaves(n.Children, out)
			continue
		}
		if n.GapCredits > 0 {
			*out = append(*out, n)
		}
	}
}

// prereqMet 该课程的所有先修是否都已通过
func (s *StudentService) prereqMet(ctx context.Context, courseID int64, passedSet map[int64]bool) (bool, error) {
	prereqs, err := s.prereqRepo.ListByCourse(ctx, courseID)
	if err != nil {
		return false, err
	}
	for _, p := range prereqs {
		if !passedSet[p.PrereqCourseID] {
			return false, nil
		}
	}
	return true, nil
}
