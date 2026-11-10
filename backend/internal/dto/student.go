package dto

// RecordItem 修读记录项
type RecordItem struct {
	ID         int64    `json:"id"`
	CourseID   int64    `json:"course_id"`
	CourseCode string   `json:"course_code"`
	CourseName string   `json:"course_name"`
	Credits    float64  `json:"credits"`
	Score      *float64 `json:"score"`
	Semester   *string  `json:"semester"`
	Status     int8     `json:"status"`
}

// AuditCategory 学分审核类别节点
type AuditCategory struct {
	ID               int64            `json:"id"`
	Name             string           `json:"name"`
	RequiredCredits  float64          `json:"required_credits"`
	CompletedCredits float64          `json:"completed_credits"`
	GapCredits       float64          `json:"gap_credits"`
	Children         []*AuditCategory `json:"children"`
}

// AuditResponse 学分审核响应
type AuditResponse struct {
	PlanID     int64            `json:"plan_id"`
	PlanName   string           `json:"plan_name"`
	Grade      string           `json:"grade"`
	Categories []*AuditCategory `json:"categories"`
}

// RecommendationItem 选课推荐项
type RecommendationItem struct {
	CourseID     int64    `json:"course_id"`
	CourseCode   string   `json:"course_code"`
	CourseName   string   `json:"course_name"`
	Credits      float64  `json:"credits"`
	CategoryID   int64    `json:"category_id"`
	CategoryName string   `json:"category_name"`
	Semesters    []string `json:"semesters"`
	Reason       string   `json:"reason"`
}
