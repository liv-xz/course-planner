package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/router"
)

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return router.Setup(config.C, testDB)
}

func doLogin(t *testing.T, r *gin.Engine, username, password string) *httptest.ResponseRecorder {
	t.Helper()
	body, _ := json.Marshal(dto.LoginRequest{Username: username, Password: password})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestLogin_Success(t *testing.T) {
	r := newTestRouter()
	w := doLogin(t, r, testUser, testPass)
	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d，body=%s", w.Code, w.Body.String())
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
			Role  string `json:"role"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if resp.Code != 0 || resp.Data.Token == "" {
		t.Fatalf("登录失败: %s", w.Body.String())
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	r := newTestRouter()
	w := doLogin(t, r, testUser, "wrong-password")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望 401，实际 %d", w.Code)
	}
}

func TestProfile_WithToken(t *testing.T) {
	r := newTestRouter()

	w := doLogin(t, r, testUser, testPass)
	var loginResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("解析登录响应失败: %v", err)
	}
	token := loginResp.Data.Token
	if token == "" {
		t.Fatal("未拿到 token")
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	pw := httptest.NewRecorder()
	r.ServeHTTP(pw, req)

	if pw.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d，body=%s", pw.Code, pw.Body.String())
	}
	var prof struct {
		Data struct {
			Username string `json:"username"`
			Role     string `json:"role"`
		} `json:"data"`
	}
	if err := json.Unmarshal(pw.Body.Bytes(), &prof); err != nil {
		t.Fatalf("解析 profile 失败: %v", err)
	}
	if prof.Data.Username != testUser {
		t.Fatalf("profile 用户名不符: %s", pw.Body.String())
	}
}

func TestProfile_NoToken(t *testing.T) {
	r := newTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("期望 401，实际 %d", w.Code)
	}
}
