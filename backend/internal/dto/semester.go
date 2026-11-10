package dto

// SemesterCreateRequest 新建学期
type SemesterCreateRequest struct {
	Code         string `json:"code" binding:"required"`
	AcademicYear string `json:"academic_year" binding:"required"`
	TermNo       int8   `json:"term_no" binding:"required,oneof=1 2"`
	SequenceNo   int    `json:"sequence_no" binding:"required,gt=0"`
	DisplayName  string `json:"display_name" binding:"required"`
}

// SemesterUpdateRequest 修改学期
type SemesterUpdateRequest struct {
	Code         string `json:"code" binding:"required"`
	AcademicYear string `json:"academic_year" binding:"required"`
	TermNo       int8   `json:"term_no" binding:"required,oneof=1 2"`
	SequenceNo   int    `json:"sequence_no" binding:"required,gt=0"`
	DisplayName  string `json:"display_name" binding:"required"`
}
