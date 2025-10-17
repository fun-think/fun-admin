package app

import (
	"context"
	"fmt"
	"fun-admin/pkg/container"
	"fun-admin/pkg/graceful"
	"fun-admin/pkg/server"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type ServiceProvider interface {
	Register(c *container.Container)
	Boot(c *container.Container) error
}

type App struct {
	name            string
	servers         []server.Server
	shutdownManager *graceful.ShutdownManager
	resourceManager *graceful.ResourceManager
	logger          *zap.Logger
	stopTimeout     time.Duration
}

type Option func(a *App)

func NewApp(opts ...Option) *App {
	// 创建默认logger
	logger, _ := zap.NewProduction()

	a := &App{
		logger:      logger,
		stopTimeout: 30 * time.Second,
	}

	// 初始化管理器
	a.shutdownManager = graceful.NewShutdownManager(logger, a.stopTimeout)
	a.resourceManager = graceful.NewResourceManager(logger)

	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithServer(servers ...server.Server) Option {
	return func(a *App) {
		a.servers = servers
	}
}

func WithName(name string) Option {
	return func(a *App) {
		a.name = name
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(a *App) {
		a.logger = logger
		// 重新创建管理器
		a.shutdownManager = graceful.NewShutdownManager(logger, a.stopTimeout)
		a.resourceManager = graceful.NewResourceManager(logger)
	}
}

func WithStopTimeout(timeout time.Duration) Option {
	return func(a *App) {
		a.stopTimeout = timeout
		if a.logger != nil {
			a.shutdownManager = graceful.NewShutdownManager(a.logger, timeout)
		}
	}
}

func (a *App) Run(ctx context.Context) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	// 注册资源管理器到关闭流程
	a.shutdownManager.Register("resources", 100, func(ctx context.Context) error {
		return a.resourceManager.CloseAll(ctx)
	})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// 启动服务器
	for i, srv := range a.servers {
		wg.Add(1)

		// 注册服务器关闭回调
		serverName := fmt.Sprintf("server-%d", i)
		a.shutdownManager.Register(serverName, 50, func(shutdownCtx context.Context) error {
			return srv.Stop(shutdownCtx)
		})

		go func(srv server.Server, name string) {
			defer wg.Done()

			a.logger.Info("Starting server", zap.String("name", name))
			err := srv.Start(ctx)
			if err != nil {
				a.logger.Error("Server start error",
					zap.String("name", name),
					zap.Error(err),
				)
			}
		}(srv, serverName)
	}

	// 等待信号或上下文取消
	select {
	case sig := <-signals:
		a.logger.Info("Received termination signal", zap.String("signal", sig.String()))
	case <-ctx.Done():
		a.logger.Info("Context canceled")
	}

	// 取消上下文，通知所有服务器停止
	cancel()

	// 执行优雅关闭
	a.logger.Info("Starting graceful shutdown")
	a.shutdownManager.Shutdown()

	// 等待所有服务器停止
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// 等待服务器停止或超时
	select {
	case <-done:
		a.logger.Info("All servers stopped gracefully")
	case <-time.After(a.stopTimeout):
		a.logger.Warn("Shutdown timeout, forcing exit")
	}

	return nil
}

// RegisterResource 注册需要管理的资源
func (a *App) RegisterResource(resource graceful.Resource) {
	a.resourceManager.Register(resource)
}

// RegisterShutdownCallback 注册关闭回调
func (a *App) RegisterShutdownCallback(name string, priority int, fn func(ctx context.Context) error) {
	a.shutdownManager.Register(name, priority, fn)
}
