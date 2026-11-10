package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// App 应用基本配置
type App struct {
	Name            string `mapstructure:"name"`
	Env             string `mapstructure:"env"`
	Port            int    `mapstructure:"port"`
	CurrentSemester string `mapstructure:"current_semester"`
}

// Logger 日志配置
type Logger struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	File   string `mapstructure:"file"`
}

// Database 数据库配置
type Database struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// JWT 认证配置
type JWT struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// LLM 大模型配置（OpenAI 兼容接口）
type LLM struct {
	Enabled        bool    `mapstructure:"enabled"`
	Provider       string  `mapstructure:"provider"`
	BaseURL        string  `mapstructure:"base_url"`
	APIKey         string  `mapstructure:"api_key"`
	Model          string  `mapstructure:"model"`
	TimeoutSeconds int     `mapstructure:"timeout_seconds"`
	Temperature    float64 `mapstructure:"temperature"`
	MaxTokens      int     `mapstructure:"max_tokens"`
}

// RAG 检索增强配置
type RAG struct {
	Enabled  bool    `mapstructure:"enabled"`
	TopK     int     `mapstructure:"top_k"`
	MinScore float64 `mapstructure:"min_score"`
}

// Security 输入安全配置
type Security struct {
	Enabled        bool     `mapstructure:"enabled"`
	MaxInputLength int      `mapstructure:"max_input_length"`
	Patterns       []string `mapstructure:"patterns"`
}

// Import 文件导入配置
type Import struct {
	MaxFileSizeMB   int      `mapstructure:"max_file_size_mb"`
	MaxRows         int      `mapstructure:"max_rows"`
	AllowedSuffixes []string `mapstructure:"allowed_suffixes"`
}

// Agent 智能体配置
type Agent struct {
	Enabled       bool   `mapstructure:"enabled"`
	MaxToolRounds int    `mapstructure:"max_tool_rounds"`
	SystemPrompt  string `mapstructure:"system_prompt"`
}

// Config 全局配置
type Config struct {
	App      App      `mapstructure:"app"`
	Logger   Logger   `mapstructure:"logger"`
	Database Database `mapstructure:"database"`
	JWT      JWT      `mapstructure:"jwt"`
	LLM      LLM      `mapstructure:"llm"`
	RAG      RAG      `mapstructure:"rag"`
	Security Security `mapstructure:"security"`
	Import   Import   `mapstructure:"import"`
	Agent    Agent    `mapstructure:"agent"`
}

// C 全局配置实例
var C Config

// Load 加载配置文件（支持 CONFIG_PATH 环境变量指定路径）
func Load() error {
	v := viper.New()

	if path := os.Getenv("CONFIG_PATH"); path != "" {
		v.SetConfigFile(path)
	} else if _, err := os.Stat("configs/config.local.yaml"); err == nil {
		v.SetConfigFile("configs/config.local.yaml")
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("configs")
		v.AddConfigPath(".")
	}

	// 环境变量覆盖：SJQU_APP_PORT -> app.port
	v.SetEnvPrefix("SJQU")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置失败: %w", err)
	}
	if err := v.Unmarshal(&C); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}
	return nil
}
