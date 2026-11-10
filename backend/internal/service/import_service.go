package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

// ImportService 修读记录导入
type ImportService struct {
	recordRepo   *repository.StudentRecordRepo
	courseRepo   *repository.CourseRepo
	semesterRepo *repository.SemesterRepo
	profileRepo  *repository.StudentProfileRepo
	cfg          config.Import
}

func NewImportService(recordRepo *repository.StudentRecordRepo, courseRepo *repository.CourseRepo,
	semesterRepo *repository.SemesterRepo, profileRepo *repository.StudentProfileRepo,
	cfg config.Import) *ImportService {
	return &ImportService{recordRepo: recordRepo, courseRepo: courseRepo, semesterRepo: semesterRepo, profileRepo: profileRepo, cfg: cfg}
}

func (s *ImportService) ImportRecords(ctx context.Context, studentID int64, filename string, r io.Reader) (*dto.ImportResult, error) {
	if _, err := s.profileRepo.GetByUserID(ctx, studentID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}
	ext := strings.ToLower(filepath.Ext(filename))
	allowed := false
	for _, a := range s.cfg.AllowedSuffixes {
		if strings.EqualFold(a, ext) {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, fmt.Errorf("只支持 %s 文件", strings.Join(s.cfg.AllowedSuffixes, "/"))
	}
	var rows [][]string
	var err error
	switch ext {
	case ".csv":
		rows, err = readCSV(r)
	case ".xlsx":
		rows, err = readXLSX(r)
	default:
		return nil, fmt.Errorf("不支持的文件类型")
	}
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("文件没有数据行")
	}
	maxRows := s.cfg.MaxRows
	if maxRows <= 0 {
		maxRows = 5000
	}
	if len(rows)-1 > maxRows {
		return nil, fmt.Errorf("单次最多导入 %d 行", maxRows)
	}
	header := map[string]int{}
	for i, h := range rows[0] {
		header[normalizeHeader(h)] = i
	}
	get := func(row []string, keys ...string) string {
		for _, k := range keys {
			if idx, ok := header[normalizeHeader(k)]; ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
		}
		return ""
	}
	result := &dto.ImportResult{Total: len(rows) - 1, Errors: []dto.ImportError{}}
	records := make([]model.StudentRecord, 0, len(rows)-1)
	seen := map[string]bool{}
	for i, row := range rows[1:] {
		rowNo := i + 2
		code := get(row, "course_code", "课程代码", "课程编号")
		name := get(row, "course_name", "课程名称", "课程名")
		if code == "" && name == "" {
			result.Errors = append(result.Errors, dto.ImportError{Row: rowNo, Reason: "课程代码和课程名称不能同时为空"})
			continue
		}
		course, err := s.findCourse(ctx, code, name)
		if err != nil {
			result.Errors = append(result.Errors, dto.ImportError{Row: rowNo, Reason: err.Error()})
			continue
		}
		semester := get(row, "semester", "学期")
		if semester == "" {
			result.Errors = append(result.Errors, dto.ImportError{Row: rowNo, Reason: "学期不能为空"})
			continue
		}
		if _, err := s.semesterRepo.GetByCode(ctx, semester); err != nil {
			result.Errors = append(result.Errors, dto.ImportError{Row: rowNo, Reason: "学期不存在: " + semester})
			continue
		}
		status, err := parseStatus(get(row, "status", "状态"))
		if err != nil {
			result.Errors = append(result.Errors, dto.ImportError{Row: rowNo, Reason: err.Error()})
			continue
		}
		var score *float64
		if v := get(row, "score", "成绩", "分数"); v != "" {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil || f < 0 || f > 100 {
				result.Errors = append(result.Errors, dto.ImportError{Row: rowNo, Reason: "成绩必须在 0~100 之间"})
				continue
			}
			score = &f
		}
		var credits *float64
		if v := get(row, "credits", "学分"); v != "" {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil || f <= 0 {
				result.Errors = append(result.Errors, dto.ImportError{Row: rowNo, Reason: "学分必须大于 0"})
				continue
			}
			credits = &f
		}
		key := fmt.Sprintf("%d|%s", course.ID, semester)
		if seen[key] {
			result.Errors = append(result.Errors, dto.ImportError{Row: rowNo, Reason: "文件中同一课程同一学期重复"})
			continue
		}
		seen[key] = true
		sem := semester
		records = append(records, model.StudentRecord{
			StudentID: studentID, CourseID: course.ID, Score: score,
			Semester: &sem, Status: status, Credits: credits,
		})
	}
	if len(result.Errors) > 0 {
		result.Failed = len(result.Errors)
		result.Success = 0
		return result, nil
	}
	if err := s.recordRepo.UpsertBatch(ctx, records); err != nil {
		return nil, err
	}
	result.Success = len(records)
	return result, nil
}

func (s *ImportService) findCourse(ctx context.Context, code, name string) (*model.Course, error) {
	if code != "" {
		if c, err := s.courseRepo.GetByCode(ctx, code); err == nil {
			return c, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if name != "" {
		list, _, err := s.courseRepo.Search(ctx, 1, 100, name, 0)
		if err != nil {
			return nil, err
		}
		var matches []model.Course
		for _, c := range list {
			if c.Name == name {
				matches = append(matches, c)
			}
		}
		if len(matches) == 1 {
			return &matches[0], nil
		}
		if len(matches) > 1 {
			return nil, fmt.Errorf("课程名称匹配到多门课程，请填写课程代码")
		}
	}
	return nil, fmt.Errorf("课程不存在: %s%s", code, name)
}

func normalizeHeader(s string) string {
	s = strings.TrimSpace(strings.TrimPrefix(s, "\ufeff"))
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "_", "")
	return s
}

func parseStatus(s string) (int8, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "", "1", "通过", "pass", "passed":
		return 1, nil
	case "0", "在读", "studying":
		return 0, nil
	case "2", "未通过", "fail", "failed":
		return 2, nil
	default:
		return 0, fmt.Errorf("状态只能是 0/1/2 或在读/通过/未通过")
	}
}

func readCSV(r io.Reader) ([][]string, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true
	return cr.ReadAll()
}

func readXLSX(r io.Reader) ([][]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel 没有工作表")
	}
	return f.GetRows(sheets[0])
}
