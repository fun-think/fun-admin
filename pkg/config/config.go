package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// NewConfig 创建配置实例（保持向后兼容）
func NewConfig(p string) *viper.Viper {
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = p
	}
	log.Printf("Loading config file: %s", envConf)

	manager := NewManager(envConf)
	return manager.GetViper()
}

// getConfig 内部配置加载函数（保持向后兼容）
func getConfig(path string) *viper.Viper {
	manager := NewManager(path)
	return manager.GetViper()
}

// Global 全局配置管理器实例
var Global *Manager

// Init 初始化全局配置管理器
func Init(configPath string) error {
	Global = NewManager(configPath)

	// 验证配置
	if err := Global.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	return nil
}

// GetManager 获取全局配置管理器
func GetManager() (*Manager, error) {
	if Global == nil {
		return nil, fmt.Errorf("configuration manager not initialized, call config.Init() first")
	}
	return Global, nil
}

// Get 获取配置值（全局函数）
func Get(key string) interface{} {
	manager, err := GetManager()
	if err != nil {
		return nil
	}
	return manager.Get(key)
}

// GetString 获取字符串配置（全局函数）
func GetString(key string) string {
	manager, err := GetManager()
	if err != nil {
		return ""
	}
	return manager.GetString(key)
}

// GetInt 获取整数配置（全局函数）
func GetInt(key string) int {
	manager, err := GetManager()
	if err != nil {
		return 0
	}
	return manager.GetInt(key)
}

// GetBool 获取布尔配置（全局函数）
func GetBool(key string) bool {
	manager, err := GetManager()
	if err != nil {
		return false
	}
	return manager.GetBool(key)
}

// GetDuration 获取时间段配置（全局函数）
func GetDuration(key string) time.Duration {
	manager, err := GetManager()
	if err != nil {
		return 0
	}
	return manager.GetDuration(key)
}

// GetStringSlice 获取字符串切片配置，全局方法
func GetStringSlice(key string) []string {
	manager, err := GetManager()
	if err != nil {
		return nil
	}
	return manager.GetStringSlice(key)
}

// IsProduction 检查是否为生产环境
func IsProduction() bool {
	return GetString("app.env") == "production"
}

// IsDevelopment 检查是否为开发环境
func IsDevelopment() bool {
	return GetString("app.env") == "development"
}

// IsDebug 检查是否开启调试模式
func IsDebug() bool {
	return GetBool("app.debug")
}
