package server

import (
	"fun-admin/internal/handler"
	"fun-admin/internal/middleware"
	"fun-admin/pkg/cache"
	"fun-admin/pkg/container"
	"fun-admin/pkg/jwt"
	"fun-admin/pkg/logger"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// NewHTTPServer 创建 HTTP 服务器
func NewHTTPServer(
	logger *logger.Logger,
	enforcer *casbin.SyncedEnforcer,
	jwt *jwt.JWT,
	c *container.Container,
) *gin.Engine {
	// 设置为发布模式
	gin.SetMode(gin.DebugMode)

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
	cacheManager := c.MustGet("cache").(cache.CacheManager)

	// 注册中间件
	registerMiddlewares(app, logger, cacheManager)

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
func registerMiddlewares(app *gin.Engine, logger *logger.Logger, cache cache.CacheManager) {
	// 日志中间件
	//app.Use(middleware.RequestLogMiddleware(logger))

	// Recovery 中间件
	app.Use(gin.Recovery())

	// 跨域中间件
	app.Use(middleware.CORSMiddleware())

	// 安全中间件
	//app.Use(middleware.SecurityMiddleware())

	// 请求频率限制中间件
	// app.Use(middleware.RateLimitMiddleware(100, time.Minute))
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
	// 注册资源管理路由
	// adminGroup := app.Group("/api/admin")
	// 为 admin API 添加权限中间件
	// adminGroup.Use(middleware.PermissionMiddleware(enforcer))

	// 将登录接口移出需要权限的路由组
	app.POST("/api/admin/login", loginHandler.Login)

	adminGroup := app.Group("/api/admin")
	// 为 admin API 添加权限中间件
	adminGroup.Use(middleware.Jwt(jwt, logger), middleware.PermissionMiddleware(enforcer))
	{
		// 注册基础路由
		adminGroup.GET("/dashboard", dashboardHandler.GetDashboard)
		adminGroup.GET("/operation-logs", operationLogHandler.GetOperationLogs)
		adminGroup.GET("/profile", profileHandler.GetProfile)
		adminGroup.PUT("/profile", profileHandler.UpdateProfile)
		adminGroup.PUT("/profile/password", profileHandler.UpdatePassword)

		// 注册导出路由
		adminGroup.GET("/export/:resource", exportHandler.ExportData)

		// 用户管理相关接口
		adminGroup.GET("/users", userHandler.GetUsers)
		adminGroup.POST("/users", userHandler.UserCreate)
		adminGroup.GET("/users/:id", userHandler.GetUser)
		adminGroup.PUT("/users/:id", userHandler.UserUpdate)
		adminGroup.DELETE("/users/:id", userHandler.UserDelete)

		// API管理相关接口
		adminGroup.GET("/apis", apiHandler.GetApis)
		adminGroup.POST("/apis", apiHandler.ApiCreate)
		adminGroup.GET("/apis/:id", apiHandler.GetApis) // 这里可能需要添加一个获取单个API的方法
		adminGroup.PUT("/apis/:id", apiHandler.ApiUpdate)
		adminGroup.DELETE("/apis/:id", apiHandler.ApiDelete)

		// 配置管理相关接口
		adminGroup.GET("/configs", configHandler.GetAllConfigs)
		adminGroup.POST("/configs", configHandler.SetConfig)
		adminGroup.GET("/configs/:id", configHandler.GetConfig)
		adminGroup.PUT("/configs/:id", configHandler.UpdateConfig)
		adminGroup.DELETE("/configs/:id", configHandler.DeleteConfig)

		// 菜单管理相关接口
		adminGroup.GET("/menus", menuHandler.GetAdminMenus)
		adminGroup.POST("/menus", menuHandler.MenuCreate)
		adminGroup.GET("/menus/:id", menuHandler.GetMenus)
		adminGroup.PUT("/menus/:id", menuHandler.MenuUpdate)
		adminGroup.DELETE("/menus/:id", menuHandler.MenuDelete)

		// 字典管理相关接口
		adminGroup.GET("/dictionaries", dictionaryHandler.ListDictionaryTypes)
		adminGroup.POST("/dictionaries", dictionaryHandler.CreateDictionaryType)
		adminGroup.GET("/dictionaries/:id", dictionaryHandler.GetDictionaryType)
		adminGroup.PUT("/dictionaries/:id", dictionaryHandler.UpdateDictionaryType)
		adminGroup.DELETE("/dictionaries/:id", dictionaryHandler.DeleteDictionaryType)

		// 文件管理相关接口
		adminGroup.POST("/upload", handler.UploadHandler)
		// 注意：文件服务中的方法需要进一步实现

		// 导入管理相关接口
		adminGroup.POST("/import/:resource", importHandler.ImportData)

		// 权限管理相关接口
		adminGroup.GET("/permissions", permissionHandler.GetUserPermissions)
		adminGroup.POST("/permissions", permissionHandler.UpdateRolePermission)
		adminGroup.GET("/permissions/role", permissionHandler.GetRolePermissions)
		adminGroup.PUT("/permissions/role", permissionHandler.UpdateRolePermission)
		adminGroup.GET("/permissions/user", permissionHandler.GetPermissionsForUser)
		adminGroup.POST("/permissions/user", permissionHandler.AddRoleForUser)
		adminGroup.DELETE("/permissions/user", permissionHandler.DeleteRoleForUser)

		// 角色管理相关接口
		adminGroup.GET("/roles", roleHandler.GetRoles)
		adminGroup.POST("/roles", roleHandler.RoleCreate)
		adminGroup.GET("/roles/:id", roleHandler.GetRoles)
		adminGroup.PUT("/roles/:id", roleHandler.RoleUpdate)
		adminGroup.DELETE("/roles/:id", roleHandler.RoleDelete)
		adminGroup.GET("/roles/all", permissionHandler.GetAllRoles)

		// 资源管理相关接口
		adminGroup.GET("/resources", resourceHandler.ListResources)
		adminGroup.GET("/resources/search", resourceHandler.GlobalSearch)
		adminGroup.GET("/resources/:slug", resourceHandler.GetResource)

		// 资源 CRUD 相关接口
		adminGroup.GET("/resource-crud/:resource", resourceCRUDHandler.List)
		adminGroup.POST("/resource-crud/:resource", resourceCRUDHandler.Create)
		adminGroup.GET("/resource-crud/:resource/:id", resourceCRUDHandler.Get)
		adminGroup.PUT("/resource-crud/:resource/:id", resourceCRUDHandler.Update)
		adminGroup.DELETE("/resource-crud/:resource/:id", resourceCRUDHandler.Delete)

		// 添加 ping 接口
		adminGroup.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// 登录接口已移到权限中间件外
		// adminGroup.POST("/login", loginHandler.Login)

		// 静态文件服务
		adminGroup.Static("/static", "./static")
		adminGroup.Static("/upload", "./upload")
	}
}
