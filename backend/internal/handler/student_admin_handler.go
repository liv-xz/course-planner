package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/dto"
	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

// StudentAdminHandler 管理端学生/修读记录接口
type StudentAdminHandler struct {
	svc            *service.StudentAdminService
	studentSvc     *service.StudentService
	importSvc      *service.ImportService
	maxImportBytes int64
}

func NewStudentAdminHandler(svc *service.StudentAdminService, studentSvc *service.StudentService,
	importSvc *service.ImportService, maxImportBytes int64) *StudentAdminHandler {
	return &StudentAdminHandler{svc: svc, studentSvc: studentSvc, importSvc: importSvc, maxImportBytes: maxImportBytes}
}

func (h *StudentAdminHandler) List(c *gin.Context) {
	var q dto.StudentListQuery
	_ = c.ShouldBindQuery(&q)
	list, total, err := h.svc.ListStudents(c.Request.Context(), q.Page, q.Size, q.Keyword)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.Page(c, list, total)
}

func (h *StudentAdminHandler) Create(c *gin.Context) {
	var req dto.StudentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}
	row, err := h.svc.CreateStudent(c.Request.Context(), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, row)
}

func (h *StudentAdminHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.StudentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}
	row, err := h.svc.UpdateStudent(c.Request.Context(), id, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, row)
}

func (h *StudentAdminHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteStudent(c.Request.Context(), id); err != nil {
		respondServiceError(c, err)
		return
	}
	response.OKMessage(c, "删除成功")
}

func (h *StudentAdminHandler) Records(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	list, err := h.svc.ListRecords(c.Request.Context(), id)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, list)
}

func (h *StudentAdminHandler) CreateRecord(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.RecordCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}
	item, err := h.svc.CreateRecord(c.Request.Context(), id, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, item)
}

func (h *StudentAdminHandler) ImportRecords(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "请选择导入文件")
		return
	}
	if h.maxImportBytes > 0 && fileHeader.Size > h.maxImportBytes {
		response.Error(c, http.StatusBadRequest, 40001, "文件过大")
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "读取文件失败")
		return
	}
	defer file.Close()
	result, err := h.importSvc.ImportRecords(c.Request.Context(), id, fileHeader.Filename, file)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, result)
}

func (h *StudentAdminHandler) UpdateRecord(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.RecordUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}
	item, err := h.svc.UpdateRecord(c.Request.Context(), id, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, item)
}

func (h *StudentAdminHandler) DeleteRecord(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteRecord(c.Request.Context(), id); err != nil {
		respondServiceError(c, err)
		return
	}
	response.OKMessage(c, "删除成功")
}

func (h *StudentAdminHandler) Audit(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	res, err := h.studentSvc.MyAudit(c.Request.Context(), id)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, res)
}

func (h *StudentAdminHandler) Recommendations(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	semester := c.Query("semester")
	list, err := h.studentSvc.MyRecommendations(c.Request.Context(), id, semester)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	response.OK(c, list)
}

func parseQueryInt(c *gin.Context, key string, def int) int {
	v, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(def)))
	if err != nil {
		return def
	}
	return v
}
