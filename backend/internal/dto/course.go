package dto

// CourseCreateRequest 新建课程
type CourseCreateRequest struct {
	Code    string  `json:"code" binding:"required,max=32"`
	Name    string  `json:"name" binding:"required,max=128"`
	Credits float64 `json:"credits" binding:"required,gt=0"`
	Nature  int8    `json:"nature" binding:"omitempty,oneof=1 2"`
}

// CourseUpdateRequest 修改课程
type CourseUpdateRequest struct {
	Name    string  `json:"name" binding:"required,max=128"`
	Credits float64 `json:"credits" binding:"required,gt=0"`
	Nature  int8    `json:"nature" binding:"omitempty,oneof=1 2"`
}

// CourseListQuery 课程列表查询参数
type CourseListQuery struct {
	Page       int    `form:"page"`
	Size       int    `form:"size"`
	Keyword    string `form:"keyword"`
	CategoryID int64  `form:"category_id"`
}
