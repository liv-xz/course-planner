package service

import (
	"context"
	"strings"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

// SecurityService 输入安全检查与日志
type SecurityService struct {
	repo *repository.SecurityLogRepo
	cfg  config.Security
}

func NewSecurityService(repo *repository.SecurityLogRepo, cfg config.Security) *SecurityService {
	return &SecurityService{repo: repo, cfg: cfg}
}

func (s *SecurityService) CheckInput(ctx context.Context, userID int64, input string) (bool, string) {
	if !s.cfg.Enabled {
		return true, ""
	}
	if s.cfg.MaxInputLength > 0 && len([]rune(input)) > s.cfg.MaxInputLength {
		s.log(ctx, userID, input, "输入过长")
		return false, "输入内容过长，请精简后重试"
	}
	lower := strings.ToLower(input)
	for _, p := range s.cfg.Patterns {
		if p == "" {
			continue
		}
		if strings.Contains(lower, strings.ToLower(p)) {
			s.log(ctx, userID, input, "命中安全规则: "+p)
			return false, "输入包含不安全的指令，已被拦截"
		}
	}
	return true, ""
}

func (s *SecurityService) log(ctx context.Context, userID int64, input, reason string) {
	in := input
	r := reason
	uid := userID
	_ = s.repo.Create(ctx, &model.SecurityLog{UserID: &uid, Input: &in, Reason: &r})
}

func (s *SecurityService) List(ctx context.Context, page, size int) ([]model.SecurityLog, int64, error) {
	return s.repo.List(ctx, page, size)
}
