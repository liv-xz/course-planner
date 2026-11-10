package test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/logger"
	"sjqu-course-planner/internal/model"
)

// 测试库（与开发库 course_planner 分开，测试结束自动删除）
const testDBName = "course_planner_test"

var (
	testDB   *gorm.DB
	adminDSN string
	dbCfg    config.Database
	testUser = "admin"
	testPass = "123456"
)

func TestMain(m *testing.M) {
	// 1. 加载配置（测试包工作目录为 backend/test）
	os.Setenv("CONFIG_PATH", "../configs/config.yaml")
	if err := config.Load(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	dbCfg = config.C.Database
	adminDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port)

	// 初始化日志（测试不写文件、只留 warn 以上，减少噪音；中间件依赖 logger.Log）
	logCfg := config.C.Logger
	logCfg.File = ""
	logCfg.Level = "warn"
	if err := logger.Init(logCfg); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 2. 每次运行重建测试库
	if err := recreateTestDB(); err != nil {
		log.Fatalf("创建测试库失败: %v", err)
	}

	// 3. 连接测试库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, testDBName)
	var err error
	testDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:         gormlogger.Default.LogMode(gormlogger.Silent),
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		log.Fatalf("连接测试库失败: %v", err)
	}

	// 4. 建表 + 种子数据
	if err := migrate(); err != nil {
		log.Fatalf("迁移失败: %v", err)
	}
	if err := seed(); err != nil {
		log.Fatalf("种子失败: %v", err)
	}

	// 5. 跑测试
	code := m.Run()

	// 6. 清理：关闭连接并删除测试库
	if sqlDB, e := testDB.DB(); e == nil {
		_ = sqlDB.Close()
	}
	dropTestDB()
	os.Exit(code)
}

func recreateTestDB() error {
	admin, err := sql.Open("mysql", adminDSN)
	if err != nil {
		return err
	}
	defer admin.Close()
	if _, err := admin.Exec("DROP DATABASE IF EXISTS " + testDBName); err != nil {
		return err
	}
	_, err = admin.Exec("CREATE DATABASE " + testDBName + " DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	return err
}

func dropTestDB() {
	admin, err := sql.Open("mysql", adminDSN)
	if err != nil {
		log.Println("清理测试库失败:", err)
		return
	}
	defer admin.Close()
	if _, err := admin.Exec("DROP DATABASE IF EXISTS " + testDBName); err != nil {
		log.Println("清理测试库失败:", err)
	}
}

func migrate() error {
	return testDB.AutoMigrate(
		&model.Role{}, &model.Major{}, &model.TrainingPlan{}, &model.User{},
		&model.StudentProfile{},
		&model.Category{}, &model.Course{}, &model.CourseCategory{},
		&model.Prerequisite{}, &model.CourseOffering{}, &model.StudentRecord{}, &model.Semester{},
		&model.CreditAudit{}, &model.ChatSession{}, &model.ChatMessage{},
		&model.SecurityLog{}, &model.KnowledgeDoc{}, &model.KnowledgeChunk{},
	)
}

func seed() error {
	roles := []model.Role{
		{Name: "student", Description: "学生"},
		{Name: "admin", Description: "教务管理员"},
	}
	if err := testDB.Create(&roles).Error; err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.MinCost)
	if err != nil {
		return err
	}
	name := "教务管理员"
	admin := model.User{
		Username:     testUser,
		PasswordHash: string(hash),
		RoleID:       roles[1].ID,
		Name:         &name,
		Status:       1,
	}
	return testDB.Create(&admin).Error
}
