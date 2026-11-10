package dto

// CategoryCreateRequest 新建课程类别
type CategoryCreateRequest struct {
	PlanID          int64   `json:"plan_id" binding:"required"`
	ParentID        *int64  `json:"parent_id"`
	Name            string  `json:"name" binding:"required,max=64"`
	RequiredCredits float64 `json:"required_credits" binding:"gte=0"`
	Sort            int     `json:"sort" binding:"gte=0"`
}

// CategoryUpdateRequest 修改课程类别
type CategoryUpdateRequest struct {
	ParentID        *int64  `json:"parent_id"`
	Name            string  `json:"name" binding:"required,max=64"`
	RequiredCredits float64 `json:"required_credits" binding:"gte=0"`
	Sort            int     `json:"sort" binding:"gte=0"`
}

// CategoryNode 类别树节点
type CategoryNode struct {
	ID              int64           `json:"id"`
	PlanID          int64           `json:"plan_id"`
	ParentID        *int64          `json:"parent_id"`
	Name            string          `json:"name"`
	RequiredCredits float64         `json:"required_credits"`
	Sort            int             `json:"sort"`
	Children        []*CategoryNode `json:"children"`
}
