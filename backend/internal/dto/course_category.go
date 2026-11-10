package dto

// CourseCategoryItem 课程归属信息
type CourseCategoryItem struct {
	PlanID          int64   `json:"plan_id"`
	PlanName        string  `json:"plan_name"`
	CategoryID      int64   `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	RequiredCredits float64 `json:"required_credits"`
}

// SetCourseCategoryRequest 设置课程在某个培养方案下的类别；category_id 为空表示解除
type SetCourseCategoryRequest struct {
	CategoryID *int64 `json:"category_id"`
}
