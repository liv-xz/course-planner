package test

import (
	"context"
	"testing"

	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

// TestRoleRepo_CRUD 角色增删查改
func TestRoleRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewRoleRepo(testDB)

	r := &model.Role{Name: "test_role", Description: "测试角色"}
	if err := repo.Create(ctx, r); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	t.Cleanup(func() { testDB.Exec("DELETE FROM role WHERE id = ?", r.ID) })

	got, err := repo.GetByID(ctx, r.ID)
	if err != nil || got.Name != "test_role" {
		t.Fatalf("查询失败: %v", err)
	}

	got.Description = "修改后"
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("更新失败: %v", err)
	}

	if err := repo.Delete(ctx, r.ID); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
	if _, err := repo.GetByID(ctx, r.ID); err == nil {
		t.Fatal("删除后仍能查到")
	}
}

// TestUserRepo_SoftDelete 用户逻辑删除
func TestUserRepo_SoftDelete(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewUserRepo(testDB)

	name := "测试用户"
	u := &model.User{
		Username:     "test_user_crud",
		PasswordHash: "x",
		RoleID:       1,
		Name:         &name,
		Status:       1,
	}
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	t.Cleanup(func() { testDB.Exec("DELETE FROM user WHERE id = ?", u.ID) })

	if _, err := repo.GetByUsername(ctx, "test_user_crud"); err != nil {
		t.Fatalf("按用户名查询失败: %v", err)
	}

	if err := repo.Delete(ctx, u.ID); err != nil {
		t.Fatalf("逻辑删除失败: %v", err)
	}
	if _, err := repo.GetByID(ctx, u.ID); err == nil {
		t.Fatal("逻辑删除后仍能查到")
	}
}

// TestPrerequisiteRepo 先修关系查询
func TestPrerequisiteRepo(t *testing.T) {
	ctx := context.Background()

	c1 := &model.Course{Code: "T101", Name: "课A", Credits: 2}
	c2 := &model.Course{Code: "T102", Name: "课B", Credits: 2}
	if err := testDB.Create(c1).Error; err != nil {
		t.Fatalf("创建课程失败: %v", err)
	}
	if err := testDB.Create(c2).Error; err != nil {
		t.Fatalf("创建课程失败: %v", err)
	}
	t.Cleanup(func() {
		testDB.Exec("DELETE FROM prerequisite WHERE course_id IN (?,?)", c1.ID, c2.ID)
		testDB.Exec("DELETE FROM course WHERE id IN (?,?)", c1.ID, c2.ID)
	})

	repo := repository.NewPrerequisiteRepo(testDB)
	p := &model.Prerequisite{CourseID: c2.ID, PrereqCourseID: c1.ID}
	if err := repo.Create(ctx, p); err != nil {
		t.Fatalf("创建先修关系失败: %v", err)
	}

	list, err := repo.ListByCourse(ctx, c2.ID)
	if err != nil {
		t.Fatalf("查询先修关系失败: %v", err)
	}
	if len(list) != 1 || list[0].PrereqCourseID != c1.ID {
		t.Fatalf("先修关系不符: %+v", list)
	}
}
