package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheStrategy 缓存策略
type CacheStrategy struct {
	cacheManager CacheManager
	redisClient  *redis.Client
}

// NewCacheStrategy 创建缓存策略
func NewCacheStrategy(cacheManager CacheManager, redisClient *redis.Client) *CacheStrategy {
	return &CacheStrategy{
		cacheManager: cacheManager,
		redisClient:  redisClient,
	}
}

// CacheUserPermissions 缓存用户权限
func (cs *CacheStrategy) CacheUserPermissions(ctx context.Context, userID uint, permissions []string) error {
	key := fmt.Sprintf("user:permissions:%d", userID)
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("序列化权限数据失败: %v", err)
	}

	// 缓存30分钟
	return cs.cacheManager.Set(ctx, key, string(data), 30*time.Minute)
}

// GetUserPermissions 获取用户权限缓存
func (cs *CacheStrategy) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	key := fmt.Sprintf("user:permissions:%d", userID)
	data, err := cs.cacheManager.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("缓存未找到")
	}

	var permissions []string
	dataStr, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("缓存数据类型错误")
	}
	if err := json.Unmarshal([]byte(dataStr), &permissions); err != nil {
		return nil, fmt.Errorf("反序列化权限数据失败: %v", err)
	}

	return permissions, nil
}

// CacheUserMenus 缓存用户菜单
func (cs *CacheStrategy) CacheUserMenus(ctx context.Context, userID uint, menus interface{}) error {
	key := fmt.Sprintf("user:menus:%d", userID)
	data, err := json.Marshal(menus)
	if err != nil {
		return fmt.Errorf("序列化菜单数据失败: %v", err)
	}

	// 缓存1小时
	return cs.cacheManager.Set(ctx, key, string(data), time.Hour)
}

// GetUserMenus 获取用户菜单缓存
func (cs *CacheStrategy) GetUserMenus(ctx context.Context, userID uint, dest interface{}) error {
	key := fmt.Sprintf("user:menus:%d", userID)
	data, err := cs.cacheManager.Get(ctx, key)
	if err != nil {
		return err
	}

	if data == nil {
		return fmt.Errorf("缓存未找到")
	}

	dataStr, ok := data.(string)
	if !ok {
		return fmt.Errorf("缓存数据类型错误")
	}

	return json.Unmarshal([]byte(dataStr), dest)
}

// CacheSystemConfig 缓存系统配置
func (cs *CacheStrategy) CacheSystemConfig(ctx context.Context, configType string, config interface{}) error {
	key := fmt.Sprintf("system:config:%s", configType)
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置数据失败: %v", err)
	}

	// 缓存24小时
	return cs.cacheManager.Set(ctx, key, string(data), 24*time.Hour)
}

// GetSystemConfig 获取系统配置缓存
func (cs *CacheStrategy) GetSystemConfig(ctx context.Context, configType string, dest interface{}) error {
	key := fmt.Sprintf("system:config:%s", configType)
	data, err := cs.cacheManager.Get(ctx, key)
	if err != nil {
		return err
	}

	if data == nil {
		return fmt.Errorf("缓存未找到")
	}

	dataStr, ok := data.(string)
	if !ok {
		return fmt.Errorf("缓存数据类型错误")
	}

	return json.Unmarshal([]byte(dataStr), dest)
}

// InvalidateUserCache 清除用户相关缓存
func (cs *CacheStrategy) InvalidateUserCache(ctx context.Context, userID uint) error {
	keys := []string{
		fmt.Sprintf("user:permissions:%d", userID),
		fmt.Sprintf("user:menus:%d", userID),
	}

	for _, key := range keys {
		if err := cs.cacheManager.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

// InvalidateSystemConfigCache 清除系统配置缓存
func (cs *CacheStrategy) InvalidateSystemConfigCache(ctx context.Context, configType string) error {
	key := fmt.Sprintf("system:config:%s", configType)
	return cs.cacheManager.Delete(ctx, key)
}

// CacheWithFallback 带回退机制的缓存获取
func (cs *CacheStrategy) CacheWithFallback(
	ctx context.Context,
	key string,
	expiration time.Duration,
	fallback func() (interface{}, error),
	dest interface{},
) error {
	// 尝试从缓存获取
	data, err := cs.cacheManager.Get(ctx, key)
	if err == nil && data != nil {
		dataStr, ok := data.(string)
		if !ok {
			return fmt.Errorf("缓存数据类型错误")
		}
		return json.Unmarshal([]byte(dataStr), dest)
	}

	// 缓存未命中，执行回退函数
	result, err := fallback()
	if err != nil {
		return fmt.Errorf("执行回退函数失败: %v", err)
	}

	// 将结果存入缓存
	cacheData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("序列化缓存数据失败: %v", err)
	}

	if err := cs.cacheManager.Set(ctx, key, string(cacheData), expiration); err != nil {
		// 缓存失败不影响主流程，只记录日志
		fmt.Printf("缓存设置失败: %v\n", err)
	}

	// 将结果赋值到目标
	return json.Unmarshal(cacheData, dest)
}
