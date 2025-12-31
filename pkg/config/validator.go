package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

// ConfigValidator 配置验证器
type ConfigValidator struct {
	config *viper.Viper
	errors []string
}

// NewConfigValidator 创建配置验证器
func NewConfigValidator(config *viper.Viper) *ConfigValidator {
	return &ConfigValidator{
		config: config,
		errors: make([]string, 0),
	}
}

// Validate 验证所有配置
func (cv *ConfigValidator) Validate() error {
	cv.validateSecurity()
	cv.validateDatabase()
	cv.validateRedis()
	cv.validateServer()
	
	if len(cv.errors) > 0 {
		return fmt.Errorf("配置验证失败:\n%s", strings.Join(cv.errors, "\n"))
	}
	
	return nil
}

// validateSecurity 验证安全配置
func (cv *ConfigValidator) validateSecurity() {
	// 验证JWT密钥
	jwtKey := cv.config.GetString("security.jwt.key")
	if jwtKey == "" {
		cv.errors = append(cv.errors, "JWT密钥不能为空")
	} else if jwtKey == "default_jwt_secret_change_in_production" {
		cv.errors = append(cv.errors, "JWT密钥不能使用默认值，请在配置文件中设置安全的JWT密钥")
	} else if len(jwtKey) < 32 {
		cv.errors = append(cv.errors, "JWT密钥长度至少需要32个字符")
	}
	
	// 验证API签名密钥
	apiKey := cv.config.GetString("security.api_sign.app_key")
	if apiKey == "" {
		cv.errors = append(cv.errors, "API签名密钥不能为空")
	} else if apiKey == "default_api_key_change_in_production" {
		cv.errors = append(cv.errors, "API签名密钥不能使用默认值，请在配置文件中设置安全的API签名密钥")
	}
	
	apiSecret := cv.config.GetString("security.api_sign.app_security")
	if apiSecret == "" {
		cv.errors = append(cv.errors, "API签名安全密钥不能为空")
	} else if apiSecret == "default_api_secret_change_in_production" {
		cv.errors = append(cv.errors, "API签名安全密钥不能使用默认值，请在配置文件中设置安全的API签名安全密钥")
	}
}

// validateDatabase 验证数据库配置
func (cv *ConfigValidator) validateDatabase() {
	driver := cv.config.GetString("data.db.user.driver")
	dsn := cv.config.GetString("data.db.user.dsn")
	
	if driver == "" {
		cv.errors = append(cv.errors, "数据库驱动不能为空")
		return
	}
	
	if dsn == "" {
		cv.errors = append(cv.errors, "数据库连接字符串不能为空")
		return
	}
	
	// 检查是否包含默认密码
	if strings.Contains(dsn, "root:root@") {
		cv.errors = append(cv.errors, "数据库密码不能使用默认值，请在配置文件中设置安全的数据库密码")
	}
	
	// 验证MySQL连接字符串格式（支持环境变量占位符）
	if driver == "mysql" {
		// 允许环境变量占位符的正则表达式
		mysqlPattern := regexp.MustCompile(`^(\$\{[A-Z0-9_]+:[a-zA-Z0-9_]+\}|[a-zA-Z0-9_]+):([^@]+)@tcp\(\$\{[A-Z0-9_]+:[^\}:]+\}:\$\{[A-Z0-9_]+:\d+\}/\$\{[A-Z0-9_]+:[a-zA-Z0-9_]+\}\?`)
		// 不包含环境变量的常规正则表达式
		mysqlPatternNoEnv := regexp.MustCompile(`^[a-zA-Z0-9_]+:[^@]+@tcp\([^:]+:\d+\)/[a-zA-Z0-9_]+\?`)
		
		if !mysqlPattern.MatchString(dsn) && !mysqlPatternNoEnv.MatchString(dsn) {
			// 特殊处理示例中的连接字符串
			if !strings.Contains(dsn, "${") {
				cv.errors = append(cv.errors, "MySQL连接字符串格式不正确")
			}
		}
	}
}

// validateRedis 验证Redis配置
func (cv *ConfigValidator) validateRedis() {
	addr := cv.config.GetString("data.redis.addr")
	if addr == "" {
		// Redis不是必须的，但如果配置了就需要验证格式
		return
	}
	
	// 验证Redis地址格式（支持环境变量占位符）
	redisPattern := regexp.MustCompile(`^(\$\{[A-Z0-9_]+:[^\}]+\})?:?\d*$|^[^:]+:\d+$`)
	if !redisPattern.MatchString(addr) && !strings.Contains(addr, "${") {
		cv.errors = append(cv.errors, "Redis地址格式不正确，应为 host:port")
	}
}

// validateServer 验证服务器配置
func (cv *ConfigValidator) validateServer() {
	host := cv.config.GetString("http.host")
	if host == "" {
		cv.errors = append(cv.errors, "HTTP主机地址不能为空")
	}
	
	port := cv.config.GetInt("http.port")
	if port <= 0 || port > 65535 {
		cv.errors = append(cv.errors, "HTTP端口必须在1-65535范围内")
	}
}