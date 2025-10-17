package admin

import (
	"sync"
)

// ResourceManager 资源管理器
type ResourceManager struct {
	mu        sync.RWMutex
	resources []Resource
	pages     []PageInterface
}

// NewResourceManager 创建资源管理器
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make([]Resource, 0),
		pages:     make([]PageInterface, 0),
	}
}

// Register 注册资源
func (rm *ResourceManager) Register(resource Resource) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.resources = append(rm.resources, resource)
}

// RegisterPage 注册自定义页面
func (rm *ResourceManager) RegisterPage(page PageInterface) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.pages = append(rm.pages, page)
}

// GetResources 获取所有资源
func (rm *ResourceManager) GetResources() []Resource {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	// 返回资源副本以避免并发问题
	resources := make([]Resource, len(rm.resources))
	copy(resources, rm.resources)

	return resources
}

// GetPages 获取所有自定义页面
func (rm *ResourceManager) GetPages() []PageInterface {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	// 返回页面副本以避免并发问题
	pages := make([]PageInterface, len(rm.pages))
	copy(pages, rm.pages)

	return pages
}

// GetResourceBySlug 根据 slug 获取资源
func (rm *ResourceManager) GetResourceBySlug(slug string) Resource {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	for _, resource := range rm.resources {
		if resource.GetSlug() == slug {
			return resource
		}
	}

	return nil
}

// GetPageBySlug 根据 slug 获取页面
func (rm *ResourceManager) GetPageBySlug(slug string) PageInterface {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	for _, page := range rm.pages {
		if page.GetSlug() == slug {
			return page
		}
	}

	return nil
}

// GetAllResourcesAndPages 获取所有资源和页面
func (rm *ResourceManager) GetAllResourcesAndPages() ([]Resource, []PageInterface) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	// 返回资源和页面副本以避免并发问题
	resources := make([]Resource, len(rm.resources))
	copy(resources, rm.resources)

	pages := make([]PageInterface, len(rm.pages))
	copy(pages, rm.pages)

	return resources, pages
}

// GlobalResourceManager 全局资源管理器实例
var GlobalResourceManager = NewResourceManager()

// Register 全局注册资源
func Register(resource Resource) {
	GlobalResourceManager.Register(resource)
}

// RegisterPage 全局注册自定义页面
func RegisterPage(page PageInterface) {
	GlobalResourceManager.RegisterPage(page)
}

// BaseResource 可复用空实现，供资源内嵌
// 资源可选择性嵌入以减少样板代码
// 注意：方法返回零值，业务层应覆盖
// 若未覆盖且被调用，通常会导致空行为或校验失败
// 使用方需实现自己的方法
type BaseResource struct{}

func (r *BaseResource) GetTitle() string            { return "" }
func (r *BaseResource) GetSlug() string             { return "" }
func (r *BaseResource) GetModel() interface{}       { return nil }
func (r *BaseResource) GetFields() []Field          { return []Field{} }
func (r *BaseResource) GetActions() []Action        { return []Action{} }
func (r *BaseResource) GetReadOnlyFields() []string { return []string{} }
func (r *BaseResource) GetColumns() []*Column       { return []*Column{} }
func (r *BaseResource) GetFilters() []*Filter       { return []*Filter{} }
