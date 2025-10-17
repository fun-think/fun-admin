package middleware

import (
	"fun-admin/api/v1"
	"fun-admin/pkg/errors"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(requestsPerMinute int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerMinute, window)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !limiter.Allow(clientIP) {
			err := errors.New(errors.CodeTooManyRequests, "请求过于频繁，请稍后再试")
			v1.HandleError(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimiter 限流器
type RateLimiter struct {
	mu       sync.RWMutex
	clients  map[string]*clientInfo
	maxReqs  int
	window   time.Duration
	cleanup  time.Duration
	stopChan chan struct{}
}

type clientInfo struct {
	requests []time.Time
	lastSeen time.Time
}

// NewRateLimiter 创建限流器
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*clientInfo),
		maxReqs:  maxRequests,
		window:   window,
		cleanup:  time.Minute * 5, // 每5分钟清理一次
		stopChan: make(chan struct{}),
	}

	go rl.cleanupClients()
	return rl
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// 获取或创建客户端信息
	client, exists := rl.clients[clientID]
	if !exists {
		client = &clientInfo{
			requests: make([]time.Time, 0, rl.maxReqs),
			lastSeen: now,
		}
		rl.clients[clientID] = client
	}

	client.lastSeen = now

	// 清理过期请求
	rl.cleanExpiredRequests(client, now)

	// 检查是否超过限制
	if len(client.requests) >= rl.maxReqs {
		return false
	}

	// 记录请求
	client.requests = append(client.requests, now)
	return true
}

// cleanExpiredRequests 清理过期请求
func (rl *RateLimiter) cleanExpiredRequests(client *clientInfo, now time.Time) {
	cutoff := now.Add(-rl.window)

	// 找到第一个未过期的请求
	validIndex := 0
	for i, requestTime := range client.requests {
		if requestTime.After(cutoff) {
			validIndex = i
			break
		}
		if i == len(client.requests)-1 {
			validIndex = len(client.requests)
		}
	}

	// 只保留未过期的请求
	if validIndex > 0 {
		copy(client.requests, client.requests[validIndex:])
		client.requests = client.requests[:len(client.requests)-validIndex]
	}
}

// cleanupClients 定期清理不活跃的客户端
func (rl *RateLimiter) cleanupClients() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			cutoff := now.Add(-rl.cleanup)

			for clientID, client := range rl.clients {
				if client.lastSeen.Before(cutoff) {
					delete(rl.clients, clientID)
				}
			}
			rl.mu.Unlock()

		case <-rl.stopChan:
			return
		}
	}
}

// Stop 停止限流器
func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

// GetStats 获取统计信息
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"total_clients":    len(rl.clients),
		"max_requests":     rl.maxReqs,
		"window_duration":  rl.window.String(),
		"cleanup_interval": rl.cleanup.String(),
	}
}

// PathBasedRateLimitMiddleware 基于路径的限流中间件
func PathBasedRateLimitMiddleware(pathLimits map[string]int, window time.Duration) gin.HandlerFunc {
	limiters := make(map[string]*RateLimiter)

	for path, limit := range pathLimits {
		limiters[path] = NewRateLimiter(limit, window)
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		limiter, exists := limiters[path]

		if exists {
			clientIP := c.ClientIP()
			if !limiter.Allow(clientIP) {
				err := errors.New(errors.CodeTooManyRequests, "请求过于频繁，请稍后再试")
				v1.HandleError(c, err)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// UserBasedRateLimitMiddleware 基于用户的限流中间件
func UserBasedRateLimitMiddleware(requestsPerMinute int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerMinute, window)

	return func(c *gin.Context) {
		// 尝试获取用户ID，如果没有则使用IP
		userID := getUserID(c)
		if userID == "" {
			userID = c.ClientIP()
		}

		if !limiter.Allow(userID) {
			err := errors.New(errors.CodeTooManyRequests, "请求过于频繁，请稍后再试")
			v1.HandleError(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}

// getUserID 从上下文获取用户ID
func getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			return uid
		}
		if uid, ok := userID.(uint); ok {
			return string(rune(uid))
		}
	}
	return ""
}
