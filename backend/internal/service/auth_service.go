package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
	"sjqu-course-planner/internal/utils"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserDisabled       = errors.New("账号已被禁用")
)

// LoginResult 登录结果
type LoginResult struct {
	User  *model.User
	Role  string
	Token string
}

// AuthService 认证业务
type AuthService struct {
	userRepo    *repository.UserRepo
	roleRepo    *repository.RoleRepo
	profileRepo *repository.StudentProfileRepo
	jwtCfg      config.JWT
}

func NewAuthService(userRepo *repository.UserRepo, roleRepo *repository.RoleRepo,
	profileRepo *repository.StudentProfileRepo, jwtCfg config.JWT) *AuthService {
	return &AuthService{userRepo: userRepo, roleRepo: roleRepo, profileRepo: profileRepo, jwtCfg: jwtCfg}
}

// Login 校验用户名密码并签发 token
func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		// 不区分"用户不存在/密码错误"，避免暴露账号是否存在
		return nil, ErrInvalidCredentials
	}
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, role.Name, s.jwtCfg.Secret, s.jwtCfg.ExpireHours)
	if err != nil {
		return nil, err
	}
	return &LoginResult{User: user, Role: role.Name, Token: token}, nil
}

// GetProfile 获取当前用户 + 学生档案（管理员无档案时 profile 为 nil）
func (s *AuthService) GetProfile(ctx context.Context, userID int64) (*model.User, *model.StudentProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return user, nil, nil
	}
	return user, profile, nil
}
