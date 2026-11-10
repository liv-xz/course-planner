package dto

// PrerequisiteCreateRequest 新增先修关系
type PrerequisiteCreateRequest struct {
	PrereqCourseID int64 `json:"prereq_course_id" binding:"required,gt=0"`
}

// OfferingCreateRequest 新增开课信息
type OfferingCreateRequest struct {
	Semester string `json:"semester" binding:"required,max=16"`
}
