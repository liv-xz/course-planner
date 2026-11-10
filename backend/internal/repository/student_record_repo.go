package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/model"
)

// StudentRecordRepo 学生修读记录数据访问
type StudentRecordRepo struct{ db *gorm.DB }

func NewStudentRecordRepo(db *gorm.DB) *StudentRecordRepo { return &StudentRecordRepo{db: db} }

func (r *StudentRecordRepo) Create(ctx context.Context, m *model.StudentRecord) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// BatchCreate 批量导入
func (r *StudentRecordRepo) BatchCreate(ctx context.Context, list []model.StudentRecord) error {
	if len(list) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(list, 200).Error
}

func (r *StudentRecordRepo) GetByID(ctx context.Context, id int64) (*model.StudentRecord, error) {
	var m model.StudentRecord
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// ListByStudent 某学生的全部修读记录
func (r *StudentRecordRepo) ListByStudent(ctx context.Context, studentID int64) ([]model.StudentRecord, error) {
	var list []model.StudentRecord
	err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Find(&list).Error
	return list, err
}

// GetByStudentCourse 某学生某门课的记录
func (r *StudentRecordRepo) GetByStudentCourse(ctx context.Context, studentID, courseID int64) (*model.StudentRecord, error) {
	var m model.StudentRecord
	if err := r.db.WithContext(ctx).
		Where("student_id = ? AND course_id = ?", studentID, courseID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *StudentRecordRepo) Update(ctx context.Context, m *model.StudentRecord) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// UpsertBatch 按 student_id + course_id + semester 幂等写入
func (r *StudentRecordRepo) UpsertBatch(ctx context.Context, list []model.StudentRecord) error {
	if len(list) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range list {
			m := list[i]
			var existing model.StudentRecord
			q := tx.Where("student_id = ? AND course_id = ?", m.StudentID, m.CourseID)
			if m.Semester != nil {
				q = q.Where("semester = ?", *m.Semester)
			} else {
				q = q.Where("semester IS NULL")
			}
			err := q.First(&existing).Error
			if err == nil {
				existing.Score = m.Score
				existing.Status = m.Status
				existing.Credits = m.Credits
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err := tx.Create(&m).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *StudentRecordRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.StudentRecord{}, id).Error
}
