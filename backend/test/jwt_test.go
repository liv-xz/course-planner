package test

import (
	"testing"

	"sjqu-course-planner/internal/utils"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret"
	token, err := utils.GenerateToken(42, "student", secret, 1)
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}
	claims, err := utils.ParseToken(token, secret)
	if err != nil {
		t.Fatalf("解析 token 失败: %v", err)
	}
	if claims.UserID != 42 || claims.Role != "student" {
		t.Fatalf("claims 不一致: %+v", claims)
	}
}

func TestParseToken_Invalid(t *testing.T) {
	if _, err := utils.ParseToken("garbage.token.value", "test-secret"); err == nil {
		t.Fatal("非法 token 应返回错误")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, err := utils.GenerateToken(1, "student", "secret-a", 1)
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}
	if _, err := utils.ParseToken(token, "secret-b"); err == nil {
		t.Fatal("不同密钥应校验失败")
	}
}

func TestParseToken_Expired(t *testing.T) {
	token, err := utils.GenerateToken(1, "student", "test-secret", -1)
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}
	if _, err := utils.ParseToken(token, "test-secret"); err == nil {
		t.Fatal("过期 token 应返回错误")
	}
}
