package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sjqu-course-planner/internal/response"
	"sjqu-course-planner/internal/service"
)

func parseID(c *gin.Context) (int64, bool) {
	return parseIDNamed(c, "id")
}

func parseIDNamed(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, 40001, "无效的 ID")
		return 0, false
	}
	return id, true
}

func queryInt(c *gin.Context, key string, def int) int {
	v, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(def)))
	if err != nil {
		return def
	}
	return v
}

func respondServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	if isNotFound(err) {
		response.Error(c, http.StatusNotFound, 40401, err.Error())
		return
	}
	if isConflict(err) {
		response.Error(c, http.StatusConflict, 40901, err.Error())
		return
	}
	if isBusiness(err) {
		response.Error(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	response.Error(c, http.StatusInternalServerError, 50001, "服务器内部错误")
}

func isNotFound(err error) bool {
	return errors.Is(err, service.ErrPlanNotFound) ||
		errors.Is(err, service.ErrCategoryNotFound) ||
		errors.Is(err, service.ErrCourseNotFound) ||
		errors.Is(err, service.ErrPrerequisiteNotFound) ||
		errors.Is(err, service.ErrOfferingNotFound) ||
		errors.Is(err, service.ErrMajorNotFound) ||
		errors.Is(err, service.ErrSemesterNotFound) ||
		errors.Is(err, service.ErrStudentProfileNotFound) ||
		errors.Is(err, service.ErrStudentNotFound)
}

func isConflict(err error) bool {
	return errors.Is(err, service.ErrPlanInUse) ||
		errors.Is(err, service.ErrCategoryInUse) ||
		errors.Is(err, service.ErrMajorInUse) ||
		errors.Is(err, service.ErrSemesterInUse) ||
		errors.Is(err, service.ErrSemesterExists) ||
		errors.Is(err, service.ErrMajorExists) ||
		errors.Is(err, service.ErrCourseCodeExists) ||
		errors.Is(err, service.ErrOfferingExists) ||
		errors.Is(err, service.ErrPrerequisiteExists) ||
		errors.Is(err, service.ErrCourseCategoryDup) ||
		errors.Is(err, service.ErrCourseDeleted) ||
		errors.Is(err, service.ErrUsernameExists) ||
		errors.Is(err, service.ErrStudentNoExists)
}

func isBusiness(err error) bool {
	return errors.Is(err, service.ErrCategoryCycle) ||
		errors.Is(err, service.ErrCategoryHasCourses) ||
		errors.Is(err, service.ErrCategoryNotLeaf) ||
		errors.Is(err, service.ErrCategoryOtherPlan) ||
		errors.Is(err, service.ErrParentCategoryBad) ||
		errors.Is(err, service.ErrCourseInvalid) ||
		errors.Is(err, service.ErrSemesterInvalid) ||
		errors.Is(err, service.ErrPrerequisiteSelf) ||
		errors.Is(err, service.ErrPrerequisiteCycle) ||
		errors.Is(err, service.ErrStudentProfileNotFound) ||
		errors.Is(err, service.ErrInvalidRecord) ||
		errors.Is(err, service.ErrRecordNotFound) ||
		errors.Is(err, service.ErrInvalidRole)
}
