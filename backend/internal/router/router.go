package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/handler"
	"sjqu-course-planner/internal/llm"
	"sjqu-course-planner/internal/middleware"
	"sjqu-course-planner/internal/repository"
	"sjqu-course-planner/internal/service"
)

// Setup 创建并配置 Gin 引擎
func Setup(cfg config.Config, db *gorm.DB) *gin.Engine {
	if cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestLogger())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		dbStatus := "unconfigured"
		if db != nil {
			if sqlDB, err := db.DB(); err == nil && sqlDB.Ping() == nil {
				dbStatus = "ok"
			} else {
				dbStatus = "error"
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "db": dbStatus})
	})

	// 数据库未配置时，不注册需要数据库的业务路由
	if db == nil {
		return r
	}

	// ===== Repositories =====
	userRepo := repository.NewUserRepo(db)
	roleRepo := repository.NewRoleRepo(db)
	profileRepo := repository.NewStudentProfileRepo(db)
	majorRepo := repository.NewMajorRepo(db)
	planRepo := repository.NewTrainingPlanRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	courseRepo := repository.NewCourseRepo(db)
	ccRepo := repository.NewCourseCategoryRepo(db)
	prereqRepo := repository.NewPrerequisiteRepo(db)
	offeringRepo := repository.NewCourseOfferingRepo(db)
	auditRepo := repository.NewCreditAuditRepo(db)
	knowledgeRepo := repository.NewKnowledgeDocRepo(db)
	recordRepo := repository.NewStudentRecordRepo(db)
	semesterRepo := repository.NewSemesterRepo(db)
	knowledgeChunkRepo := repository.NewKnowledgeChunkRepo(db)
	securityRepo := repository.NewSecurityLogRepo(db)
	chatSessionRepo := repository.NewChatSessionRepo(db)
	chatMessageRepo := repository.NewChatMessageRepo(db)

	// ===== Services =====
	authSvc := service.NewAuthService(userRepo, roleRepo, profileRepo, cfg.JWT)
	planSvc := service.NewPlanService(planRepo, categoryRepo, auditRepo, knowledgeRepo, majorRepo)
	categorySvc := service.NewCategoryService(categoryRepo, planRepo, ccRepo, auditRepo)
	courseSvc := service.NewCourseService(courseRepo, ccRepo)
	prereqSvc := service.NewPrerequisiteService(prereqRepo, courseRepo)
	offeringSvc := service.NewOfferingService(offeringRepo, courseRepo, semesterRepo)
	studentSvc := service.NewStudentService(profileRepo, planRepo, categoryRepo, courseRepo, ccRepo, prereqRepo, offeringRepo, recordRepo, semesterRepo)
	majorSvc := service.NewMajorService(majorRepo, planRepo, profileRepo)
	semesterSvc := service.NewSemesterService(semesterRepo)
	ccSvc := service.NewCourseCategoryService(ccRepo, courseRepo, categoryRepo, planRepo)
	studentAdminSvc := service.NewStudentAdminService(userRepo, roleRepo, profileRepo, majorRepo, recordRepo, courseRepo, semesterRepo)
	importSvc := service.NewImportService(recordRepo, courseRepo, semesterRepo, profileRepo, cfg.Import)
	securitySvc := service.NewSecurityService(securityRepo, cfg.Security)
	ragSvc := service.NewRAGService(knowledgeChunkRepo, cfg.RAG)
	llmClient := llm.NewClient(cfg.LLM)
	agentSvc := service.NewAgentService(studentSvc, courseRepo, prereqRepo, offeringRepo, recordRepo, ragSvc, securitySvc, chatSessionRepo, chatMessageRepo, llmClient, cfg.Agent)
	statsSvc := service.NewStatsService(profileRepo, planRepo, courseRepo, majorRepo, semesterRepo, categoryRepo, studentSvc)

	// ===== Handlers =====
	authHandler := handler.NewAuthHandler(authSvc)
	planHandler := handler.NewPlanHandler(planSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	courseHandler := handler.NewCourseHandler(courseSvc)
	prereqHandler := handler.NewPrerequisiteHandler(prereqSvc)
	offeringHandler := handler.NewOfferingHandler(offeringSvc)
	studentHandler := handler.NewStudentHandler(studentSvc)
	majorHandler := handler.NewMajorHandler(majorSvc)
	semesterHandler := handler.NewSemesterHandler(semesterSvc)
	ccHandler := handler.NewCourseCategoryHandler(ccSvc)
	studentAdminHandler := handler.NewStudentAdminHandler(studentAdminSvc, studentSvc, importSvc, int64(cfg.Import.MaxFileSizeMB)*1024*1024)
	statsHandler := handler.NewStatsHandler(statsSvc)
	securityHandler := handler.NewSecurityHandler(securitySvc)
	chatHandler := handler.NewChatHandler(agentSvc)

	// ===== 路由 =====
	api := r.Group("/api/v1")

	// 公开
	api.POST("/auth/login", authHandler.Login)

	// 需要登录（读接口：所有登录用户）
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(cfg.JWT.Secret))
	protected.GET("/auth/profile", authHandler.Profile)

	protected.GET("/plans", planHandler.List)
	protected.GET("/plans/:id", planHandler.Get)
	protected.GET("/plans/:id/categories", categoryHandler.ListTree)
	protected.GET("/majors", majorHandler.List)
	protected.GET("/semesters", semesterHandler.List)

	protected.GET("/courses", courseHandler.List)
	protected.GET("/courses/:id", courseHandler.Get)
	protected.GET("/courses/:id/prerequisites", prereqHandler.List)
	protected.GET("/courses/:id/offerings", offeringHandler.List)
	protected.GET("/courses/:id/categories", ccHandler.ListByCourse)

	// 学生端（仅 student 角色，按 JWT 用户隔离）
	student := protected.Group("")
	student.Use(middleware.StudentOnly())
	student.GET("/records", studentHandler.Records)
	student.GET("/audit/me", studentHandler.Audit)
	student.GET("/recommendations", studentHandler.Recommendations)
	student.POST("/chat/completions", chatHandler.Chat)
	student.GET("/chat/sessions", chatHandler.Sessions)
	student.GET("/chat/sessions/:id/messages", chatHandler.Messages)
	student.DELETE("/chat/sessions/:id", chatHandler.DeleteSession)
	student.GET("/chat/status", chatHandler.Status)

	// 仅管理员（写接口）
	admin := protected.Group("")
	admin.Use(middleware.AdminOnly())

	admin.POST("/plans", planHandler.Create)
	admin.PUT("/plans/:id", planHandler.Update)
	admin.DELETE("/plans/:id", planHandler.Delete)

	admin.POST("/majors", majorHandler.Create)
	admin.PUT("/majors/:id", majorHandler.Update)
	admin.DELETE("/majors/:id", majorHandler.Delete)

	admin.POST("/semesters", semesterHandler.Create)
	admin.PUT("/semesters/:id", semesterHandler.Update)
	admin.DELETE("/semesters/:id", semesterHandler.Delete)

	admin.POST("/categories", categoryHandler.Create)
	admin.PUT("/categories/:id", categoryHandler.Update)
	admin.DELETE("/categories/:id", categoryHandler.Delete)

	admin.POST("/courses", courseHandler.Create)
	admin.PUT("/courses/:id", courseHandler.Update)
	admin.DELETE("/courses/:id", courseHandler.Delete)

	admin.POST("/courses/:id/prerequisites", prereqHandler.Create)
	admin.DELETE("/prerequisites/:id", prereqHandler.Delete)

	admin.POST("/courses/:id/offerings", offeringHandler.Create)
	admin.DELETE("/offerings/:id", offeringHandler.Delete)

	admin.PUT("/plans/:id/courses/:courseId/category", ccHandler.SetForPlan)

	admin.GET("/students", studentAdminHandler.List)
	admin.POST("/students", studentAdminHandler.Create)
	admin.PUT("/students/:id", studentAdminHandler.Update)
	admin.DELETE("/students/:id", studentAdminHandler.Delete)
	admin.GET("/students/:id/records", studentAdminHandler.Records)
	admin.POST("/students/:id/records", studentAdminHandler.CreateRecord)
	admin.PUT("/records/:id", studentAdminHandler.UpdateRecord)
	admin.DELETE("/records/:id", studentAdminHandler.DeleteRecord)
	admin.GET("/students/:id/audit", studentAdminHandler.Audit)
	admin.GET("/students/:id/recommendations", studentAdminHandler.Recommendations)
	admin.POST("/students/:id/records/import", studentAdminHandler.ImportRecords)

	admin.GET("/stats/overview", statsHandler.Overview)
	admin.GET("/stats/credit-gap", statsHandler.CreditGap)
	admin.GET("/security/logs", securityHandler.Logs)

	return r
}
