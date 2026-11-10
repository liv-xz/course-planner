package dto

// OverviewStats 管理端总览
type OverviewStats struct {
	Students   int64 `json:"students"`
	Plans      int64 `json:"plans"`
	Courses    int64 `json:"courses"`
	Majors     int64 `json:"majors"`
	Semesters  int64 `json:"semesters"`
	Categories int64 `json:"categories"`
}

// CategoryGapStat 某个类别的缺口分布
type CategoryGapStat struct {
	PlanID          int64   `json:"plan_id"`
	PlanName        string  `json:"plan_name"`
	CategoryID      int64   `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	RequiredCredits float64 `json:"required_credits"`
	StudentCount    int     `json:"student_count"`
	GapStudentCount int     `json:"gap_student_count"`
	TotalGap        float64 `json:"total_gap"`
	AverageGap      float64 `json:"average_gap"`
}
