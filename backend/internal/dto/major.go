package dto

// MajorCreateRequest 新建专业
type MajorCreateRequest struct {
	Code string `json:"code" binding:"required,max=32"`
	Name string `json:"name" binding:"required,max=64"`
}

// MajorUpdateRequest 修改专业
type MajorUpdateRequest struct {
	Code string `json:"code" binding:"required,max=32"`
	Name string `json:"name" binding:"required,max=64"`
}
