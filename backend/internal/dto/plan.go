package dto

// PlanCreateRequest 新建培养方案
type PlanCreateRequest struct {
	MajorID      int64    `json:"major_id" binding:"required"`
	Grade        string   `json:"grade" binding:"required,len=4,numeric"`
	Name         string   `json:"name" binding:"required,max=128"`
	TotalCredits *float64 `json:"total_credits" binding:"omitempty,gt=0"`
}

// PlanUpdateRequest 修改培养方案
type PlanUpdateRequest struct {
	Grade        string   `json:"grade" binding:"required,len=4,numeric"`
	Name         string   `json:"name" binding:"required,max=128"`
	TotalCredits *float64 `json:"total_credits" binding:"omitempty,gt=0"`
}
