package config

import (
	"fmt"
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
	fmt.Println("load conf file:", envConf)

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
func Init(configPath string) {
	Global = NewManager(configPath)

	// 验证配置
	if err := Global.Validate(); err != nil {
		panic("Configuration validation failed: " + err.Error())
	}
}

// GetManager 获取全局配置管理器
func GetManager() *Manager {
	if Global == nil {
		panic("Configuration manager not initialized. Call config.Init() first.")
	}
	return Global
}

// Get 获取配置值（全局函数）
func Get(key string) interface{} {
	return GetManager().Get(key)
}

// GetString 获取字符串配置（全局函数）
func GetString(key string) string {
	return GetManager().GetString(key)
}

// GetInt 获取整数配置（全局函数）
func GetInt(key string) int {
	return GetManager().GetInt(key)
}

// GetBool 获取布尔配置（全局函数）
func GetBool(key string) bool {
	return GetManager().GetBool(key)
}

// GetDuration 获取时间段配置（全局函数）
func GetDuration(key string) time.Duration {
	return GetManager().GetDuration(key)
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
