package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container 依赖注入容器
type Container struct {
	mu        sync.RWMutex
	services  map[string]interface{}
	bindings  map[string]interface{}
	instances map[string]interface{}
	once      map[string]*sync.Once
}

// New 创建新的容器
func New() *Container {
	return &Container{
		services:  make(map[string]interface{}),
		bindings:  make(map[string]interface{}),
		instances: make(map[string]interface{}),
		once:      make(map[string]*sync.Once),
	}
}

// Singleton 单例绑定
func (c *Container) Singleton(name string, factory interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	factoryType := reflect.TypeOf(factory)
	if factoryType.Kind() != reflect.Func {
		panic(fmt.Sprintf("factory must be a function, got %T", factory))
	}

	c.bindings[name] = factory
	c.once[name] = &sync.Once{}
}

// Bind 绑定服务
func (c *Container) Bind(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// Get 获取服务
func (c *Container) Get(name string) (interface{}, error) {
	// 先尝试从已创建的服务中获取
	c.mu.RLock()
	if service, exists := c.services[name]; exists {
		c.mu.RUnlock()
		return service, nil
	}
	
	// 尝试从单例实例中获取
	if service, exists := c.instances[name]; exists {
		c.mu.RUnlock()
		return service, nil
	}
	c.mu.RUnlock()

	// 查找工厂函数
	c.mu.RLock()
	factory, exists := c.bindings[name]
	once, onceExists := c.once[name]
	c.mu.RUnlock()

	if !exists || !onceExists {
		return nil, fmt.Errorf("service '%s' not found", name)
	}

	// 使用 sync.Once 确保单例只创建一次
	var instance interface{}
	var err error
	
	once.Do(func() {
		// 获取工厂函数的反射值
		factoryValue := reflect.ValueOf(factory)
		factoryType := factoryValue.Type()

		// 准备调用参数
		var args []reflect.Value
		if factoryType.NumIn() > 0 {
			// 工厂函数需要参数，传入容器本身
			args = append(args, reflect.ValueOf(c))
		}

		// 调用工厂函数创建实例
		results := factoryValue.Call(args)
		if len(results) == 0 {
			err = fmt.Errorf("factory function must return at least one value")
			return
		}
		
		instance = results[0].Interface()
		
		// 将新创建的服务存储到实例映射中
		c.mu.Lock()
		c.instances[name] = instance
		c.mu.Unlock()
	})
	
	if err != nil {
		return nil, err
	}
	
	// 如果 instance 仍未被设置，从映射中再次获取
	if instance == nil {
		c.mu.RLock()
		instance = c.instances[name]
		c.mu.RUnlock()
	}

	return instance, nil
}

// MustGet 获取服务，如果不存在则panic
func (c *Container) MustGet(name string) interface{} {
	service, err := c.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}

// Has 检查服务是否存在
func (c *Container) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.services[name]
	if exists {
		return true
	}

	_, exists = c.instances[name]
	if exists {
		return true
	}

	_, exists = c.bindings[name]
	return exists
}

// Clear 清空所有服务
func (c *Container) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.services = make(map[string]interface{})
	c.bindings = make(map[string]interface{})
	c.instances = make(map[string]interface{})
	c.once = make(map[string]*sync.Once)
}