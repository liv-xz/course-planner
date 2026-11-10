package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"sjqu-course-planner/internal/config"
)

var (
	// Log 结构化 logger（推荐使用）
	Log *zap.Logger
	// Sugared 便捷 logger（类似 log.Printf，也可用）
	Sugared *zap.SugaredLogger
)

// Init 按配置初始化 zap：
//   - 控制台输出 console 格式（带颜色）
//   - 文件输出 json 格式（便于排查）
//   - 级别由配置决定
func Init(cfg config.Logger) error {
	// 1. 解析日志级别
	level := zapcore.InfoLevel
	if cfg.Level != "" {
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			return err
		}
	}

	// 2. 控制台输出（console 编码器）
	consoleEncoder := newConsoleEncoder()
	cores := []zapcore.Core{
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	}

	// 3. 文件输出（json 编码器，若配置了路径）
	if cfg.File != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.File), 0o755); err != nil {
			return err
		}
		f, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return err
		}
		cores = append(cores, zapcore.NewCore(newJSONEncoder(), zapcore.AddSync(f), level))
	}

	// 4. 多路输出
	core := zapcore.NewTee(cores...)
	Log = zap.New(core, zap.AddCaller())
	Sugared = Log.Sugar()
	return nil
}

// Sync 刷新缓冲（进程退出前调用）
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

func baseEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	return cfg
}

func newConsoleEncoder() zapcore.Encoder {
	cfg := baseEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(cfg)
}

func newJSONEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(baseEncoderConfig())
}
