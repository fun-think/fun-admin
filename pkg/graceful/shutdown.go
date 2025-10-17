package graceful

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// ShutdownManager 优雅关闭管理器
type ShutdownManager struct {
	logger    *zap.Logger
	timeout   time.Duration
	callbacks []ShutdownCallback
	mu        sync.RWMutex
}

// ShutdownCallback 关闭回调函数
type ShutdownCallback struct {
	Name     string
	Priority int // 优先级，数字越小优先级越高
	Fn       func(ctx context.Context) error
}

// NewShutdownManager 创建关闭管理器
func NewShutdownManager(logger *zap.Logger, timeout time.Duration) *ShutdownManager {
	return &ShutdownManager{
		logger:    logger,
		timeout:   timeout,
		callbacks: make([]ShutdownCallback, 0),
	}
}

// Register 注册关闭回调
func (sm *ShutdownManager) Register(name string, priority int, fn func(ctx context.Context) error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	callback := ShutdownCallback{
		Name:     name,
		Priority: priority,
		Fn:       fn,
	}

	// 按优先级插入
	inserted := false
	for i, existing := range sm.callbacks {
		if priority < existing.Priority {
			// 在当前位置插入
			sm.callbacks = append(sm.callbacks[:i], append([]ShutdownCallback{callback}, sm.callbacks[i:]...)...)
			inserted = true
			break
		}
	}

	if !inserted {
		sm.callbacks = append(sm.callbacks, callback)
	}
}

// Wait 等待关闭信号并执行关闭流程
func (sm *ShutdownManager) Wait() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	sm.logger.Info("Received shutdown signal, starting graceful shutdown...")

	sm.Shutdown()
}

// Shutdown 执行关闭流程
func (sm *ShutdownManager) Shutdown() {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	sm.logger.Info("Starting graceful shutdown process",
		zap.Duration("timeout", sm.timeout),
		zap.Int("callbacks", len(sm.callbacks)),
	)

	// 执行所有关闭回调
	for _, callback := range sm.callbacks {
		sm.executeCallback(ctx, callback)
	}

	sm.logger.Info("Graceful shutdown completed")
}

// executeCallback 执行单个关闭回调
func (sm *ShutdownManager) executeCallback(ctx context.Context, callback ShutdownCallback) {
	start := time.Now()

	sm.logger.Info("Executing shutdown callback",
		zap.String("name", callback.Name),
		zap.Int("priority", callback.Priority),
	)

	// 为每个回调创建独立的超时上下文
	callbackCtx, cancel := context.WithTimeout(ctx, sm.timeout/2)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- callback.Fn(callbackCtx)
	}()

	select {
	case err := <-done:
		duration := time.Since(start)
		if err != nil {
			sm.logger.Error("Shutdown callback failed",
				zap.String("name", callback.Name),
				zap.Error(err),
				zap.Duration("duration", duration),
			)
		} else {
			sm.logger.Info("Shutdown callback completed",
				zap.String("name", callback.Name),
				zap.Duration("duration", duration),
			)
		}
	case <-callbackCtx.Done():
		sm.logger.Error("Shutdown callback timeout",
			zap.String("name", callback.Name),
			zap.Duration("timeout", sm.timeout/2),
		)
	}
}

// Resource 资源接口
type Resource interface {
	Close() error
	Name() string
}

// ResourceManager 资源管理器
type ResourceManager struct {
	resources []Resource
	mu        sync.RWMutex
	logger    *zap.Logger
}

// NewResourceManager 创建资源管理器
func NewResourceManager(logger *zap.Logger) *ResourceManager {
	return &ResourceManager{
		resources: make([]Resource, 0),
		logger:    logger,
	}
}

// Register 注册资源
func (rm *ResourceManager) Register(resource Resource) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.resources = append(rm.resources, resource)
	rm.logger.Info("Resource registered", zap.String("name", resource.Name()))
}

// CloseAll 关闭所有资源
func (rm *ResourceManager) CloseAll(ctx context.Context) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	var errors []error

	// 反向关闭资源（后进先出）
	for i := len(rm.resources) - 1; i >= 0; i-- {
		resource := rm.resources[i]

		rm.logger.Info("Closing resource", zap.String("name", resource.Name()))

		done := make(chan error, 1)
		go func(r Resource) {
			done <- r.Close()
		}(resource)

		select {
		case err := <-done:
			if err != nil {
				rm.logger.Error("Failed to close resource",
					zap.String("name", resource.Name()),
					zap.Error(err),
				)
				errors = append(errors, err)
			} else {
				rm.logger.Info("Resource closed successfully", zap.String("name", resource.Name()))
			}
		case <-ctx.Done():
			rm.logger.Error("Resource close timeout", zap.String("name", resource.Name()))
			errors = append(errors, fmt.Errorf("timeout closing resource: %s", resource.Name()))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to close %d resources", len(errors))
	}

	return nil
}

// HTTPServerResource HTTP服务器资源包装器
type HTTPServerResource struct {
	server interface {
		Shutdown(ctx context.Context) error
	}
	name string
}

// NewHTTPServerResource 创建HTTP服务器资源
func NewHTTPServerResource(server interface {
	Shutdown(ctx context.Context) error
}, name string) *HTTPServerResource {
	return &HTTPServerResource{
		server: server,
		name:   name,
	}
}

func (h *HTTPServerResource) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return h.server.Shutdown(ctx)
}

func (h *HTTPServerResource) Name() string {
	return h.name
}

// DatabaseResource 数据库资源包装器
type DatabaseResource struct {
	db interface {
		Close() error
	}
	name string
}

// NewDatabaseResource 创建数据库资源
func NewDatabaseResource(db interface{ Close() error }, name string) *DatabaseResource {
	return &DatabaseResource{
		db:   db,
		name: name,
	}
}

func (d *DatabaseResource) Close() error {
	return d.db.Close()
}

func (d *DatabaseResource) Name() string {
	return d.name
}

// CacheResource 缓存资源包装器
type CacheResource struct {
	cache interface {
		Close() error
	}
	name string
}

// NewCacheResource 创建缓存资源
func NewCacheResource(cache interface{ Close() error }, name string) *CacheResource {
	return &CacheResource{
		cache: cache,
		name:  name,
	}
}

func (c *CacheResource) Close() error {
	return c.cache.Close()
}

func (c *CacheResource) Name() string {
	return c.name
}

// WorkerResource 工作协程资源包装器
type WorkerResource struct {
	stopChan chan struct{}
	doneChan chan struct{}
	name     string
}

// NewWorkerResource 创建工作协程资源
func NewWorkerResource(stopChan, doneChan chan struct{}, name string) *WorkerResource {
	return &WorkerResource{
		stopChan: stopChan,
		doneChan: doneChan,
		name:     name,
	}
}

func (w *WorkerResource) Close() error {
	close(w.stopChan)

	// 等待工作协程完成
	select {
	case <-w.doneChan:
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("worker %s shutdown timeout", w.name)
	}
}

func (w *WorkerResource) Name() string {
	return w.name
}
