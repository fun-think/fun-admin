package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheManager 缓存管理器接口
type CacheManager interface {
	// Get 获取缓存值
	Get(ctx context.Context, key string) (interface{}, error)

	// Set 设置缓存值
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Delete 删除缓存
	Delete(ctx context.Context, key string) error

	// Exists 检查缓存是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// Flush 清空所有缓存
	Flush(ctx context.Context) error
}

// RedisCacheManager Redis缓存管理器
type RedisCacheManager struct {
	client *redis.Client
}

// NewRedisCacheManager 创建Redis缓存管理器
func NewRedisCacheManager(client *redis.Client) *RedisCacheManager {
	return &RedisCacheManager{
		client: client,
	}
}

// Get 获取缓存值
func (r *RedisCacheManager) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // 缓存未找到
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Set 设置缓存值
func (r *RedisCacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Delete 删除缓存
func (r *RedisCacheManager) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists 检查缓存是否存在
func (r *RedisCacheManager) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Flush 清空所有缓存
func (r *RedisCacheManager) Flush(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

// MemoryCacheManager 内存缓存管理器
type MemoryCacheManager struct {
	data map[string]interface{}
	ttl  map[string]time.Time
}

// NewMemoryCacheManager 创建内存缓存管理器
func NewMemoryCacheManager() *MemoryCacheManager {
	return &MemoryCacheManager{
		data: make(map[string]interface{}),
		ttl:  make(map[string]time.Time),
	}
}

// Get 获取缓存值
func (m *MemoryCacheManager) Get(ctx context.Context, key string) (interface{}, error) {
	// 检查是否过期
	if ttl, exists := m.ttl[key]; exists {
		if time.Now().After(ttl) {
			// 过期，删除缓存
			delete(m.data, key)
			delete(m.ttl, key)
			return nil, nil
		}
	}

	if val, exists := m.data[key]; exists {
		return val, nil
	}

	return nil, nil
}

// Set 设置缓存值
func (m *MemoryCacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.data[key] = value
	if expiration > 0 {
		m.ttl[key] = time.Now().Add(expiration)
	}
	return nil
}

// Delete 删除缓存
func (m *MemoryCacheManager) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	delete(m.ttl, key)
	return nil
}

// Exists 检查缓存是否存在
func (m *MemoryCacheManager) Exists(ctx context.Context, key string) (bool, error) {
	// 检查是否过期
	if ttl, exists := m.ttl[key]; exists {
		if time.Now().After(ttl) {
			// 过期，删除缓存
			delete(m.data, key)
			delete(m.ttl, key)
			return false, nil
		}
	}

	_, exists := m.data[key]
	return exists, nil
}

// Flush 清空所有缓存
func (m *MemoryCacheManager) Flush(ctx context.Context) error {
	m.data = make(map[string]interface{})
	m.ttl = make(map[string]time.Time)
	return nil
}
