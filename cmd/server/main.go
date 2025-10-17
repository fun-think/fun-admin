package main

import (
	"context"
	"flag"
	"fmt"
	"fun-admin/internal/resources"
	"fun-admin/pkg/admin"
	"fun-admin/pkg/admin/i18n"
	"fun-admin/pkg/app"
	"fun-admin/pkg/cache"
	"fun-admin/pkg/config"
	"fun-admin/pkg/container"
	"fun-admin/pkg/logger"
	"fun-admin/pkg/server"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// @title           Nunu Example API
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8000
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	log := logger.NewLogger(conf)

	// 确保上传目录存在
	if err := os.MkdirAll("storage/uploads", os.ModePerm); err != nil {
		log.Error("创建上传目录失败", zap.Error(err))
	}

	// 初始化国际化资源
	initI18n()

	// 创建依赖注入容器
	containerManager := initContainer(*envConf)

	// 初始化缓存管理器
	cacheManager := initCache(conf)

	// 获取基础服务实例
	db := containerManager.MustGet("database").(*gorm.DB)
	syncedEnforcer := containerManager.MustGet("enforcer").(*casbin.SyncedEnforcer)

	// 初始化 admin 系统
	initAdmin(containerManager, cacheManager, log, db, syncedEnforcer)

	// 初始化服务器实例并直接作为 server.Server 接口使用
	httpServer := containerManager.MustGet("http_server").(server.Server)
	jobServer := containerManager.MustGet("job_server").(server.Server)

	var servers []server.Server
	servers = append(servers, httpServer, jobServer)

	appManager := app.NewApp(
		app.WithServer(servers...),
		app.WithName("fun-server"),
	)

	log.Info("server start", zap.String("host", fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))))
	log.Info("docs addr", zap.String("addr", fmt.Sprintf("http://%s:%d/swagger/index.html", conf.GetString("http.host"), conf.GetInt("http.port"))))
	if err := appManager.Run(context.Background()); err != nil {
		panic(err)
	}
}

// initAdmin 初始化 admin 系统
func initAdmin(container *container.Container, cache cache.CacheManager, logger *logger.Logger, db *gorm.DB, enforcer *casbin.SyncedEnforcer) {
	// 注册资源
	admin.GlobalResourceManager.Register(resources.NewOperationLogResource())
	admin.GlobalResourceManager.Register(resources.NewUserResource())

	// 注册自定义页面

	// 注意：所有 handler 现在都通过容器注册和获取，无需在此处手动初始化
}

// initI18n 初始化国际化资源
func initI18n() {
	// 添加中文资源
	i18n.AddResource("zh-CN", i18n.ChineseResources)

	// 添加英文资源
	i18n.AddResource("en", i18n.EnglishResources)
}

// initContainer 初始化依赖注入容器
func initContainer(configPath string) *container.Container {
	c := container.New()

	// 注册服务提供者
	providers := []app.ServiceProvider{
		&app.CoreServiceProvider{},
		&app.RepositoryServiceProvider{},
		&app.ServiceServiceProvider{},
		&app.HandlerServiceProvider{},
		&app.ServerServiceProvider{},
	}

	// 注册所有服务
	for _, provider := range providers {
		provider.Register(c)
	}

	// 启动所有服务
	for _, provider := range providers {
		if err := provider.Boot(c); err != nil {
			panic(fmt.Sprintf("Failed to boot service provider: %v", err))
		}
	}

	return c
}

// initCache 初始化缓存管理器
func initCache(conf *viper.Viper) cache.CacheManager {
	// 检查是否配置了Redis
	redisAddr := conf.GetString("redis.addr")
	if redisAddr != "" {
		// 使用Redis缓存
		rdb := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: conf.GetString("redis.password"),
			DB:       conf.GetInt("redis.db"),
		})

		return cache.NewRedisCacheManager(rdb)
	}

	// 使用内存缓存
	return cache.NewMemoryCacheManager()
}
