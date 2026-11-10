package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

var (
	ErrStudentNotFound = errors.New("学生不存在")
	ErrUsernameExists  = errors.New("用户名已存在")
	ErrStudentNoExists = errors.New("学号已存在")
	ErrInvalidRole     = errors.New("角色配置错误")
	ErrInvalidRecord   = errors.New("修读记录参数不合法")
	ErrRecordNotFound  = errors.New("修读记录不存在")
)

var gradePattern = regexp.MustCompile(`^\d{4}$`)

// StudentAdminService 管理端学生与修读记录业务
type StudentAdminService struct {
	userRepo     *repository.UserRepo
	roleRepo     *repository.RoleRepo
	profileRepo  *repository.StudentProfileRepo
	majorRepo    *repository.MajorRepo
	recordRepo   *repository.StudentRecordRepo
	courseRepo   *repository.CourseRepo
	semesterRepo *repository.SemesterRepo
}

func NewStudentAdminService(userRepo *repository.UserRepo, roleRepo *repository.RoleRepo,
	profileRepo *repository.StudentProfileRepo, majorRepo *repository.MajorRepo,
	recordRepo *repository.StudentRecordRepo, courseRepo *repository.CourseRepo,
	semesterRepo *repository.SemesterRepo) *StudentAdminService {
	return &StudentAdminService{
		userRepo: userRepo, roleRepo: roleRepo, profileRepo: profileRepo, majorRepo: majorRepo,
		recordRepo: recordRepo, courseRepo: courseRepo, semesterRepo: semesterRepo,
	}
}

func (s *StudentAdminService) ListStudents(ctx context.Context, page, size int, keyword string) ([]repository.StudentRow, int64, error) {
	return s.profileRepo.ListWithUser(ctx, page, size, keyword)
}

func (s *StudentAdminService) CreateStudent(ctx context.Context, req dto.StudentCreateRequest) (*repository.StudentRow, error) {
	if !gradePattern.MatchString(req.Grade) {
		return nil, ErrInvalidRecord
	}
	if _, err := s.majorRepo.GetByID(ctx, req.MajorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMajorNotFound
		}
		return nil, err
	}
	if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, ErrUsernameExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if _, err := s.profileRepo.GetByStudentNo(ctx, req.StudentNo); err == nil {
		return nil, ErrStudentNoExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	role, err := s.roleRepo.GetByName(ctx, "student")
	if err != nil {
		return nil, ErrInvalidRole
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	name := req.Name
	status := req.Status
	if status != 0 {
		status = 1
	}
	u := &model.User{Username: req.Username, PasswordHash: string(hash), RoleID: role.ID, Name: &name, Status: status}
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}
	sp := &model.StudentProfile{UserID: u.ID, MajorID: req.MajorID, StudentNo: req.StudentNo, Grade: req.Grade}
	if err := s.profileRepo.Create(ctx, sp); err != nil {
		_ = s.userRepo.Delete(ctx, u.ID)
		return nil, err
	}
	return s.studentRow(ctx, u.ID)
}

func (s *StudentAdminService) UpdateStudent(ctx context.Context, userID int64, req dto.StudentUpdateRequest) (*repository.StudentRow, error) {
	if !gradePattern.MatchString(req.Grade) {
		return nil, ErrInvalidRecord
	}
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}
	sp, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}
	if _, err := s.majorRepo.GetByID(ctx, req.MajorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMajorNotFound
		}
		return nil, err
	}
	if old, err := s.profileRepo.GetByStudentNo(ctx, req.StudentNo); err == nil && old.UserID != userID {
		return nil, ErrStudentNoExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	name := req.Name
	u.Name = &name
	u.Status = req.Status
	if u.Status != 0 {
		u.Status = 1
	}
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}
	sp.MajorID = req.MajorID
	sp.StudentNo = req.StudentNo
	sp.Grade = req.Grade
	if err := s.profileRepo.Update(ctx, sp); err != nil {
		return nil, err
	}
	return s.studentRow(ctx, userID)
}

func (s *StudentAdminService) DeleteStudent(ctx context.Context, userID int64) error {
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrStudentNotFound
		}
		return err
	}
	return s.userRepo.Delete(ctx, userID)
}

func (s *StudentAdminService) studentRow(ctx context.Context, userID int64) (*repository.StudentRow, error) {
	row, err := s.profileRepo.GetRowByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}
	return row, nil
}

func (s *StudentAdminService) ListRecords(ctx context.Context, studentID int64) ([]dto.RecordItem, error) {
	if _, err := s.profileRepo.GetByUserID(ctx, studentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}
	records, err := s.recordRepo.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.RecordItem, 0, len(records))
	for _, r := range records {
		c, err := s.courseRepo.GetByIDIncludingDeleted(ctx, r.CourseID)
		if err != nil {
			return nil, err
		}
		credits := c.Credits
		if r.Credits != nil {
			credits = *r.Credits
		}
		out = append(out, dto.RecordItem{
			ID: r.ID, CourseID: r.CourseID, CourseCode: c.Code, CourseName: c.Name,
			Credits: credits, Score: r.Score, Semester: r.Semester, Status: r.Status,
		})
	}
	return out, nil
}

func (s *StudentAdminService) validateRecord(ctx context.Context, studentID, courseID int64, score *float64, semester *string, status int8, credits *float64) (*model.Course, error) {
	if _, err := s.profileRepo.GetByUserID(ctx, studentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}
	c, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}
	if status != 0 && status != 1 && status != 2 {
		return nil, ErrInvalidRecord
	}
	if score != nil && (*score < 0 || *score > 100) {
		return nil, ErrInvalidRecord
	}
	if credits != nil && *credits <= 0 {
		return nil, ErrInvalidRecord
	}
	if semester != nil && strings.TrimSpace(*semester) != "" {
		if _, err := s.semesterRepo.GetByCode(ctx, *semester); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrSemesterNotFound
			}
			return nil, err
		}
	}
	return c, nil
}

func (s *StudentAdminService) CreateRecord(ctx context.Context, studentID int64, req dto.RecordCreateRequest) (*dto.RecordItem, error) {
	if _, err := s.validateRecord(ctx, studentID, req.CourseID, req.Score, req.Semester, req.Status, req.Credits); err != nil {
		return nil, err
	}
	r := &model.StudentRecord{
		StudentID: studentID, CourseID: req.CourseID, Score: req.Score,
		Semester: req.Semester, Status: req.Status, Credits: req.Credits,
	}
	if err := s.recordRepo.Create(ctx, r); err != nil {
		return nil, err
	}
	list, err := s.ListRecords(ctx, studentID)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].ID == r.ID {
			return &list[i], nil
		}
	}
	return nil, ErrRecordNotFound
}

func (s *StudentAdminService) UpdateRecord(ctx context.Context, recordID int64, req dto.RecordUpdateRequest) (*dto.RecordItem, error) {
	r, err := s.recordRepo.GetByID(ctx, recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	if _, err := s.validateRecord(ctx, r.StudentID, req.CourseID, req.Score, req.Semester, req.Status, req.Credits); err != nil {
		return nil, err
	}
	r.CourseID = req.CourseID
	r.Score = req.Score
	r.Semester = req.Semester
	r.Status = req.Status
	r.Credits = req.Credits
	if err := s.recordRepo.Update(ctx, r); err != nil {
		return nil, err
	}
	list, err := s.ListRecords(ctx, r.StudentID)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].ID == r.ID {
			return &list[i], nil
		}
	}
	return nil, ErrRecordNotFound
}

func (s *StudentAdminService) DeleteRecord(ctx context.Context, recordID int64) error {
	if _, err := s.recordRepo.GetByID(ctx, recordID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	return s.recordRepo.Delete(ctx, recordID)
}
