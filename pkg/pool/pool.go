package pool

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrPoolClosed = errors.New("pool is closed")
	ErrPoolFull   = errors.New("pool is full")
)

// Pool 通用资源池接口
type Pool interface {
	Get(ctx context.Context) (interface{}, error)
	Put(resource interface{}) error
	Close() error
	Len() int
	Cap() int
}

// ResourceFactory 资源工厂函数
type ResourceFactory func() (interface{}, error)

// ResourceValidator 资源验证函数
type ResourceValidator func(resource interface{}) bool

// Config 池配置
type Config struct {
	MaxSize       int               // 最大资源数
	MinSize       int               // 最小资源数
	MaxIdleTime   time.Duration     // 最大空闲时间
	ValidateOnGet bool              // 获取时验证
	ValidateOnPut bool              // 放回时验证
	Factory       ResourceFactory   // 资源工厂
	Validator     ResourceValidator // 资源验证器
	Destructor    func(interface{}) // 资源销毁器
}

// SimplePool 简单资源池实现
type SimplePool struct {
	config    Config
	resources chan interface{}
	mu        sync.RWMutex
	closed    bool
	created   int
	lastUsed  map[interface{}]time.Time
}

// NewSimplePool 创建简单资源池
func NewSimplePool(config Config) (*SimplePool, error) {
	if config.MaxSize <= 0 {
		return nil, errors.New("max size must be greater than 0")
	}

	if config.MinSize < 0 {
		config.MinSize = 0
	}

	if config.MinSize > config.MaxSize {
		config.MinSize = config.MaxSize
	}

	if config.MaxIdleTime <= 0 {
		config.MaxIdleTime = 30 * time.Minute
	}

	pool := &SimplePool{
		config:    config,
		resources: make(chan interface{}, config.MaxSize),
		lastUsed:  make(map[interface{}]time.Time),
	}

	// 预创建最小数量的资源
	for i := 0; i < config.MinSize; i++ {
		resource, err := pool.createResource()
		if err != nil {
			return nil, err
		}
		pool.resources <- resource
	}

	// 启动清理协程
	go pool.cleanup()

	return pool, nil
}

// Get 获取资源
func (p *SimplePool) Get(ctx context.Context) (interface{}, error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil, ErrPoolClosed
	}
	p.mu.RUnlock()

	select {
	case resource := <-p.resources:
		// 验证资源
		if p.config.ValidateOnGet && p.config.Validator != nil {
			if !p.config.Validator(resource) {
				// 资源无效，销毁并创建新的
				p.destroyResource(resource)
				return p.createAndReturn()
			}
		}

		p.mu.Lock()
		p.lastUsed[resource] = time.Now()
		p.mu.Unlock()

		return resource, nil

	case <-ctx.Done():
		return nil, ctx.Err()

	default:
		// 池中没有资源，尝试创建新的
		return p.createAndReturn()
	}
}

// Put 放回资源
func (p *SimplePool) Put(resource interface{}) error {
	if resource == nil {
		return nil
	}

	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		p.destroyResource(resource)
		return ErrPoolClosed
	}
	p.mu.RUnlock()

	// 验证资源
	if p.config.ValidateOnPut && p.config.Validator != nil {
		if !p.config.Validator(resource) {
			p.destroyResource(resource)
			return nil
		}
	}

	select {
	case p.resources <- resource:
		p.mu.Lock()
		p.lastUsed[resource] = time.Now()
		p.mu.Unlock()
		return nil
	default:
		// 池已满，销毁资源
		p.destroyResource(resource)
		return nil
	}
}

// Close 关闭资源池
func (p *SimplePool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true

	// 销毁所有资源
	close(p.resources)
	for resource := range p.resources {
		p.destroyResource(resource)
	}

	return nil
}

// Len 当前资源数量
func (p *SimplePool) Len() int {
	return len(p.resources)
}

// Cap 池容量
func (p *SimplePool) Cap() int {
	return p.config.MaxSize
}

// createAndReturn 创建并返回新资源
func (p *SimplePool) createAndReturn() (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil, ErrPoolClosed
	}

	if p.created >= p.config.MaxSize {
		return nil, ErrPoolFull
	}

	resource, err := p.createResource()
	if err != nil {
		return nil, err
	}

	p.lastUsed[resource] = time.Now()
	return resource, nil
}

// createResource 创建资源
func (p *SimplePool) createResource() (interface{}, error) {
	if p.config.Factory == nil {
		return nil, errors.New("resource factory is required")
	}

	resource, err := p.config.Factory()
	if err != nil {
		return nil, err
	}

	p.created++
	return resource, nil
}

// destroyResource 销毁资源
func (p *SimplePool) destroyResource(resource interface{}) {
	p.mu.Lock()
	delete(p.lastUsed, resource)
	p.created--
	p.mu.Unlock()

	if p.config.Destructor != nil {
		p.config.Destructor(resource)
	}
}

// cleanup 清理过期资源
func (p *SimplePool) cleanup() {
	ticker := time.NewTicker(p.config.MaxIdleTime / 2)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.RLock()
		if p.closed {
			p.mu.RUnlock()
			return
		}
		p.mu.RUnlock()

		p.cleanupExpired()
	}
}

// cleanupExpired 清理过期资源
func (p *SimplePool) cleanupExpired() {
	now := time.Now()
	cutoff := now.Add(-p.config.MaxIdleTime)

	// 收集需要清理的资源
	var toCleanup []interface{}

	for {
		select {
		case resource := <-p.resources:
			p.mu.RLock()
			lastUsed := p.lastUsed[resource]
			minSize := p.config.MinSize
			currentSize := len(p.resources) + 1
			p.mu.RUnlock()

			// 检查是否需要清理
			if lastUsed.Before(cutoff) && currentSize > minSize {
				toCleanup = append(toCleanup, resource)
			} else {
				// 放回资源
				p.resources <- resource
				return // 退出清理循环
			}
		default:
			// 没有更多资源需要检查
			goto cleanup
		}
	}

cleanup:
	// 清理收集到的资源
	for _, resource := range toCleanup {
		p.destroyResource(resource)
	}
}

// WorkerPool 工作协程池
type WorkerPool struct {
	maxWorkers  int
	taskQueue   chan func()
	workerQueue chan chan func()
	quit        chan struct{}
	wg          sync.WaitGroup
	mu          sync.RWMutex
	started     bool
}

// NewWorkerPool 创建工作协程池
func NewWorkerPool(maxWorkers int, queueSize int) *WorkerPool {
	return &WorkerPool{
		maxWorkers:  maxWorkers,
		taskQueue:   make(chan func(), queueSize),
		workerQueue: make(chan chan func(), maxWorkers),
		quit:        make(chan struct{}),
	}
}

// Start 启动工作池
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.started {
		return
	}

	wp.started = true

	// 启动调度器
	go wp.dispatcher()

	// 启动工作协程
	for i := 0; i < wp.maxWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Submit 提交任务
func (wp *WorkerPool) Submit(task func()) error {
	wp.mu.RLock()
	started := wp.started
	wp.mu.RUnlock()

	if !started {
		return errors.New("worker pool not started")
	}

	select {
	case wp.taskQueue <- task:
		return nil
	default:
		return errors.New("task queue is full")
	}
}

// Stop 停止工作池
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.started {
		return
	}

	close(wp.quit)
	wp.wg.Wait()
	wp.started = false
}

// dispatcher 任务调度器
func (wp *WorkerPool) dispatcher() {
	for {
		select {
		case task := <-wp.taskQueue:
			// 获取可用的工作协程
			select {
			case workerChannel := <-wp.workerQueue:
				// 分配任务给工作协程
				workerChannel <- task
			case <-wp.quit:
				return
			}
		case <-wp.quit:
			return
		}
	}
}

// worker 工作协程
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	// 创建工作协程的任务通道
	taskChannel := make(chan func())

	for {
		// 将工作协程注册到工作队列
		select {
		case wp.workerQueue <- taskChannel:
		case <-wp.quit:
			return
		}

		// 等待任务或退出信号
		select {
		case task := <-taskChannel:
			// 执行任务
			func() {
				defer func() {
					if r := recover(); r != nil {
						// 任务panic处理
					}
				}()
				task()
			}()
		case <-wp.quit:
			return
		}
	}
}
