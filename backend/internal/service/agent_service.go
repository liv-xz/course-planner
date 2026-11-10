package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/llm"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

type EmitFunc func(event string, data any)

// AgentService 选课规划智能体
type AgentService struct {
	studentSvc   *StudentService
	courseRepo   *repository.CourseRepo
	prereqRepo   *repository.PrerequisiteRepo
	offeringRepo *repository.CourseOfferingRepo
	recordRepo   *repository.StudentRecordRepo
	rag          *RAGService
	security     *SecurityService
	sessionRepo  *repository.ChatSessionRepo
	msgRepo      *repository.ChatMessageRepo
	llm          *llm.Client
	cfg          config.Agent
}

func NewAgentService(studentSvc *StudentService, courseRepo *repository.CourseRepo,
	prereqRepo *repository.PrerequisiteRepo, offeringRepo *repository.CourseOfferingRepo,
	recordRepo *repository.StudentRecordRepo, rag *RAGService, security *SecurityService,
	sessionRepo *repository.ChatSessionRepo, msgRepo *repository.ChatMessageRepo,
	llmClient *llm.Client, cfg config.Agent) *AgentService {
	return &AgentService{
		studentSvc: studentSvc, courseRepo: courseRepo, prereqRepo: prereqRepo,
		offeringRepo: offeringRepo, recordRepo: recordRepo, rag: rag, security: security,
		sessionRepo: sessionRepo, msgRepo: msgRepo, llm: llmClient, cfg: cfg,
	}
}

func (s *AgentService) ChatStream(ctx context.Context, userID, sessionID int64, input string, emit EmitFunc) error {
	if strings.TrimSpace(input) == "" {
		emit("error", map[string]any{"message": "请输入问题"})
		return nil
	}
	ok, reason := s.security.CheckInput(ctx, userID, input)
	if !ok {
		emit("error", map[string]any{"message": reason})
		emit("done", map[string]any{"session_id": sessionID})
		return nil
	}
	session, err := s.ensureSession(ctx, userID, sessionID, input)
	if err != nil {
		return err
	}
	sessionID = session.ID
	emit("meta", map[string]any{"session_id": sessionID, "model": s.llmModel()})
	s.saveMessage(ctx, sessionID, "user", input, nil)

	answer, err := s.run(ctx, userID, sessionID, input, emit)
	if err != nil {
		answer = s.fallbackAnswer(ctx, userID, input, emit, err)
	}
	if answer == "" {
		answer = "暂时没有生成有效回答，请换个问法。"
	}
	s.saveMessage(ctx, sessionID, "assistant", answer, nil)
	emit("done", map[string]any{"session_id": sessionID})
	return nil
}

func (s *AgentService) run(ctx context.Context, userID, sessionID int64, input string, emit EmitFunc) (string, error) {
	if !s.llm.Ready() {
		return "", llm.ErrNotConfigured
	}
	hist, err := s.msgRepo.ListBySession(ctx, sessionID)
	if err != nil {
		return "", err
	}
	ragHits, _ := s.rag.Search(ctx, input)
	system := s.cfg.SystemPrompt
	if system == "" {
		system = "你是上海建桥学院选课规划助手。学分、先修、开课学期以工具返回为准，不得编造。"
	}
	if len(ragHits) > 0 {
		parts := make([]string, 0, len(ragHits))
		for _, h := range ragHits {
			parts = append(parts, "【"+h.Title+"】"+h.Content)
		}
		system += "\n\n参考知识：\n" + strings.Join(parts, "\n")
	}
	messages := []llm.Message{{Role: "system", Content: system}}
	for _, m := range hist {
		if m.Role != "user" && m.Role != "assistant" {
			continue
		}
		if m.Content != nil {
			messages = append(messages, llm.Message{Role: m.Role, Content: *m.Content})
		}
	}
	// hist 已包含当前用户消息，避免重复
	rounds := s.cfg.MaxToolRounds
	if rounds <= 0 {
		rounds = 3
	}
	tools := agentTools()
	for i := 0; i < rounds; i++ {
		msg, err := s.llm.Chat(ctx, messages, tools)
		if err != nil {
			return "", err
		}
		if len(msg.ToolCalls) == 0 {
			answer := strings.TrimSpace(msg.Content)
			if answer == "" {
				return "", fmt.Errorf("模型没有返回文本")
			}
			emitText(emit, answer)
			return answer, nil
		}
		messages = append(messages, llm.Message{Role: "assistant", Content: msg.Content, ToolCalls: msg.ToolCalls})
		for _, tc := range msg.ToolCalls {
			emit("tool_call", map[string]any{"name": tc.Function.Name, "arguments": tc.Function.Arguments})
			result, err := s.executeTool(ctx, userID, tc.Function.Name, tc.Function.Arguments)
			if err != nil {
				result = map[string]any{"error": err.Error()}
			}
			emit("tool_result", map[string]any{"name": tc.Function.Name, "result": result})
			b, _ := json.Marshal(result)
			messages = append(messages, llm.Message{Role: "tool", ToolCallID: tc.ID, Name: tc.Function.Name, Content: string(b)})
		}
	}
	return "", fmt.Errorf("工具调用轮次超限")
}

func (s *AgentService) fallbackAnswer(ctx context.Context, userID int64, input string, emit EmitFunc, cause error) string {
	_ = cause
	var b strings.Builder
	ragHits, _ := s.rag.Search(ctx, input)
	needAudit := strings.Contains(input, "学分") || strings.Contains(input, "缺口") || strings.Contains(input, "差")
	needRec := strings.Contains(input, "推荐") || strings.Contains(input, "选课") || strings.Contains(input, "下学期") || strings.Contains(input, "课程")
	if needAudit || (!needRec && len(ragHits) == 0) {
		emit("tool_call", map[string]any{"name": "audit_credit_gap", "arguments": "{}"})
		audit, err := s.studentSvc.MyAudit(ctx, userID)
		if err == nil {
			emit("tool_result", map[string]any{"name": "audit_credit_gap", "result": audit})
			b.WriteString("按你的培养方案，当前学分情况：\n")
			for _, c := range audit.Categories {
				b.WriteString(fmt.Sprintf("- %s：要求 %.1f，已完成 %.1f，缺口 %.1f\n", c.Name, c.RequiredCredits, c.CompletedCredits, c.GapCredits))
			}
		} else {
			b.WriteString("暂时无法读取你的学分审核结果：" + err.Error() + "\n")
		}
	}
	if needRec {
		emit("tool_call", map[string]any{"name": "list_candidate_courses", "arguments": "{}"})
		list, err := s.studentSvc.MyRecommendations(ctx, userID, "")
		if err == nil {
			emit("tool_result", map[string]any{"name": "list_candidate_courses", "result": list})
			if len(list) == 0 {
				b.WriteString("当前没有满足硬约束的候选课程。\n")
			} else {
				b.WriteString("\n候选课程（按学分和缺口筛选）：\n")
				limit := len(list)
				if limit > 8 {
					limit = 8
				}
				for _, c := range list[:limit] {
					b.WriteString(fmt.Sprintf("- %s %s（%.1f 学分，%s，开课：%s）\n", c.CourseCode, c.CourseName, c.Credits, c.CategoryName, strings.Join(c.Semesters, "、")))
				}
			}
		} else {
			b.WriteString("暂时无法生成推荐：" + err.Error() + "\n")
		}
	}
	if len(ragHits) > 0 {
		b.WriteString("\n相关规则：\n")
		for _, h := range ragHits {
			b.WriteString("- " + h.Content + "\n")
		}
	}
	if b.Len() == 0 {
		b.WriteString("我可以帮你查学分缺口、筛选课程、检查先修关系。你可以问：“我还差哪些学分？”或“下学期推荐什么课？”")
	}
	answer := strings.TrimSpace(b.String())
	emitText(emit, answer)
	return answer
}

func (s *AgentService) executeTool(ctx context.Context, userID int64, name, args string) (any, error) {
	switch name {
	case "audit_credit_gap", "get_my_credits":
		return s.studentSvc.MyAudit(ctx, userID)
	case "list_candidate_courses":
		var in struct {
			Semester   string  `json:"semester"`
			MaxCredits float64 `json:"max_credits"`
		}
		_ = json.Unmarshal([]byte(args), &in)
		list, err := s.studentSvc.MyRecommendations(ctx, userID, in.Semester)
		if err != nil {
			return nil, err
		}
		if in.MaxCredits > 0 {
			out := make([]any, 0)
			total := 0.0
			for _, c := range list {
				if total+c.Credits > in.MaxCredits {
					continue
				}
				total += c.Credits
				out = append(out, c)
			}
			return out, nil
		}
		return list, nil
	case "check_prerequisite":
		var in struct {
			CourseID   int64  `json:"course_id"`
			CourseName string `json:"course_name"`
		}
		_ = json.Unmarshal([]byte(args), &in)
		return s.checkCourse(ctx, userID, in.CourseID, in.CourseName)
	case "search_rules":
		var in struct {
			Query string `json:"query"`
		}
		_ = json.Unmarshal([]byte(args), &in)
		if in.Query == "" {
			in.Query = args
		}
		return s.rag.Search(ctx, in.Query)
	default:
		return nil, fmt.Errorf("未知工具: %s", name)
	}
}

func (s *AgentService) checkCourse(ctx context.Context, userID, courseID int64, courseName string) (any, error) {
	var course *model.Course
	var err error
	if courseID > 0 {
		course, err = s.courseRepo.GetByID(ctx, courseID)
	} else {
		if strings.TrimSpace(courseName) == "" {
			return nil, fmt.Errorf("请提供课程 ID 或课程名称")
		}
		list, _, e := s.courseRepo.Search(ctx, 1, 100, courseName, 0)
		if e != nil {
			return nil, e
		}
		if len(list) == 0 {
			return nil, fmt.Errorf("没有找到课程: %s", courseName)
		}
		course = &list[0]
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("课程不存在")
		}
		return nil, err
	}
	records, err := s.recordRepo.ListByStudent(ctx, userID)
	if err != nil {
		return nil, err
	}
	passed := map[int64]bool{}
	for _, r := range records {
		if r.Status == 1 {
			passed[r.CourseID] = true
		}
	}
	prereqs, err := s.prereqRepo.ListByCourse(ctx, course.ID)
	if err != nil {
		return nil, err
	}
	missing := []string{}
	for _, p := range prereqs {
		if !passed[p.PrereqCourseID] {
			if c, err := s.courseRepo.GetByIDIncludingDeleted(ctx, p.PrereqCourseID); err == nil {
				missing = append(missing, c.Name)
			}
		}
	}
	offerings, _ := s.offeringRepo.ListByCourse(ctx, course.ID)
	semesters := []string{}
	for _, o := range offerings {
		semesters = append(semesters, o.Semester)
	}
	return map[string]any{
		"course_id": course.ID, "course_code": course.Code, "course_name": course.Name,
		"credits": course.Credits, "passed": passed[course.ID],
		"prerequisites_met": len(missing) == 0, "missing_prerequisites": missing,
		"semesters": semesters,
	}, nil
}

func (s *AgentService) ensureSession(ctx context.Context, userID, sessionID int64, input string) (*model.ChatSession, error) {
	if sessionID > 0 {
		session, err := s.sessionRepo.GetByID(ctx, sessionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("会话不存在")
			}
			return nil, err
		}
		if session.UserID != userID {
			return nil, fmt.Errorf("无权访问该会话")
		}
		return session, nil
	}
	title := input
	if r := []rune(title); len(r) > 30 {
		title = string(r[:30])
	}
	session := &model.ChatSession{UserID: userID, Title: &title}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *AgentService) saveMessage(ctx context.Context, sessionID int64, role, content string, toolCalls []llm.ToolCall) {
	var tc *string
	if len(toolCalls) > 0 {
		if b, err := json.Marshal(toolCalls); err == nil {
			v := string(b)
			tc = &v
		}
	}
	_ = s.msgRepo.Create(ctx, &model.ChatMessage{SessionID: sessionID, Role: role, Content: &content, ToolCalls: tc})
}

func (s *AgentService) ListSessions(ctx context.Context, userID int64) ([]model.ChatSession, error) {
	return s.sessionRepo.ListByUser(ctx, userID)
}

func (s *AgentService) ListMessages(ctx context.Context, userID, sessionID int64) ([]model.ChatMessage, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("会话不存在")
		}
		return nil, err
	}
	if session.UserID != userID {
		return nil, fmt.Errorf("无权访问该会话")
	}
	return s.msgRepo.ListBySession(ctx, sessionID)
}

func (s *AgentService) DeleteSession(ctx context.Context, userID, sessionID int64) error {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("会话不存在")
		}
		return err
	}
	if session.UserID != userID {
		return fmt.Errorf("无权访问该会话")
	}
	if err := s.msgRepo.DeleteBySession(ctx, sessionID); err != nil {
		return err
	}
	return s.sessionRepo.Delete(ctx, sessionID)
}

func (s *AgentService) Mode() string {
	if s.llm.Ready() {
		return "llm"
	}
	return "rule-fallback"
}

func (s *AgentService) llmModel() string {
	if s.llm.Ready() {
		return "llm"
	}
	return "rule-fallback"
}

func emitText(emit EmitFunc, text string) {
	runes := []rune(text)
	for i := 0; i < len(runes); i += 20 {
		end := i + 20
		if end > len(runes) {
			end = len(runes)
		}
		emit("delta", map[string]any{"text": string(runes[i:end])})
		time.Sleep(8 * time.Millisecond)
	}
}

func agentTools() []llm.Tool {
	mk := func(name, desc string, params map[string]any) llm.Tool {
		t := llm.Tool{Type: "function"}
		t.Function.Name = name
		t.Function.Description = desc
		t.Function.Parameters = params
		return t
	}
	obj := func(props map[string]any) map[string]any {
		return map[string]any{"type": "object", "properties": props, "additionalProperties": false}
	}
	return []llm.Tool{
		mk("audit_credit_gap", "查询当前登录学生的学分审核和缺口，不需要任何参数", obj(map[string]any{})),
		mk("list_candidate_courses", "按学分缺口、先修和开课学期筛选候选课程", obj(map[string]any{
			"semester":    map[string]any{"type": "string", "description": "日历学期，如 2024-2025-1，可空"},
			"max_credits": map[string]any{"type": "number", "description": "总学分上限，可空"},
		})),
		mk("check_prerequisite", "检查当前学生能否选修某门课", obj(map[string]any{
			"course_id":   map[string]any{"type": "integer", "description": "课程 ID，可空"},
			"course_name": map[string]any{"type": "string", "description": "课程名称，可空"},
		})),
		mk("search_rules", "检索培养方案和选课规则", obj(map[string]any{
			"query": map[string]any{"type": "string", "description": "检索问题"},
		})),
	}
}
