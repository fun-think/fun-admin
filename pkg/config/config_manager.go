package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Manager 配置管理器
type Manager struct {
	v    *viper.Viper
	mu   sync.RWMutex
	path string
}

// NewManager 创建配置管理器
func NewManager(configPath string) *Manager {
	m := &Manager{
		v:    viper.New(),
		path: configPath,
	}

	m.initDefaults()
	m.loadConfig(configPath)
	return m
}

// initDefaults 初始化默认配置
func (m *Manager) initDefaults() {
	// 服务器配置
	m.v.SetDefault("app.name", "fun-admin")
	m.v.SetDefault("app.version", "1.0.0")
	m.v.SetDefault("app.env", "development")
	m.v.SetDefault("app.debug", true)

	// HTTP服务器配置
	m.v.SetDefault("http.host", "0.0.0.0")
	m.v.SetDefault("http.port", 8080)
	m.v.SetDefault("http.read_timeout", "30s")
	m.v.SetDefault("http.write_timeout", "30s")
	m.v.SetDefault("http.idle_timeout", "120s")
	m.v.SetDefault("http.cors.allowed_origins", []string{})
	m.v.SetDefault("http.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	m.v.SetDefault("http.cors.allowed_headers", []string{"Content-Type", "Authorization", "Accept", "Origin", "X-Requested-With"})
	m.v.SetDefault("http.cors.allow_credentials", false)
	m.v.SetDefault("http.cors.max_age", "2h")

	// 数据库配置
	m.v.SetDefault("data.db.user.driver", "sqlite")
	m.v.SetDefault("data.db.user.dsn", "data/database.db")
	m.v.SetDefault("data.db.debug", false)
	m.v.SetDefault("data.db.max_open_conns", 100)
	m.v.SetDefault("data.db.max_idle_conns", 10)
	m.v.SetDefault("data.db.max_lifetime", "1h")
	m.v.SetDefault("data.db.max_idle_time", "30m")

	// Redis配置
	m.v.SetDefault("data.redis.addr", "")
	m.v.SetDefault("data.redis.password", "")
	m.v.SetDefault("data.redis.db", 0)
	m.v.SetDefault("redis.pool_size", 10)
	m.v.SetDefault("redis.min_idle_conns", 5)

	// JWT配置
	m.v.SetDefault("security.jwt.key", "your-secret-key")
	m.v.SetDefault("jwt.expires", "24h")
	m.v.SetDefault("jwt.refresh_expires", "168h")

	// 日志配置
	m.v.SetDefault("logger.level", "info")
	m.v.SetDefault("logger.encoding", "json")
	m.v.SetDefault("logger.output_paths", []string{"stdout"})
	m.v.SetDefault("logger.error_output_paths", []string{"stderr"})
	m.v.SetDefault("logger.file.enabled", false)
	m.v.SetDefault("logger.file.path", "logs/app.logger")
	m.v.SetDefault("logger.file.max_size", 100)
	m.v.SetDefault("logger.file.max_backups", 3)
	m.v.SetDefault("logger.file.max_age", 28)

	// 缓存配置
	m.v.SetDefault("cache.type", "memory")
	m.v.SetDefault("cache.ttl", "1h")
	m.v.SetDefault("cache.cleanup_interval", "10m")

	// 上传配置
	m.v.SetDefault("upload.path", "storage/uploads")
	m.v.SetDefault("upload.max_size", "10MB")
	m.v.SetDefault("upload.allowed_types", []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx"})

	// 邮件配置
	m.v.SetDefault("mail.driver", "smtp")
	m.v.SetDefault("mail.host", "smtp.gmail.com")
	m.v.SetDefault("mail.port", 587)
	m.v.SetDefault("mail.username", "")
	m.v.SetDefault("mail.password", "")
	m.v.SetDefault("mail.from_address", "noreply@example.com")
	m.v.SetDefault("mail.from_name", "Fun Admin")
}

// loadConfig 加载配置文件
func (m *Manager) loadConfig(configPath string) {
	if configPath == "" {
		configPath = "config/local.yml"
	}

	// 设置配置文件
	m.v.SetConfigFile(configPath)

	// 设置环境变量前缀
	m.v.SetEnvPrefix("FUN_ADMIN")
	m.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	m.v.AutomaticEnv()

	// 创建配置文件目录
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		log.Printf("Failed to create config directory: %v", err)
	}

	// 读取配置文件
	if err := m.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，使用默认配置并创建配置文件
			log.Printf("Config file not found, creating default config at: %s", configPath)
			m.createDefaultConfigFile(configPath)
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	}
}

// createDefaultConfigFile 创建默认配置文件
func (m *Manager) createDefaultConfigFile(configPath string) {
	if err := m.v.WriteConfigAs(configPath); err != nil {
		log.Printf("Failed to create default config file: %v", err)
	}
}

// Get 获取配置值
func (m *Manager) Get(key string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.Get(key)
}

// GetString 获取字符串配置
func (m *Manager) GetString(key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.GetString(key)
}

// GetInt 获取整数配置
func (m *Manager) GetInt(key string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.GetInt(key)
}

// GetBool 获取布尔配置
func (m *Manager) GetBool(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.GetBool(key)
}

// GetFloat64 获取浮点数配置
func (m *Manager) GetFloat64(key string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.GetFloat64(key)
}

// GetDuration 获取时间段配置
func (m *Manager) GetDuration(key string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.GetDuration(key)
}

// GetStringSlice 获取字符串切片配置
func (m *Manager) GetStringSlice(key string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.GetStringSlice(key)
}

// GetStringMap 获取字符串映射配置
func (m *Manager) GetStringMap(key string) map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.GetStringMap(key)
}

// Set 设置配置值
func (m *Manager) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.v.Set(key, value)
}

// IsSet 检查配置键是否存在
func (m *Manager) IsSet(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.IsSet(key)
}

// AllKeys 获取所有配置键
func (m *Manager) AllKeys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.AllKeys()
}

// AllSettings 获取所有配置
func (m *Manager) AllSettings() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.AllSettings()
}

// Unmarshal 将配置解析到结构体
func (m *Manager) Unmarshal(rawVal interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.Unmarshal(rawVal)
}

// UnmarshalKey 将指定键的配置解析到结构体
func (m *Manager) UnmarshalKey(key string, rawVal interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.v.UnmarshalKey(key, rawVal)
}

// WriteConfig 保存配置到文件
func (m *Manager) WriteConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.v.WriteConfig()
}

// WriteConfigAs 保存配置到指定文件
func (m *Manager) WriteConfigAs(filename string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.v.WriteConfigAs(filename)
}

// GetViper 获取原始viper实例（慎用）
func (m *Manager) GetViper() *viper.Viper {
	return m.v
}

// Reload 重新加载配置
func (m *Manager) Reload() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.v.ReadInConfig()
}

// Validate 验证配置
func (m *Manager) Validate() error {
	// 验证必要的配置项
	requiredKeys := []string{
		"app.name",
		"http.port",
		"security.jwt.key",
	}

	for _, key := range requiredKeys {
		if !m.v.IsSet(key) {
			return fmt.Errorf("required config key '%s' is not set", key)
		}
	}

	// 验证JWT secret长度
	if len(m.GetString("security.jwt.key")) < 32 {
		return fmt.Errorf("security.jwt.key must be at least 32 characters long")
	}

	// 验证端口范围
	port := m.GetInt("http.port")
	if port < 1 || port > 65535 {
		return fmt.Errorf("http.port must be between 1 and 65535")
	}

	return nil
}

// Environment 配置环境结构
type Environment struct {
	App      AppConfig      `mapstructure:"app"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	Database DatabaseConfig `mapstructure:"data"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"logger"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Mail     MailConfig     `mapstructure:"mail"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
	Debug   bool   `mapstructure:"debug"`
}

type HTTPConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	DB struct {
		User struct {
			Driver string `mapstructure:"driver"`
			DSN    string `mapstructure:"dsn"`
		} `mapstructure:"user"`
		Debug        bool          `mapstructure:"debug"`
		MaxOpenConns int           `mapstructure:"max_open_conns"`
		MaxIdleConns int           `mapstructure:"max_idle_conns"`
		MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
		MaxIdleTime  time.Duration `mapstructure:"max_idle_time"`
	} `mapstructure:"db"`
}

type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type JWTConfig struct {
	Secret         string        `mapstructure:"key"`
	Expires        time.Duration `mapstructure:"expires"`
	RefreshExpires time.Duration `mapstructure:"refresh_expires"`
}

type LogConfig struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
	File             struct {
		Enabled    bool   `mapstructure:"enabled"`
		Path       string `mapstructure:"path"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxAge     int    `mapstructure:"max_age"`
	} `mapstructure:"file"`
}

type CacheConfig struct {
	Type            string        `mapstructure:"type"`
	TTL             time.Duration `mapstructure:"ttl"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

type UploadConfig struct {
	Path         string   `mapstructure:"path"`
	MaxSize      string   `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
}

type MailConfig struct {
	Driver      string `mapstructure:"driver"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	FromAddress string `mapstructure:"from_address"`
	FromName    string `mapstructure:"from_name"`
}
