package service

import (
	"context"
	"fmt"
	"fun-admin/internal/repository"
	"fun-admin/pkg/admin"
	"fun-admin/pkg/cache"
	"strings"
)

// ResourceService 资源服务层
type ResourceService struct {
	resourceRepository *repository.ResourceRepository
	resourceManager    *admin.ResourceManager
	exportService      *ExportService
	cacheManager       cache.CacheManager
}

// NewResourceService 创建资源服务层
func NewResourceService(
	resourceRepository *repository.ResourceRepository,
	resourceManager *admin.ResourceManager,
	cacheManager cache.CacheManager,
) *ResourceService {
	return &ResourceService{
		resourceRepository: resourceRepository,
		resourceManager:    resourceManager,
		exportService:      NewExportService(),
		cacheManager:       cacheManager,
	}
}

// Create 创建资源记录
func (s *ResourceService) Create(ctx context.Context, resourceSlug string, data map[string]interface{}) (map[string]interface{}, error) {
	resource := s.resourceManager.GetResourceBySlug(resourceSlug)
	if resource == nil {
		return nil, &ResourceNotFoundError{ResourceSlug: resourceSlug}
	}
	if auth, ok := resource.(admin.Authorizable); ok {
		if err := auth.CanCreate(ctx, data); err != nil {
			return nil, err
		}
	}
	if hook, ok := resource.(admin.CreateHook); ok {
		if err := hook.BeforeCreate(ctx, data); err != nil {
			return nil, err
		}
	}
	errors := admin.ValidateResourceData(resource, data)
	if len(errors) > 0 {
		return nil, &ValidationError{Errors: errors}
	}
	if err := s.resourceRepository.Create(ctx, resourceSlug, data); err != nil {
		return nil, err
	}
	if hook, ok := resource.(admin.CreateHook); ok {
		if err := hook.AfterCreate(ctx, data); err != nil {
			return nil, err
		}
	}
	s.clearResourceCache(ctx, resourceSlug)
	return data, nil
}

// Update 更新资源记录
func (s *ResourceService) Update(ctx context.Context, resourceSlug string, id interface{}, data map[string]interface{}) error {
	resource := s.resourceManager.GetResourceBySlug(resourceSlug)
	if resource == nil {
		return &ResourceNotFoundError{ResourceSlug: resourceSlug}
	}
	if auth, ok := resource.(admin.Authorizable); ok {
		if err := auth.CanUpdate(ctx, id, data); err != nil {
			return err
		}
	}
	if hook, ok := resource.(admin.UpdateHook); ok {
		if err := hook.BeforeUpdate(ctx, id, data); err != nil {
			return err
		}
	}
	errors := admin.ValidateResourceData(resource, data)
	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}
	if err := s.resourceRepository.Update(ctx, resourceSlug, id, data); err != nil {
		return err
	}
	if hook, ok := resource.(admin.UpdateHook); ok {
		if err := hook.AfterUpdate(ctx, id, data); err != nil {
			return err
		}
	}
	s.clearResourceCache(ctx, resourceSlug)
	cacheKey := s.getRecordCacheKey(resourceSlug, id)
	s.cacheManager.Delete(ctx, cacheKey)
	return nil
}

// Delete 删除资源记录
func (s *ResourceService) Delete(ctx context.Context, resourceSlug string, id interface{}) error {
	resource := s.resourceManager.GetResourceBySlug(resourceSlug)
	if resource == nil {
		return &ResourceNotFoundError{ResourceSlug: resourceSlug}
	}
	if auth, ok := resource.(admin.Authorizable); ok {
		if err := auth.CanDelete(ctx, id); err != nil {
			return err
		}
	}
	if hook, ok := resource.(admin.DeleteHook); ok {
		if err := hook.BeforeDelete(ctx, id); err != nil {
			return err
		}
	}
	if err := s.resourceRepository.Delete(ctx, resourceSlug, id); err != nil {
		return err
	}
	if hook, ok := resource.(admin.DeleteHook); ok {
		if err := hook.AfterDelete(ctx, id); err != nil {
			return err
		}
	}
	s.clearResourceCache(ctx, resourceSlug)
	cacheKey := s.getRecordCacheKey(resourceSlug, id)
	s.cacheManager.Delete(ctx, cacheKey)
	return nil
}

// DeleteBatch 批量删除资源记录
func (s *ResourceService) DeleteBatch(ctx context.Context, resourceSlug string, ids []interface{}) (int64, error) {
	// 获取资源配置
	resource := s.resourceManager.GetResourceBySlug(resourceSlug)
	if resource == nil {
		return 0, &ResourceNotFoundError{ResourceSlug: resourceSlug}
	}

	// 批量删除记录
	affected, err := s.resourceRepository.DeleteBatch(ctx, resourceSlug, ids)
	if err != nil {
		return 0, err
	}

	// 清除相关缓存
	s.clearResourceCache(ctx, resourceSlug)

	// 清除每条记录的缓存
	for _, id := range ids {
		cacheKey := s.getRecordCacheKey(resourceSlug, id)
		s.cacheManager.Delete(ctx, cacheKey)
	}

	return affected, nil
}

// Get 获取资源记录详情
func (s *ResourceService) Get(ctx context.Context, resourceSlug string, id interface{}) (map[string]interface{}, error) {
	// 检查缓存
	cacheKey := s.getRecordCacheKey(resourceSlug, id)
	if cached, err := s.cacheManager.Get(ctx, cacheKey); err == nil && cached != nil {
		if result, ok := cached.(map[string]interface{}); ok {
			return result, nil
		}
	}

	// 获取资源配置
	resource := s.resourceManager.GetResourceBySlug(resourceSlug)
	if resource == nil {
		return nil, &ResourceNotFoundError{ResourceSlug: resourceSlug}
	}

	// 获取关联字段信息
	relationships := s.getRelationships(resource)

	// 获取记录（包含关联数据）
	if len(relationships) > 0 {
		result, err := s.resourceRepository.FindByIDWithRelationships(ctx, resourceSlug, id, relationships)
		if err != nil {
			return nil, err
		}

		// 缓存结果
		s.cacheManager.Set(ctx, cacheKey, result, cache.DefaultExpiration)

		return result, nil
	}

	// 获取记录
	result, err := s.resourceRepository.FindByID(ctx, resourceSlug, id)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	s.cacheManager.Set(ctx, cacheKey, result, cache.DefaultExpiration)

	return result, nil
}

// List 获取资源记录列表
func (s *ResourceService) List(
	ctx context.Context,
	resourceSlug string,
	page, pageSize int,
	filters map[string]interface{},
	search map[string]interface{},
	orderBy string,
	orderDirection string,
) ([]map[string]interface{}, int64, error) {
	resource := s.resourceManager.GetResourceBySlug(resourceSlug)
	if resource != nil {
		if auth, ok := resource.(admin.Authorizable); ok {
			if err := auth.CanList(ctx); err != nil {
				return nil, 0, err
			}
		}
	}
	// 解析布尔值字符串
	for k, v := range filters {
		if sv, ok := v.(string); ok {
			lv := strings.ToLower(sv)
			if lv == "true" {
				filters[k] = true
			}
			if lv == "false" {
				filters[k] = false
			}
		}
	}
	// 白名单过滤：filters/search/orderBy
	filters = s.sanitizeFilters(resource, filters)
	search = s.sanitizeSearch(resource, search)

	// 排序字段白名单与默认排序
	orderBy, orderDirection = s.sanitizeOrder(resource, orderBy, orderDirection)

	cacheKey := s.getListCacheKey(resourceSlug, page, pageSize, filters, search, orderBy, orderDirection)
	if cached, err := s.cacheManager.Get(ctx, cacheKey); err == nil && cached != nil {
		if result, ok := cached.(map[string]interface{}); ok {
			if items, ok := result["items"].([]map[string]interface{}); ok {
				if total, ok := result["total"].(int64); ok {
					return items, total, nil
				}
			}
		}
	}
	results, total, err := s.resourceRepository.ListWithRelationshipsAndFilters(
		ctx, resourceSlug, page, pageSize, s.getRelationships(s.resourceManager.GetResourceBySlug(resourceSlug)), filters, search, orderBy, orderDirection)
	if err != nil {
		return nil, 0, err
	}
	cacheData := map[string]interface{}{"items": results, "total": total}
	s.cacheManager.Set(ctx, cacheKey, cacheData, cache.DefaultExpiration)
	return results, total, nil
}

// Export 导出资源数据
func (s *ResourceService) Export(
	ctx context.Context,
	resourceSlug string,
	filters map[string]interface{}, // 精确过滤条件
	search map[string]interface{}, // 模糊搜索条件
	orderBy string, // 排序字段
	orderDirection string, // 排序方向 ASC/DESC
	format string, // 导出格式 (csv, excel)
) ([]byte, string, error) {
	// 获取资源配置
	resource := s.resourceManager.GetResourceBySlug(resourceSlug)
	if resource == nil {
		return nil, "", &ResourceNotFoundError{ResourceSlug: resourceSlug}
	}

	// 检查资源是否支持导出
	if exportable, ok := resource.(admin.Exportable); ok {
		if !exportable.IsExportable() {
			return nil, "", fmt.Errorf("资源不支持导出功能")
		}
	}

	// 白名单过滤与默认排序
	filters = s.sanitizeFilters(resource, filters)
	search = s.sanitizeSearch(resource, search)
	orderBy, orderDirection = s.sanitizeOrder(resource, orderBy, orderDirection)

	// 获取所有数据（不分页）
	results, _, err := s.resourceRepository.ListWithRelationshipsAndFilters(
		ctx, resourceSlug, 1, 10000, s.getRelationships(resource), filters, search, orderBy, orderDirection)
	if err != nil {
		return nil, "", err
	}

	// 获取字段配置
	fields := resource.GetFields()

	// 构建表头映射
	headers := make(map[string]string)
	for _, field := range fields {
		headers[field.GetName()] = field.GetLabel()
	}

	// 处理导出数据
	var exportedData []byte
	var filename string

	switch strings.ToLower(format) {
	case "excel", "xlsx":
		exportedData, err = s.exportService.ExportToExcel(results, headers)
		if err != nil {
			return nil, "", err
		}
		filename = s.exportService.GenerateFileName(resource.GetTitle(), "xlsx")

	case "csv":
		fallthrough
	default:
		exportedData, err = s.exportService.ExportToCSV(results, headers)
		if err != nil {
			return nil, "", err
		}
		filename = s.exportService.GenerateFileName(resource.GetTitle(), "csv")
	}

	return exportedData, filename, nil
}

// Restore 恢复软删除
func (s *ResourceService) Restore(ctx context.Context, resourceSlug string, id interface{}) error {
	resource := s.resourceManager.GetResourceBySlug(resourceSlug)
	if resource == nil {
		return &ResourceNotFoundError{ResourceSlug: resourceSlug}
	}
	// 可加 Authorizable 针对恢复的权限（此处复用 Update 或 Delete 权限约定）
	if auth, ok := resource.(admin.Authorizable); ok {
		if err := auth.CanUpdate(ctx, id, map[string]interface{}{}); err != nil {
			return err
		}
	}
	if err := s.resourceRepository.Restore(ctx, resourceSlug, id); err != nil {
		return err
	}
	s.clearResourceCache(ctx, resourceSlug)
	return nil
}

// ForceDelete 强制删除（硬删）
func (s *ResourceService) ForceDelete(ctx context.Context, resourceSlug string, id interface{}) error {
	resource := s.resourceManager.GetResourceBySlug(resourceSlug)
	if resource == nil {
		return &ResourceNotFoundError{ResourceSlug: resourceSlug}
	}
	if auth, ok := resource.(admin.Authorizable); ok {
		if err := auth.CanDelete(ctx, id); err != nil {
			return err
		}
	}
	if err := s.resourceRepository.ForceDelete(ctx, resourceSlug, id); err != nil {
		return err
	}
	s.clearResourceCache(ctx, resourceSlug)
	return nil
}

// getRelationships 获取资源的关联字段信息
func (s *ResourceService) getRelationships(resource admin.Resource) map[string]string {
	relationships := make(map[string]string)

	// 遍历字段，查找关联字段
	fields := resource.GetFields()
	for _, field := range fields {
		if relField, ok := field.(*admin.RelationshipField); ok {
			relationships[relField.GetName()] = relField.RelatedResource
		}
	}

	return relationships
}

// sanitizeFilters 仅保留资源声明的精确过滤字段，且转换日期范围等约定键
func (s *ResourceService) sanitizeFilters(resource admin.Resource, filters map[string]interface{}) map[string]interface{} {
	if filters == nil {
		return map[string]interface{}{}
	}
	sanitized := make(map[string]interface{})
	// 先处理日期范围快捷键，转换到 search 由上游处理已移除，这里仅清理占位
	if _, ok := filters["created_at_from"]; ok {
		// 留给上层转换，仓储使用标准 where 片段
	}
	if _, ok := filters["created_at_to"]; ok {
		// 同上
	}
	var allowed map[string]struct{}
	if f, ok := resource.(admin.Filterable); ok {
		allowed = toSet(f.GetFilterableFields())
	} else {
		allowed = toSet(s.getFieldNames(resource))
	}
	for k, v := range filters {
		if k == "created_at_from" || k == "created_at_to" || k == "trashed" {
			sanitized[k] = v
			continue
		}
		if _, ok := allowed[k]; ok {
			sanitized[k] = v
		}
	}
	return sanitized
}

// sanitizeSearch 仅保留资源声明的模糊搜索字段，并支持 created_at 范围到 search 的迁移
func (s *ResourceService) sanitizeSearch(resource admin.Resource, search map[string]interface{}) map[string]interface{} {
	if search == nil {
		search = make(map[string]interface{})
	}
	var allowed map[string]struct{}
	if se, ok := resource.(admin.Searchable); ok {
		allowed = toSet(se.GetSearchableFields())
	} else {
		allowed = toSet(s.getFieldNames(resource))
	}
	sanitized := make(map[string]interface{})
	for k, v := range search {
		if _, ok := allowed[k]; ok {
			sanitized[k] = v
		}
	}
	return sanitized
}

// sanitizeOrder 校验排序字段并回落到默认排序
func (s *ResourceService) sanitizeOrder(resource admin.Resource, orderBy, orderDirection string) (string, string) {
	normalizeDir := func(dir string) string {
		if strings.ToUpper(dir) == "ASC" {
			return "ASC"
		}
		return "DESC"
	}
	// 来自资源的允许排序字段，若未实现则禁用客户端排序
	if so, ok := resource.(admin.Sortable); ok {
		allowed := toSet(so.GetSortableFields())
		if _, ok := allowed[orderBy]; !ok {
			orderBy = ""
		}
	} else {
		orderBy = ""
	}
	if orderBy == "" {
		if def, ok := resource.(admin.DefaultOrder); ok {
			f, d := def.GetDefaultOrder()
			orderBy, orderDirection = f, normalizeDir(d)
		}
	} else {
		orderDirection = normalizeDir(orderDirection)
	}
	return orderBy, orderDirection
}

func toSet(list []string) map[string]struct{} {
	s := make(map[string]struct{}, len(list))
	for _, v := range list {
		s[v] = struct{}{}
	}
	return s
}

// getFieldNames 提取资源字段名集合
func (s *ResourceService) getFieldNames(resource admin.Resource) []string {
	fields := resource.GetFields()
	names := make([]string, 0, len(fields))
	for _, f := range fields {
		names = append(names, f.GetName())
	}
	return names
}

// clearResourceCache 清除资源相关缓存
func (s *ResourceService) clearResourceCache(ctx context.Context, resourceSlug string) {
	// 这里可以实现更复杂的缓存清除逻辑
	// 例如：清除所有与该资源相关的列表缓存
	// 简化处理：清除所有缓存
	s.cacheManager.Flush(ctx)
}

// getRecordCacheKey 生成记录缓存键
func (s *ResourceService) getRecordCacheKey(resourceSlug string, id interface{}) string {
	return "resource:" + resourceSlug + ":record:" + s.interfaceToString(id)
}

// getListCacheKey 生成列表缓存键
func (s *ResourceService) getListCacheKey(
	resourceSlug string,
	page, pageSize int,
	filters map[string]interface{},
	search map[string]interface{},
	orderBy string,
	orderDirection string,
) string {
	key := "resource:" + resourceSlug + ":list:" +
		"page-" + s.intToString(page) +
		":size-" + s.intToString(pageSize)

	if orderBy != "" {
		key += ":order-" + orderBy + "-" + orderDirection
	}

	// 添加过滤条件到键中
	for k, v := range filters {
		key += ":filter-" + k + "-" + s.interfaceToString(v)
	}

	// 添加搜索条件到键中
	for k, v := range search {
		key += ":search-" + k + "-" + s.interfaceToString(v)
	}

	return key
}

// interfaceToString 将interface{}转换为字符串
func (s *ResourceService) interfaceToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return s.intToString(val)
	case int64:
		return s.int64ToString(val)
	default:
		return ""
	}
}

// intToString 将int转换为字符串
func (s *ResourceService) intToString(v int) string {
	return fmt.Sprintf("%d", v)
}

// int64ToString 将int64转换为字符串
func (s *ResourceService) int64ToString(v int64) string {
	return fmt.Sprintf("%d", v)
}

// ResourceNotFoundError 资源不存在错误
type ResourceNotFoundError struct {
	ResourceSlug string
}

func (e *ResourceNotFoundError) Error() string {
	return "resource not found: " + e.ResourceSlug
}

// ValidationError 验证错误
type ValidationError struct {
	Errors map[string][]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
