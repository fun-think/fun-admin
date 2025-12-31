package app

import (
	"context"
	"fun-admin/internal/handler"
	"fun-admin/internal/job"
	"fun-admin/internal/repository"
	internalserver "fun-admin/internal/server"
	"fun-admin/internal/service"
	"fun-admin/pkg/admin"
	"fun-admin/pkg/cache"
	"fun-admin/pkg/config"
	"fun-admin/pkg/container"
	"fun-admin/pkg/jwt"
	"fun-admin/pkg/logger"
	httpserver "fun-admin/pkg/server/http"
	"fun-admin/pkg/sid"
	"time"

	server "fun-admin/pkg/server"

	"github.com/casbin/casbin/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CoreServiceProvider 核心服务提供者
type CoreServiceProvider struct{}

func (p *CoreServiceProvider) Register(c *container.Container) {
	// 注册配置
	c.Singleton("config", func(c *container.Container) *viper.Viper {
		return config.NewConfig("config/local.yml")
	})

	// 注册日志
	c.Singleton("logger", func(c *container.Container) *logger.Logger {
		conf := c.MustGet("config").(*viper.Viper)
		return logger.NewLogger(conf)
	})

	// 注册数据库
	c.Singleton("database", func(c *container.Container) *gorm.DB {
		conf := c.MustGet("config").(*viper.Viper)
		log := c.MustGet("logger").(*logger.Logger)
		return repository.NewDB(conf, log)
	})

	// 注册事务
	c.Singleton("transaction", func(c *container.Container) repository.Transaction {
		repo := c.MustGet("repository").(*repository.Repository)
		return repository.NewTransaction(repo)
	})

	// 注册 Casbin 执行器
	c.Singleton("enforcer", func(c *container.Container) *casbin.SyncedEnforcer {
		conf := c.MustGet("config").(*viper.Viper)
		db := c.MustGet("database").(*gorm.DB)
		log := c.MustGet("logger").(*logger.Logger)
		return repository.NewCasbinEnforcer(conf, log, db)
	})

	// 注册缓存管理器
	c.Singleton("cache", func(c *container.Container) cache.CacheManager {
		// 简化实现，使用内存缓存
		return cache.NewMemoryCacheManager()
	})

	// 注册 JWT
	c.Singleton("jwt", func(c *container.Container) *jwt.JWT {
		conf := c.MustGet("config").(*viper.Viper)
		return jwt.NewJwt(conf)
	})

	// 注册 SID
	c.Singleton("sid", func(c *container.Container) *sid.Sid {
		sidInstance, err := sid.NewSid()
		if err != nil {
			log := c.MustGet("logger").(*logger.Logger)
			log.Fatal("failed to create sid instance", zap.Error(err))
			return nil
		}
		return sidInstance
	})
}

func (p *CoreServiceProvider) Boot(c *container.Container) error {
	return nil
}

// RepositoryServiceProvider 仓储服务提供者
type RepositoryServiceProvider struct{}

func (p *RepositoryServiceProvider) Register(c *container.Container) {
	// 注册基础仓储
	c.Singleton("repository", func(c *container.Container) *repository.Repository {
		db := c.MustGet("database").(*gorm.DB)
		log := c.MustGet("logger").(*logger.Logger)
		enforcer := c.MustGet("enforcer").(*casbin.SyncedEnforcer)
		return repository.NewRepository(log, db, enforcer)
	})

	// 注册资源仓储
	c.Singleton("resource_repository", func(c *container.Container) *repository.ResourceRepository {
		repo := c.MustGet("repository").(*repository.Repository)
		return repository.NewResourceRepository(*repo)
	})

	// 注册角色仓储
	c.Singleton("role_repository", func(c *container.Container) repository.RoleRepository {
		log := c.MustGet("logger").(*logger.Logger)
		db := c.MustGet("database").(*gorm.DB)
		enforcer := c.MustGet("enforcer").(*casbin.SyncedEnforcer)
		return repository.NewRoleRepository(log, db, enforcer)
	})

	// 注册用户仓储
	c.Singleton("user_repository", func(c *container.Container) repository.UserRepository {
		log := c.MustGet("logger").(*logger.Logger)
		db := c.MustGet("database").(*gorm.DB)
		enforcer := c.MustGet("enforcer").(*casbin.SyncedEnforcer)
		return repository.NewUserRepository(log, db, enforcer)
	})

	// 注册菜单仓储
	c.Singleton("menu_repository", func(c *container.Container) repository.MenuRepository {
		log := c.MustGet("logger").(*logger.Logger)
		db := c.MustGet("database").(*gorm.DB)
		return repository.NewMenuRepository(log, db)
	})

	// 注册接口仓储
	c.Singleton("api_repository", func(c *container.Container) repository.ApiRepository {
		log := c.MustGet("logger").(*logger.Logger)
		db := c.MustGet("database").(*gorm.DB)
		return repository.NewApiRepository(log, db)
	})

	// 注册操作日志仓储
	c.Singleton("operation_log_repository", func(c *container.Container) repository.OperationLogRepository {
		log := c.MustGet("logger").(*logger.Logger)
		db := c.MustGet("database").(*gorm.DB)
		return repository.NewOperationLogRepository(log, db)
	})

	// 注册仪表板仓储
	c.Singleton("dashboard_repository", func(c *container.Container) repository.DashboardRepository {
		log := c.MustGet("logger").(*logger.Logger)
		db := c.MustGet("database").(*gorm.DB)
		return repository.NewDashboardRepository(log, db)
	})

	// 注册个人资料仓储
	c.Singleton("profile_repository", func(c *container.Container) repository.ProfileRepository {
		log := c.MustGet("logger").(*logger.Logger)
		db := c.MustGet("database").(*gorm.DB)
		return repository.NewProfileRepository(log, db)
	})

	// 注册字典仓储
	c.Singleton("dictionary_repository", func(c *container.Container) repository.DictionaryRepository {
		log := c.MustGet("logger").(*logger.Logger)
		db := c.MustGet("database").(*gorm.DB)
		return repository.NewDictionaryRepository(log, db)
	})

	// 注册配置仓储
	c.Singleton("config_repository", func(c *container.Container) repository.ConfigRepository {
		repo := c.MustGet("repository").(*repository.Repository)
		return repository.NewConfigRepository(repo)
	})

	// 注册权限仓储
	c.Singleton("permission_repository", func(c *container.Container) repository.PermissionRepository {
		log := c.MustGet("logger").(*logger.Logger)
		db := c.MustGet("database").(*gorm.DB)
		enforcer := c.MustGet("enforcer").(*casbin.SyncedEnforcer)
		return repository.NewPermissionRepository(log, db, enforcer)
	})

	// 注册登录仓储
	c.Singleton("login_repository", func(c *container.Container) repository.LoginRepository {
		log := c.MustGet("logger").(*logger.Logger)
		db := c.MustGet("database").(*gorm.DB)
		return repository.NewLoginRepository(log, db)
	})

}

func (p *RepositoryServiceProvider) Boot(c *container.Container) error {
	return nil
}

// ServiceServiceProvider 服务提供者
type ServiceServiceProvider struct{}

func (p *ServiceServiceProvider) Register(c *container.Container) {
	// 注册基础服务
	c.Singleton("base_service", func(c *container.Container) *service.Service {
		log := c.MustGet("logger").(*logger.Logger)
		sidObj := c.MustGet("sid").(*sid.Sid)
		jwtObj := c.MustGet("jwt").(*jwt.JWT)
		tm := c.MustGet("repository").(*repository.Repository)
		return service.NewService(tm, log, sidObj, jwtObj)
	})

	// 注册登录服务
	c.Singleton("login_service", func(c *container.Container) service.LoginService {
		baseService := c.MustGet("base_service").(*service.Service)
		userRepo := c.MustGet("user_repository").(repository.UserRepository)
		loginRepo := c.MustGet("login_repository").(repository.LoginRepository)
		return service.NewLoginService(baseService, userRepo, loginRepo)
	})

	// 注册用户服务
	c.Singleton("user_service", func(c *container.Container) service.UserService {
		baseService := c.MustGet("base_service").(*service.Service)
		userRepo := c.MustGet("user_repository").(repository.UserRepository)
		roleRepo := c.MustGet("role_repository").(repository.RoleRepository)
		permissionRepo := c.MustGet("permission_repository").(repository.PermissionRepository) // 添加权限仓库
		return service.NewUserService(baseService, userRepo, roleRepo, permissionRepo)
	})

	// 注册资源服务
	c.Singleton("resource_service", func(c *container.Container) *service.ResourceService {
		resourceRepo := c.MustGet("resource_repository").(*repository.ResourceRepository)
		resourceManager := admin.GlobalResourceManager
		cacheManager := c.MustGet("cache").(cache.CacheManager)
		return service.NewResourceService(resourceRepo, resourceManager, cacheManager)
	})

	// 注册API服务
	c.Singleton("api_service", func(c *container.Container) service.ApiService {
		baseService := c.MustGet("base_service").(*service.Service)
		apiRepo := c.MustGet("api_repository").(repository.ApiRepository)
		return service.NewApiService(baseService, apiRepo)
	})

	// 注册配置服务
	c.Singleton("config_service", func(c *container.Container) *service.ConfigService {
		repo := c.MustGet("repository").(*repository.Repository)
		configRepo := c.MustGet("config_repository").(repository.ConfigRepository)
		return service.NewConfigService(repo, configRepo)
	})

	// 注册菜单服务
	c.Singleton("menu_service", func(c *container.Container) service.MenuService {
		baseService := c.MustGet("base_service").(*service.Service)
		menuRepo := c.MustGet("menu_repository").(repository.MenuRepository)
		userRepo := c.MustGet("user_repository").(repository.UserRepository)
		permissionService := c.MustGet("permission_service").(service.PermissionService)
		return service.NewMenuService(baseService, menuRepo, userRepo, permissionService)
	})

	// 注册字典服务
	c.Singleton("dictionary_service", func(c *container.Container) *service.DictionaryService {
		dictionaryRepo := c.MustGet("dictionary_repository").(repository.DictionaryRepository)
		log := c.MustGet("logger").(*logger.Logger)
		cacheObj := c.MustGet("cache").(cache.CacheManager)
		return service.NewDictionaryService(dictionaryRepo, log, cacheObj)
	})

	// 注册文件服务
	c.Singleton("file_service", func(c *container.Container) *service.FileService {
		log := c.MustGet("logger").(*logger.Logger)
		conf := c.MustGet("config").(*viper.Viper)
		return service.NewFileService(log, conf)
	})

	// 注册导入服务
	c.Singleton("import_service", func(c *container.Container) *service.ImportService {
		log := c.MustGet("logger").(*logger.Logger)
		return service.NewImportService(log)
	})

	// 注册权限服务
	c.Singleton("permission_service", func(c *container.Container) service.PermissionService {
		baseService := c.MustGet("base_service").(*service.Service)
		roleRepo := c.MustGet("role_repository").(repository.RoleRepository)
		userRepo := c.MustGet("user_repository").(repository.UserRepository)
		permissionRepo := c.MustGet("permission_repository").(repository.PermissionRepository) // 添加权限仓库
		return service.NewPermissionService(baseService, roleRepo, userRepo, permissionRepo)
	})

	// 注册角色服务
	c.Singleton("role_service", func(c *container.Container) service.RoleService {
		baseService := c.MustGet("base_service").(*service.Service)
		roleRepo := c.MustGet("role_repository").(repository.RoleRepository)
		permissionRepo := c.MustGet("permission_repository").(repository.PermissionRepository) // 添加权限仓库
		return service.NewRoleService(baseService, roleRepo, permissionRepo)
	})

	// 注册操作日志服务
	c.Singleton("operation_log_service", func(c *container.Container) service.OperationLogServiceInterface {
		baseService := c.MustGet("base_service").(*service.Service)
		operationLogRepo := c.MustGet("operation_log_repository").(repository.OperationLogRepository)
		return service.NewOperationLogService(baseService, operationLogRepo)
	})

	// 注册仪表板服务
	c.Singleton("dashboard_service", func(c *container.Container) service.DashboardServiceInterface {
		baseService := c.MustGet("base_service").(*service.Service)
		log := c.MustGet("logger").(*logger.Logger)
		dashboardRepo := c.MustGet("dashboard_repository").(repository.DashboardRepository)
		cacheObj := c.MustGet("cache").(cache.CacheManager)
		return service.NewDashboardService(baseService, log.Logger, dashboardRepo, cacheObj)
	})

	// 注册个人资料服务
	c.Singleton("profile_service", func(c *container.Container) service.ProfileServiceInterface {
		baseService := c.MustGet("base_service").(*service.Service)
		log := c.MustGet("logger").(*logger.Logger)
		profileRepo := c.MustGet("profile_repository").(repository.ProfileRepository)
		return service.NewProfileService(baseService, log.Logger, profileRepo)
	})
}

func (p *ServiceServiceProvider) Boot(c *container.Container) error {
	return nil
}

// HandlerServiceProvider 处理器服务提供者
type HandlerServiceProvider struct{}

func (p *HandlerServiceProvider) Register(c *container.Container) {
	// 注册基础处理器
	c.Singleton("handler", func(c *container.Container) *handler.Handler {
		log := c.MustGet("logger").(*logger.Logger)
		return handler.NewHandler(log)
	})

	// 注册登录处理器
	c.Singleton("login_handler", func(c *container.Container) *handler.LoginHandler {
		handlerInstance := c.MustGet("handler").(*handler.Handler)
		loginService := c.MustGet("login_service").(service.LoginService)
		return handler.NewLoginHandler(handlerInstance, loginService)
	})

	// 注册用户处理器
	c.Singleton("user_handler", func(c *container.Container) *handler.UserHandler {
		handlerInstance := c.MustGet("handler").(*handler.Handler)
		userService := c.MustGet("user_service").(service.UserService)
		return handler.NewUserHandler(handlerInstance, userService)
	})

	// 注册资源处理器
	c.Singleton("resource_handler", func(c *container.Container) *handler.ResourceHandler {
		resourceRepo := c.MustGet("resource_repository").(*repository.ResourceRepository)
		return handler.NewResourceHandler(resourceRepo)
	})

	// 注册导出处理器
	c.Singleton("export_handler", func(c *container.Container) *handler.ExportHandler {
		handlerInstance := c.MustGet("handler").(*handler.Handler)
		resourceService := c.MustGet("resource_service").(*service.ResourceService)
		log := c.MustGet("logger").(*logger.Logger)
		return handler.NewExportHandler(handlerInstance, resourceService, log.Logger)
	})

	// 注册API处理器
	c.Singleton("api_handler", func(c *container.Container) *handler.ApiHandler {
		handlerInstance := c.MustGet("handler").(*handler.Handler)
		apiService := c.MustGet("api_service").(service.ApiService)
		return handler.NewApiHandler(handlerInstance, apiService)
	})

	// 注册配置处理器
	c.Singleton("config_handler", func(c *container.Container) *handler.ConfigHandler {
		configService := c.MustGet("config_service").(*service.ConfigService)
		return handler.NewConfigHandler(configService)
	})

	// 注册菜单处理器
	c.Singleton("menu_handler", func(c *container.Container) *handler.MenuHandler {
		handlerInstance := c.MustGet("handler").(*handler.Handler)
		menuService := c.MustGet("menu_service").(service.MenuService)
		return handler.NewMenuHandler(handlerInstance, menuService)
	})

	// 注册字典处理器
	c.Singleton("dictionary_handler", func(c *container.Container) *handler.DictionaryHandler {
		dictionaryService := c.MustGet("dictionary_service").(*service.DictionaryService)
		return handler.NewDictionaryHandler(dictionaryService)
	})

	// 注册文件处理器
	c.Singleton("file_handler", func(c *container.Container) *handler.FileHandler {
		log := c.MustGet("logger").(*logger.Logger)
		fileService := c.MustGet("file_service").(*service.FileService)
		return handler.NewFileHandler(log, fileService)
	})

	// 注册导入处理器
	c.Singleton("import_handler", func(c *container.Container) *handler.ImportHandler {
		importService := c.MustGet("import_service").(*service.ImportService)
		resourceService := c.MustGet("resource_service").(*service.ResourceService)
		return handler.NewImportHandler(importService, resourceService)
	})

	// 注册权限处理器
	c.Singleton("permission_handler", func(c *container.Container) *handler.PermissionHandler {
		handlerInstance := c.MustGet("handler").(*handler.Handler)
		permissionService := c.MustGet("permission_service").(service.PermissionService)
		return handler.NewPermissionHandler(handlerInstance, permissionService)
	})

	// 注册个人资料处理器
	c.Singleton("profile_handler", func(c *container.Container) *handler.ProfileHandler {
		handlerInstance := c.MustGet("handler").(*handler.Handler)
		profileService := c.MustGet("profile_service").(service.ProfileServiceInterface)
		return handler.NewProfileHandler(handlerInstance, profileService)
	})

	// 注册角色处理器
	c.Singleton("role_handler", func(c *container.Container) *handler.RoleHandler {
		handlerInstance := c.MustGet("handler").(*handler.Handler)
		roleService := c.MustGet("role_service").(service.RoleService)
		return handler.NewRoleHandler(handlerInstance, roleService)
	})

	// 注册仪表盘处理器
	c.Singleton("dashboard_handler", func(c *container.Container) *handler.DashboardHandler {
		handlerInstance := c.MustGet("handler").(*handler.Handler)
		dashboardService := c.MustGet("dashboard_service").(service.DashboardServiceInterface)
		return handler.NewDashboardHandler(handlerInstance, dashboardService)
	})

	// 注册操作日志处理器
	c.Singleton("operation_log_handler", func(c *container.Container) *handler.OperationLogHandler {
		handlerInstance := c.MustGet("handler").(*handler.Handler)
		operationLogService := c.MustGet("operation_log_service").(service.OperationLogServiceInterface)
		return handler.NewOperationLogHandler(handlerInstance, operationLogService)
	})

	// 注册资源CRUD处理器
	c.Singleton("resource_crud_handler", func(c *container.Container) *handler.ResourceCRUDHandler {
		resourceService := c.MustGet("resource_service").(*service.ResourceService)
		return handler.NewResourceCRUDHandler(resourceService)
	})
}

func (p *HandlerServiceProvider) Boot(c *container.Container) error {
	return nil
}

// 定义缓存适配器结构体
type cacheAdapter struct {
	cm cache.CacheManager
}

// 实现 Get 方法
func (ca *cacheAdapter) Get(ctx context.Context, key string) (interface{}, error) {
	return ca.cm.Get(ctx, key)
}

// 实现 Set 方法
func (ca *cacheAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return ca.cm.Set(ctx, key, value, expiration)
}

// 实现 Delete 方法
func (ca *cacheAdapter) Delete(ctx context.Context, key string) error {
	return ca.cm.Delete(ctx, key)
}

// 实现 Exists 方法
func (ca *cacheAdapter) Exists(ctx context.Context, key string) (bool, error) {
	return ca.cm.Exists(ctx, key)
}

// 实现 DeleteByPrefix 方法
func (ca *cacheAdapter) DeleteByPrefix(ctx context.Context, prefix string) error {
	return ca.cm.DeleteByPrefix(ctx, prefix)
}

// 实现 Flush 方法
func (ca *cacheAdapter) Flush(ctx context.Context) error {
	return ca.cm.Flush(ctx)
}

// ServerServiceProvider 服务器服务提供者
type ServerServiceProvider struct{}

func (p *ServerServiceProvider) Register(c *container.Container) {
	// 注册HTTP服务器
	c.Singleton("http_server", func(c *container.Container) server.Server {
		log := c.MustGet("logger").(*logger.Logger)
		conf := c.MustGet("config").(*viper.Viper)
		jwtObj := c.MustGet("jwt").(*jwt.JWT)
		enforcer := c.MustGet("enforcer").(*casbin.SyncedEnforcer)
		// 创建 HTTP 服务器
		engine := internalserver.NewHTTPServer(
			log,
			enforcer,
			jwtObj,
			c, // 传递容器，让 NewHTTPServer 自己从容器中获取 handler
		)

		// 包装为实现了 server.Server 接口的服务器
		httpServer := httpserver.NewServer(
			engine,
			log,
			httpserver.WithServerHost(conf.GetString("http.host")),
			httpserver.WithServerPort(conf.GetInt("http.port")),
		)

		return httpServer
	})

	// 注册Job服务器
	c.Singleton("job_server", func(c *container.Container) server.Server {
		log := c.MustGet("logger").(*logger.Logger)
		transaction := c.MustGet("transaction").(repository.Transaction)
		sidObj := c.MustGet("sid").(*sid.Sid)
		userRepo := c.MustGet("user_repository").(repository.UserRepository)

		baseJob := job.NewJob(transaction, log, sidObj)
		userJob := job.NewUserJob(baseJob, userRepo)
		jobServer := internalserver.NewJobServer(log, userJob)

		return jobServer
	})
}

func (p *ServerServiceProvider) Boot(c *container.Container) error {
	return nil
}
