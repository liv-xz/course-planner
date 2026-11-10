package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"sjqu-course-planner/internal/config"
)

// IsConfigured 判断是否配置了数据库（host/user/name 均为空则视为未配置）
func IsConfigured(cfg config.Database) bool {
	return cfg.Host != "" && cfg.User != "" && cfg.Name != ""
}

// DSN 生成 MySQL 数据源字符串
func DSN(cfg config.Database) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.Charset)
}

// Init 初始化 GORM 连接（含连接池与 ping 检查）
func Init(cfg config.Database) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(DSN(cfg)), &gorm.Config{
		Logger:         gormlogger.Default.LogMode(gormlogger.Info),
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层连接失败: %w", err)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库 ping 失败: %w", err)
	}
	return db, nil
}
