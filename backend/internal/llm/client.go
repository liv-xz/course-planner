package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"sjqu-course-planner/internal/config"
)

var ErrNotConfigured = errors.New("LLM 未配置")

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

type Tool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string         `json:"name"`
		Description string         `json:"description"`
		Parameters  map[string]any `json:"parameters"`
	} `json:"function"`
}

type Client struct {
	cfg  config.LLM
	http *http.Client
}

func NewClient(cfg config.LLM) *Client {
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &Client{cfg: cfg, http: &http.Client{Timeout: timeout}}
}

func (c *Client) Ready() bool {
	return c.cfg.Enabled && strings.TrimSpace(c.cfg.APIKey) != "" && strings.TrimSpace(c.cfg.Model) != ""
}

func (c *Client) Chat(ctx context.Context, messages []Message, tools []Tool) (*Message, error) {
	if !c.Ready() {
		return nil, ErrNotConfigured
	}
	endpoint := strings.TrimRight(c.cfg.BaseURL, "/")
	if !strings.HasSuffix(endpoint, "/chat/completions") {
		endpoint += "/chat/completions"
	}
	body := map[string]any{
		"model":       c.cfg.Model,
		"messages":    messages,
		"temperature": c.cfg.Temperature,
	}
	if c.cfg.MaxTokens > 0 {
		body["max_tokens"] = c.cfg.MaxTokens
	}
	if len(tools) > 0 {
		body["tools"] = tools
		body["tool_choice"] = "auto"
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("LLM 请求失败: HTTP %d %s", resp.StatusCode, truncate(string(data), 500))
	}
	var parsed struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("解析 LLM 响应失败: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("LLM 没有返回结果")
	}
	return &parsed.Choices[0].Message, nil
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
