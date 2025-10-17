package middleware

import (
	"github.com/gin-gonic/gin"
)

// Chain 中间件链
type Chain struct {
	middlewares []gin.HandlerFunc
}

// NewChain 创建中间件链
func NewChain(middlewares ...gin.HandlerFunc) *Chain {
	return &Chain{
		middlewares: middlewares,
	}
}

// Append 追加中间件
func (c *Chain) Append(middlewares ...gin.HandlerFunc) *Chain {
	return &Chain{
		middlewares: append(c.middlewares, middlewares...),
	}
}

// Prepend 前置中间件
func (c *Chain) Prepend(middlewares ...gin.HandlerFunc) *Chain {
	return &Chain{
		middlewares: append(middlewares, c.middlewares...),
	}
}

// Then 执行中间件链并返回最终处理器
func (c *Chain) Then(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 创建中间件处理链
		middlewares := append(c.middlewares, handler)

		// 执行中间件链
		for _, mw := range middlewares {
			mw(ctx)
			if ctx.IsAborted() {
				return
			}
		}
	}
}

// ThenFunc 执行中间件链并返回处理器函数
func (c *Chain) ThenFunc(handlerFunc func(*gin.Context)) gin.HandlerFunc {
	return c.Then(gin.HandlerFunc(handlerFunc))
}

// GetMiddlewares 获取所有中间件
func (c *Chain) GetMiddlewares() []gin.HandlerFunc {
	return c.middlewares
}

// Length 获取中间件数量
func (c *Chain) Length() int {
	return len(c.middlewares)
}

// Clone 克隆中间件链
func (c *Chain) Clone() *Chain {
	middlewares := make([]gin.HandlerFunc, len(c.middlewares))
	copy(middlewares, c.middlewares)
	return &Chain{middlewares: middlewares}
}

// MiddlewareBuilder 中间件构建器
type MiddlewareBuilder struct {
	chains map[string]*Chain
}

// NewMiddlewareBuilder 创建中间件构建器
func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		chains: make(map[string]*Chain),
	}
}

// DefineChain 定义命名中间件链
func (mb *MiddlewareBuilder) DefineChain(name string, middlewares ...gin.HandlerFunc) {
	mb.chains[name] = NewChain(middlewares...)
}

// GetChain 获取命名中间件链
func (mb *MiddlewareBuilder) GetChain(name string) *Chain {
	if chain, exists := mb.chains[name]; exists {
		return chain.Clone()
	}
	return NewChain()
}

// CombineChains 合并多个命名中间件链
func (mb *MiddlewareBuilder) CombineChains(names ...string) *Chain {
	var allMiddlewares []gin.HandlerFunc

	for _, name := range names {
		if chain, exists := mb.chains[name]; exists {
			allMiddlewares = append(allMiddlewares, chain.middlewares...)
		}
	}

	return NewChain(allMiddlewares...)
}

// RouteMiddleware 路由级中间件配置
type RouteMiddleware struct {
	Path        string            `json:"path"`
	Methods     []string          `json:"methods"`
	Middlewares []string          `json:"middlewares"`
	Exclude     []string          `json:"exclude"`
	Priority    int               `json:"priority"`
	Conditions  map[string]string `json:"conditions"`
}

// MiddlewareManager 中间件管理器
type MiddlewareManager struct {
	builder *MiddlewareBuilder
	routes  []RouteMiddleware
}

// NewMiddlewareManager 创建中间件管理器
func NewMiddlewareManager() *MiddlewareManager {
	return &MiddlewareManager{
		builder: NewMiddlewareBuilder(),
		routes:  make([]RouteMiddleware, 0),
	}
}

// RegisterMiddleware 注册中间件到构建器
func (mm *MiddlewareManager) RegisterMiddleware(name string, middlewares ...gin.HandlerFunc) {
	mm.builder.DefineChain(name, middlewares...)
}

// AddRouteMiddleware 添加路由中间件配置
func (mm *MiddlewareManager) AddRouteMiddleware(config RouteMiddleware) {
	mm.routes = append(mm.routes, config)
}

// BuildChainForRoute 为特定路由构建中间件链
func (mm *MiddlewareManager) BuildChainForRoute(path, method string) *Chain {
	var chains []string

	// 匹配路由配置
	for _, route := range mm.routes {
		if mm.matchRoute(route, path, method) {
			// 添加中间件（排除exclude列表中的）
			for _, mw := range route.Middlewares {
				excluded := false
				for _, ex := range route.Exclude {
					if ex == mw {
						excluded = true
						break
					}
				}
				if !excluded {
					chains = append(chains, mw)
				}
			}
		}
	}

	return mm.builder.CombineChains(chains...)
}

// matchRoute 匹配路由
func (mm *MiddlewareManager) matchRoute(route RouteMiddleware, path, method string) bool {
	// 路径匹配（支持通配符）
	if !mm.matchPath(route.Path, path) {
		return false
	}

	// 方法匹配
	if len(route.Methods) > 0 {
		methodMatched := false
		for _, m := range route.Methods {
			if m == method || m == "*" {
				methodMatched = true
				break
			}
		}
		if !methodMatched {
			return false
		}
	}

	return true
}

// matchPath 匹配路径（简单通配符支持）
func (mm *MiddlewareManager) matchPath(pattern, path string) bool {
	if pattern == "*" || pattern == path {
		return true
	}

	// 支持 /api/* 这样的通配符
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(path) >= len(prefix) && path[:len(prefix)] == prefix
	}

	return false
}

// ConditionalMiddleware 条件中间件
type ConditionalMiddleware struct {
	condition  func(*gin.Context) bool
	middleware gin.HandlerFunc
}

// NewConditionalMiddleware 创建条件中间件
func NewConditionalMiddleware(condition func(*gin.Context) bool, middleware gin.HandlerFunc) *ConditionalMiddleware {
	return &ConditionalMiddleware{
		condition:  condition,
		middleware: middleware,
	}
}

// Handler 返回中间件处理器
func (cm *ConditionalMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cm.condition(c) {
			cm.middleware(c)
		} else {
			c.Next()
		}
	}
}

// GroupMiddleware 路由组中间件配置
type GroupMiddleware struct {
	Group       string            `json:"group"`
	Middlewares []string          `json:"middlewares"`
	Priority    int               `json:"priority"`
	Conditions  map[string]string `json:"conditions"`
}

// Advanced middleware utilities

// OnlyIf 条件中间件工厂
func OnlyIf(condition func(*gin.Context) bool, middleware gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if condition(c) {
			middleware(c)
		} else {
			c.Next()
		}
	}
}

// SkipIf 跳过条件中间件工厂
func SkipIf(condition func(*gin.Context) bool, middleware gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !condition(c) {
			middleware(c)
		} else {
			c.Next()
		}
	}
}

// ForMethods 仅对指定方法应用中间件
func ForMethods(methods []string, middleware gin.HandlerFunc) gin.HandlerFunc {
	methodMap := make(map[string]bool)
	for _, method := range methods {
		methodMap[method] = true
	}

	return OnlyIf(func(c *gin.Context) bool {
		return methodMap[c.Request.Method]
	}, middleware)
}

// ForPaths 仅对指定路径应用中间件
func ForPaths(paths []string, middleware gin.HandlerFunc) gin.HandlerFunc {
	pathMap := make(map[string]bool)
	for _, path := range paths {
		pathMap[path] = true
	}

	return OnlyIf(func(c *gin.Context) bool {
		return pathMap[c.Request.URL.Path]
	}, middleware)
}

// ExceptPaths 排除指定路径的中间件
func ExceptPaths(paths []string, middleware gin.HandlerFunc) gin.HandlerFunc {
	pathMap := make(map[string]bool)
	for _, path := range paths {
		pathMap[path] = true
	}

	return SkipIf(func(c *gin.Context) bool {
		return pathMap[c.Request.URL.Path]
	}, middleware)
}

// LoggingChain 日志中间件链
func LoggingChain() *Chain {
	return NewChain(
		gin.Logger(),
		gin.Recovery(),
	)
}

// AuthChain 认证中间件链
func AuthChain(authMiddleware, rbacMiddleware gin.HandlerFunc) *Chain {
	return NewChain(
		authMiddleware,
		rbacMiddleware,
	)
}

// SecurityChain 安全中间件链
func SecurityChain(corsMiddleware, securityMiddleware, rateLimitMiddleware gin.HandlerFunc) *Chain {
	return NewChain(
		corsMiddleware,
		securityMiddleware,
		rateLimitMiddleware,
	)
}
