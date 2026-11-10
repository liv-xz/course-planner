package dto

// StudentCreateRequest 管理端新建学生
type StudentCreateRequest struct {
	Username  string `json:"username" binding:"required,max=64"`
	Password  string `json:"password" binding:"required,min=6"`
	Name      string `json:"name" binding:"max=32"`
	StudentNo string `json:"student_no" binding:"required,max=32"`
	MajorID   int64  `json:"major_id" binding:"required"`
	Grade     string `json:"grade" binding:"required,len=4,numeric"`
	Status    int8   `json:"status"`
}

// StudentUpdateRequest 管理端修改学生
type StudentUpdateRequest struct {
	Name      string `json:"name" binding:"max=32"`
	StudentNo string `json:"student_no" binding:"required,max=32"`
	MajorID   int64  `json:"major_id" binding:"required"`
	Grade     string `json:"grade" binding:"required,len=4,numeric"`
	Status    int8   `json:"status"`
}

// StudentListQuery 学生列表查询
type StudentListQuery struct {
	Page    int    `form:"page"`
	Size    int    `form:"size"`
	Keyword string `form:"keyword"`
}

// RecordCreateRequest 管理端新增修读记录
type RecordCreateRequest struct {
	CourseID int64    `json:"course_id" binding:"required"`
	Score    *float64 `json:"score"`
	Semester *string  `json:"semester"`
	Status   int8     `json:"status"`
	Credits  *float64 `json:"credits"`
}

// RecordUpdateRequest 管理端修改修读记录
type RecordUpdateRequest struct {
	CourseID int64    `json:"course_id" binding:"required"`
	Score    *float64 `json:"score"`
	Semester *string  `json:"semester"`
	Status   int8     `json:"status"`
	Credits  *float64 `json:"credits"`
}
