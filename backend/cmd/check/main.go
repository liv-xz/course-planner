package main

import (
	"context"
	"fmt"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/database"
	"sjqu-course-planner/internal/repository"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}
	db, err := database.Init(config.C.Database)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	type check struct {
		name string
		fn   func() (int64, error)
	}
	checks := []check{
		{"role", func() (int64, error) { _, n, e := repository.NewRoleRepo(db).List(ctx, 1, 1); return n, e }},
		{"user", func() (int64, error) { _, n, e := repository.NewUserRepo(db).List(ctx, 1, 1); return n, e }},
		{"major", func() (int64, error) { _, n, e := repository.NewMajorRepo(db).List(ctx, 1, 1); return n, e }},
		{"training_plan", func() (int64, error) { _, n, e := repository.NewTrainingPlanRepo(db).List(ctx, 1, 1); return n, e }},
		{"category", func() (int64, error) { _, n, e := repository.NewCategoryRepo(db).List(ctx, 1, 1); return n, e }},
		{"course", func() (int64, error) { _, n, e := repository.NewCourseRepo(db).List(ctx, 1, 1); return n, e }},
		{"course_category", func() (int64, error) {
			l, e := repository.NewCourseCategoryRepo(db).ListByCourse(ctx, 0)
			return int64(len(l)), e
		}},
		{"prerequisite", func() (int64, error) {
			l, e := repository.NewPrerequisiteRepo(db).ListByCourse(ctx, 0)
			return int64(len(l)), e
		}},
		{"course_offering", func() (int64, error) {
			l, e := repository.NewCourseOfferingRepo(db).ListByCourse(ctx, 0)
			return int64(len(l)), e
		}},
		{"student_record", func() (int64, error) {
			l, e := repository.NewStudentRecordRepo(db).ListByStudent(ctx, 0)
			return int64(len(l)), e
		}},
		{"credit_audit", func() (int64, error) {
			l, e := repository.NewCreditAuditRepo(db).ListByStudent(ctx, 0)
			return int64(len(l)), e
		}},
		{"chat_session", func() (int64, error) {
			l, e := repository.NewChatSessionRepo(db).ListByUser(ctx, 0)
			return int64(len(l)), e
		}},
		{"chat_message", func() (int64, error) {
			l, e := repository.NewChatMessageRepo(db).ListBySession(ctx, 0)
			return int64(len(l)), e
		}},
		{"security_log", func() (int64, error) { _, n, e := repository.NewSecurityLogRepo(db).List(ctx, 1, 1); return n, e }},
		{"knowledge_doc", func() (int64, error) { _, n, e := repository.NewKnowledgeDocRepo(db).List(ctx, 1, 1); return n, e }},
	}
	for _, c := range checks {
		n, err := c.fn()
		if err != nil {
			fmt.Printf("%-16s ERROR: %v\n", c.name, err)
			continue
		}
		fmt.Printf("%-16s OK (count=%d)\n", c.name, n)
	}
}
