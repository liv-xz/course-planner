package main

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/database"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

// 临时种子程序：插入角色 + 一个 admin 测试用户（可重复执行）
func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}
	db, err := database.Init(config.C.Database)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	roleRepo := repository.NewRoleRepo(db)
	userRepo := repository.NewUserRepo(db)

	ensureRole := func(name, desc string) *model.Role {
		r, err := roleRepo.GetByName(ctx, name)
		if err == nil {
			fmt.Println("角色已存在:", name)
			return r
		}
		r = &model.Role{Name: name, Description: desc}
		if err := roleRepo.Create(ctx, r); err != nil {
			panic(err)
		}
		fmt.Println("角色已创建:", name, "id =", r.ID)
		return r
	}

	ensureRole("student", "学生")
	adminRole := ensureRole("admin", "教务管理员")

	if _, err := userRepo.GetByUsername(ctx, "admin"); err == nil {
		fmt.Println("用户已存在: admin")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	name := "教务管理员"
	u := &model.User{
		Username:     "admin",
		PasswordHash: string(hash),
		RoleID:       adminRole.ID,
		Name:         &name,
		Status:       1,
	}
	if err := userRepo.Create(ctx, u); err != nil {
		panic(err)
	}
	fmt.Println("用户已创建: admin / 123456, id =", u.ID)
}
