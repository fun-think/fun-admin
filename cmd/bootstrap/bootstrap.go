package bootstrap

import (
	"fmt"
	"fun-admin/internal/resources"
	"fun-admin/pkg/admin"
	"fun-admin/pkg/admin/i18n"
	"fun-admin/pkg/app"
	"fun-admin/pkg/cache"
	"fun-admin/pkg/config"
	"fun-admin/pkg/container"
	"fun-admin/pkg/logger"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitI18n registers translation resources.
func InitI18n() {
	i18n.AddResource("zh-CN", i18n.ChineseResources)
	i18n.AddResource("en", i18n.EnglishResources)
}

// InitContainer builds the IoC container with all service providers.
func InitContainer(configPath string) *container.Container {
	c := container.New()
	providers := []app.ServiceProvider{
		&app.CoreServiceProvider{},
		&app.RepositoryServiceProvider{},
		&app.ServiceServiceProvider{},
		&app.HandlerServiceProvider{},
		&app.ServerServiceProvider{},
	}

	for _, provider := range providers {
		provider.Register(c)
	}
	for _, provider := range providers {
		if err := provider.Boot(c); err != nil {
			panic(fmt.Sprintf("Failed to boot service provider: %v", err))
		}
	}
	return c
}

// InitCache initializes cache based on config.
func InitCache(conf *viper.Viper) cache.CacheManager {
	redisAddr := conf.GetString("data.redis.addr")
	if redisAddr != "" {
		rdb := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: conf.GetString("data.redis.password"),
			DB:       conf.GetInt("data.redis.db"),
		})
		return cache.NewRedisCacheManager(rdb)
	}
	return cache.NewMemoryCacheManager()
}

// InitAdmin registers built-in admin resources.
func InitAdmin(container *container.Container, cache cache.CacheManager, logger *logger.Logger, db *gorm.DB, enforcer *casbin.SyncedEnforcer) {
	admin.GlobalResourceManager.Register(resources.NewOperationLogResource())
	admin.GlobalResourceManager.Register(resources.NewUserResource())
	admin.GlobalResourceManager.Register(resources.NewCrudTableResource())
	admin.GlobalResourceManager.Register(resources.NewDictionaryTypeResource())
	admin.GlobalResourceManager.Register(resources.NewDictionaryDataResource())
}

// EnsureStorage prepares upload directories.
func EnsureStorage(log *logger.Logger) {
	if err := os.MkdirAll("storage/uploads", os.ModePerm); err != nil {
		log.Error("创建上传目录失败", zap.Error(err))
	}
}

// LoadConfig loads configuration and validates it.
func LoadConfig(confPath string) *viper.Viper {
	conf := config.NewConfig(confPath)
	validator := config.NewConfigValidator(conf)
	if err := validator.Validate(); err != nil {
		fmt.Printf("配置校验失败: %v\n", err)
		os.Exit(1)
	}
	return conf
}
