package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/database"
	"sjqu-course-planner/internal/logger"
	"sjqu-course-planner/internal/router"
)

func main() {
	// 1. 加载配置
	if err := config.Load(); err != nil {
		fmt.Println("加载配置失败:", err)
		os.Exit(1)
	}
	cfg := config.C

	// 2. 初始化日志
	if err := logger.Init(cfg.Logger); err != nil {
		fmt.Println("初始化日志失败:", err)
		os.Exit(1)
	}
	defer logger.Sync()
	logger.Log.Info("应用启动",
		zap.String("name", cfg.App.Name),
		zap.String("env", cfg.App.Env),
	)

	// 3. 连接数据库（未配置则跳过）
	var db *gorm.DB
	if database.IsConfigured(cfg.Database) {
		var err error
		db, err = database.Init(cfg.Database)
		if err != nil {
			logger.Log.Fatal("数据库连接失败", zap.Error(err))
		}
		logger.Log.Info("数据库连接成功")
	} else {
		logger.Log.Warn("数据库未配置，跳过连接（可在 configs/config.yaml 填写 database 项）")
	}

	// 4. 构建路由
	r := router.Setup(cfg, db)

	// 5. 启动 HTTP 服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: r,
	}
	go func() {
		logger.Log.Info("HTTP 服务启动", zap.Int("port", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal("HTTP 服务异常", zap.Error(err))
		}
	}()

	// 6. 等待退出信号，优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("收到退出信号，开始优雅关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("优雅关闭失败", zap.Error(err))
	}
	logger.Log.Info("服务已退出")
}
