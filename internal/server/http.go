package server

import (
	"fun-admin/internal/handler"
	"fun-admin/internal/middleware"
	"fun-admin/internal/repository"
	"fun-admin/pkg/container"
	"fun-admin/pkg/jwt"
	"fun-admin/pkg/logger"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// NewHTTPServer 创建 HTTP 服务器
func NewHTTPServer(
	logger *logger.Logger,
	enforcer *casbin.SyncedEnforcer,
	jwt *jwt.JWT,
	c *container.Container,
) *gin.Engine {
	conf := c.MustGet("config").(*viper.Viper)

	// 根据环境切换 Gin 模式（默认 Debug，生产环境 Release）
	env := strings.ToLower(conf.GetString("env"))
	if env == "" {
		env = strings.ToLower(conf.GetString("app.env"))
	}
	if env == "prod" || env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化 gin 引擎
	app := gin.New()

	// 从容器中获取 handler
	loginHandler := c.MustGet("login_handler").(*handler.LoginHandler)
	userHandler := c.MustGet("user_handler").(*handler.UserHandler)
	exportHandler := c.MustGet("export_handler").(*handler.ExportHandler)
	dashboardHandler := c.MustGet("dashboard_handler").(*handler.DashboardHandler)
	operationLogHandler := c.MustGet("operation_log_handler").(*handler.OperationLogHandler)
	profileHandler := c.MustGet("profile_handler").(*handler.ProfileHandler)
	apiHandler := c.MustGet("api_handler").(*handler.ApiHandler)
	configHandler := c.MustGet("config_handler").(*handler.ConfigHandler)
	menuHandler := c.MustGet("menu_handler").(*handler.MenuHandler)
	dictionaryHandler := c.MustGet("dictionary_handler").(*handler.DictionaryHandler)
	fileHandler := c.MustGet("file_handler").(*handler.FileHandler)
	importHandler := c.MustGet("import_handler").(*handler.ImportHandler)
	permissionHandler := c.MustGet("permission_handler").(*handler.PermissionHandler)
	roleHandler := c.MustGet("role_handler").(*handler.RoleHandler)
	resourceHandler := c.MustGet("resource_handler").(*handler.ResourceHandler)
	resourceCRUDHandler := c.MustGet("resource_crud_handler").(*handler.ResourceCRUDHandler)
	repo := c.MustGet("repository").(*repository.Repository)
	db := c.MustGet("database").(*gorm.DB)
	mwManager := middleware.NewManager(logger, db, enforcer, repo, conf)

	// 注册中间件
	registerMiddlewares(app, mwManager)

	// 注册管理后台路由
	registerAdminRoutes(
		app,
		enforcer,
		jwt,
		dashboardHandler,
		operationLogHandler,
		profileHandler,
		exportHandler,
		userHandler,
		apiHandler,
		configHandler,
		menuHandler,
		dictionaryHandler,
		fileHandler,
		importHandler,
		permissionHandler,
		roleHandler,
		resourceHandler,
		resourceCRUDHandler,
		loginHandler,
		logger,
	)

	return app
}

// registerMiddlewares 注册中间件
func registerMiddlewares(app *gin.Engine, manager *middleware.Manager) {
	app.Use(manager.SetupRequestLogMiddleware())
	app.Use(manager.SetupOperationLogMiddleware())
	app.Use(gin.Recovery())
	app.Use(manager.SetupCORSMiddleware())
}

// registerAdminRoutes 注册管理后台路由
func registerAdminRoutes(
	app *gin.Engine,
	enforcer *casbin.SyncedEnforcer,
	jwt *jwt.JWT,
	dashboardHandler *handler.DashboardHandler,
	operationLogHandler *handler.OperationLogHandler,
	profileHandler *handler.ProfileHandler,
	exportHandler *handler.ExportHandler,
	// 新增的 Handler
	userHandler *handler.UserHandler,
	apiHandler *handler.ApiHandler,
	configHandler *handler.ConfigHandler,
	menuHandler *handler.MenuHandler,
	dictionaryHandler *handler.DictionaryHandler,
	fileHandler *handler.FileHandler,
	importHandler *handler.ImportHandler,
	permissionHandler *handler.PermissionHandler,
	roleHandler *handler.RoleHandler,
	resourceHandler *handler.ResourceHandler,
	resourceCRUDHandler *handler.ResourceCRUDHandler,
	// 公共路由需要的 Handler
	loginHandler *handler.LoginHandler,
	logger *logger.Logger,
) {
	// 将登录接口移出需要权限的路由组
	app.POST("/api/admin/login", loginHandler.Login)

	adminGroup := app.Group("/api/admin")
	// 为 admin API 添加权限中间件
	adminGroup.Use(middleware.Jwt(jwt, logger), middleware.PermissionMiddleware(enforcer))
	{
		// 注册基础路由
		adminGroup.GET("/v1/dashboard", dashboardHandler.GetDashboard)
		adminGroup.GET("/v1/operation-logs", operationLogHandler.GetOperationLogs)
		adminGroup.GET("/v1/operation-logs/stats", operationLogHandler.GetOperationLogStats)
		adminGroup.GET("/v1/operation-logs/:id", operationLogHandler.GetOperationLog)
		adminGroup.DELETE("/v1/operation-logs/:id", operationLogHandler.DeleteOperationLog)
		adminGroup.DELETE("/v1/operation-logs", operationLogHandler.BatchDeleteOperationLogs)
		adminGroup.DELETE("/v1/operation-logs/clear", operationLogHandler.ClearOperationLogs)
		adminGroup.GET("/v1/profile", profileHandler.GetProfile)
		adminGroup.PUT("/v1/profile", profileHandler.UpdateProfile)
		adminGroup.PUT("/v1/profile/password", profileHandler.UpdatePassword)

		// 注册导出路由
		adminGroup.GET("/v1/export/:resource", exportHandler.ExportData)

		// 用户管理相关接口
		adminGroup.GET("/v1/users", userHandler.GetUsers)
		adminGroup.POST("/v1/users", userHandler.UserCreate)
		adminGroup.GET("/v1/users/:id", userHandler.GetUser)
		adminGroup.PUT("/v1/users/:id", userHandler.UserUpdate)
		adminGroup.DELETE("/v1/users/:id", userHandler.UserDelete)

		// API管理相关接口
		adminGroup.GET("/v1/apis", apiHandler.GetApis)
		adminGroup.POST("/v1/apis", apiHandler.ApiCreate)
		adminGroup.GET("/v1/apis/:id", apiHandler.ApiGet)
		adminGroup.PUT("/v1/apis/:id", apiHandler.ApiUpdate)
		adminGroup.DELETE("/v1/apis/:id", apiHandler.ApiDelete)

		// 配置管理相关接口
		adminGroup.GET("/v1/configs", configHandler.GetAllConfigs)
		adminGroup.POST("/v1/configs", configHandler.SetConfig)
		adminGroup.GET("/v1/configs/:id", configHandler.GetConfig)
		adminGroup.PUT("/v1/configs/:id", configHandler.UpdateConfig)
		adminGroup.DELETE("/v1/configs/:id", configHandler.DeleteConfig)

		// 菜单管理相关接口
		adminGroup.GET("/v1/menus", menuHandler.GetAdminMenus)
		adminGroup.POST("/v1/menus", menuHandler.MenuCreate)
		adminGroup.GET("/v1/menus/:id", menuHandler.GetMenus)
		adminGroup.PUT("/v1/menus/:id", menuHandler.MenuUpdate)
		adminGroup.DELETE("/v1/menus/:id", menuHandler.MenuDelete)

		// 字典管理相关接口
		adminGroup.GET("/v1/dictionaries", dictionaryHandler.ListDictionaryTypes)
		adminGroup.POST("/v1/dictionaries", dictionaryHandler.CreateDictionaryType)
		adminGroup.GET("/v1/dictionaries/:id", dictionaryHandler.GetDictionaryType)
		adminGroup.PUT("/v1/dictionaries/:id", dictionaryHandler.UpdateDictionaryType)
		// 文件服务接口
		adminGroup.POST("/v1/upload", fileHandler.Upload)
		adminGroup.GET("/v1/files", fileHandler.List)
		adminGroup.GET("/v1/files/info", fileHandler.Info)
		adminGroup.DELETE("/v1/files", fileHandler.Delete)

		// 导入管理相关接口
		adminGroup.POST("/v1/import/:resource", importHandler.ImportData)

		// 权限管理相关接口
		adminGroup.GET("/v1/permissions", permissionHandler.GetUserPermissions)
		adminGroup.POST("/v1/permissions", permissionHandler.UpdateRolePermission)
		adminGroup.GET("/v1/permissions/role", permissionHandler.GetRolePermissions)
		adminGroup.PUT("/v1/permissions/role", permissionHandler.UpdateRolePermission)
		adminGroup.GET("/v1/permissions/user", permissionHandler.GetPermissionsForUser)
		adminGroup.POST("/v1/permissions/user", permissionHandler.AddRoleForUser)
		adminGroup.DELETE("/v1/permissions/user", permissionHandler.DeleteRoleForUser)

		// 角色管理相关接口
		adminGroup.GET("/v1/roles", roleHandler.GetRoles)
		adminGroup.POST("/v1/roles", roleHandler.RoleCreate)
		adminGroup.GET("/v1/roles/:id", roleHandler.GetRoles)
		adminGroup.PUT("/v1/roles/:id", roleHandler.RoleUpdate)
		adminGroup.DELETE("/v1/roles/:id", roleHandler.RoleDelete)
		adminGroup.GET("/v1/roles/all", permissionHandler.GetAllRoles)

		// 资源管理相关接口
		adminGroup.GET("/v1/resources", resourceHandler.ListResources)
		adminGroup.GET("/v1/resources/search", resourceHandler.GlobalSearch)
		adminGroup.GET("/v1/resources/:slug", resourceHandler.GetResource)

		// 资源 CRUD 相关接口
		adminGroup.GET("/v1/resource-crud/:resource", resourceCRUDHandler.List)
		adminGroup.POST("/v1/resource-crud/:resource", resourceCRUDHandler.Create)
		adminGroup.GET("/v1/resource-crud/:resource/:id", resourceCRUDHandler.Get)
		adminGroup.PUT("/v1/resource-crud/:resource/:id", resourceCRUDHandler.Update)
		adminGroup.DELETE("/v1/resource-crud/:resource/:id", resourceCRUDHandler.Delete)
		adminGroup.POST("/v1/resource-crud/:resource/actions/:action", resourceCRUDHandler.RunAction)

		// 添加 ping 接口
		adminGroup.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// 登录接口已移到权限中间件外
		// adminGroup.POST("/login", loginHandler.Login)

		// 静态文件服务
		adminGroup.Static("/v1/static", "./static")
		adminGroup.Static("/v1/upload", "./upload")
	}
}
