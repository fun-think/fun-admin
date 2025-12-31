package service

import (
	"context"
	"errors"
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
	permittedData, err := s.enforceWritableFields(ctx, resource, data)
	if err != nil {
		return nil, err
	}
	data = permittedData
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
	permittedData, err := s.enforceWritableFields(ctx, resource, data)
	if err != nil {
		return err
	}
	data = permittedData
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

	resource := s.resourceManager.GetResourceBySlug(resourceSlug)

	if resource == nil {

		return nil, &ResourceNotFoundError{ResourceSlug: resourceSlug}

	}

	cacheKey := s.getRecordCacheKey(resourceSlug, id)

	if cached, err := s.cacheManager.Get(ctx, cacheKey); err == nil && cached != nil {

		if result, ok := cached.(map[string]interface{}); ok {

			return s.filterReadableRecord(ctx, resource, result), nil

		}

	}

	relationships := s.getRelationships(resource)

	var (
		result map[string]interface{}

		err error
	)

	if len(relationships) > 0 {

		result, err = s.resourceRepository.FindByIDWithRelationships(ctx, resourceSlug, id, relationships)

	} else {

		result, err = s.resourceRepository.FindByID(ctx, resourceSlug, id)

	}

	if err != nil {

		return nil, err

	}

	s.cacheManager.Set(ctx, cacheKey, result, cache.DefaultExpiration)

	return s.filterReadableRecord(ctx, resource, result), nil

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

	var relationships map[string]string
	if resource != nil {
		relationships = s.getRelationships(resource)
	}

	cacheKey := s.getListCacheKey(resourceSlug, page, pageSize, filters, search, orderBy, orderDirection)
	if cached, err := s.cacheManager.Get(ctx, cacheKey); err == nil && cached != nil {
		if result, ok := cached.(map[string]interface{}); ok {
			if items, ok := result["items"].([]map[string]interface{}); ok {
				if total, ok := result["total"].(int64); ok {
					return s.filterReadableList(ctx, resource, items), total, nil
				}
			}
		}
	}
	results, total, err := s.resourceRepository.ListWithRelationshipsAndFilters(
		ctx, resourceSlug, page, pageSize, relationships, filters, search, orderBy, orderDirection)
	if err != nil {
		return nil, 0, err
	}
	cacheData := map[string]interface{}{"items": results, "total": total}
	s.cacheManager.Set(ctx, cacheKey, cacheData, cache.DefaultExpiration)
	return s.filterReadableList(ctx, resource, results), total, nil
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

// RunAction 执行资源动作（供 Frontend 调用）
func (s *ResourceService) RunAction(
	ctx context.Context,
	resourceSlug string,
	actionName string,
	ids []interface{},
	params map[string]interface{},
) (interface{}, error) {
	resource := s.resourceManager.GetResourceBySlug(resourceSlug)
	if resource == nil {
		return nil, &ResourceNotFoundError{ResourceSlug: resourceSlug}
	}
	executor, ok := resource.(admin.ActionExecutor)
	if !ok {
		return s.handleBuiltInAction(ctx, resourceSlug, actionName, ids, params)
	}
	return executor.RunAction(ctx, actionName, ids, params)
}

func (s *ResourceService) handleBuiltInAction(
	ctx context.Context,
	resourceSlug string,
	actionName string,
	ids []interface{},
	params map[string]interface{},
) (interface{}, error) {
	if resourceSlug != "crud_items" {
		return nil, ErrActionNotSupported
	}
	switch actionName {
	case "reset_values":
		updated := 0
		for _, id := range ids {
			if err := s.resourceRepository.Update(ctx, resourceSlug, id, map[string]interface{}{"value": "", "remark": ""}); err == nil {
				updated++
			}
		}
		return map[string]interface{}{"updated": updated}, nil
	case "bulk_delete":
		count, err := s.resourceRepository.DeleteBatch(ctx, resourceSlug, ids)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{"deleted": count}, nil
	default:
		return nil, ErrActionNotSupported
	}
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

func (s *ResourceService) getReadableFieldSet(ctx context.Context, resource admin.Resource) map[string]struct{} {
	if resource == nil {
		return nil
	}
	var fields []string
	if provider, ok := any(resource).(admin.FieldPermissionProvider); ok {
		perms := provider.GetFieldPermissions(ctx)
		if len(perms.Readable) > 0 {
			fields = append(fields, perms.Readable...)
		}
	}
	if len(fields) == 0 {
		fields = append(fields, s.getFieldNames(resource)...)
	}
	defaults := []string{"id", "created_at", "updated_at"}
	fields = append(fields, defaults...)
	return toSet(fields)
}

func (s *ResourceService) getWritableFieldSet(ctx context.Context, resource admin.Resource) map[string]struct{} {
	if resource == nil {
		return nil
	}
	var fields []string
	if provider, ok := any(resource).(admin.FieldPermissionProvider); ok {
		perms := provider.GetFieldPermissions(ctx)
		if len(perms.Writable) > 0 {
			fields = append(fields, perms.Writable...)
		}
	}
	if len(fields) == 0 {
		for _, field := range resource.GetFields() {
			fields = append(fields, field.GetName())
		}
	}
	return toSet(fields)
}

func (s *ResourceService) getReadOnlyFieldSet(resource admin.Resource) map[string]struct{} {
	if resource == nil {
		return nil
	}
	return toSet(resource.GetReadOnlyFields())
}

func (s *ResourceService) filterReadableRecord(ctx context.Context, resource admin.Resource, record map[string]interface{}) map[string]interface{} {
	if record == nil {
		return nil
	}
	readable := s.getReadableFieldSet(ctx, resource)
	return s.keepFields(record, readable)
}

func (s *ResourceService) filterReadableList(ctx context.Context, resource admin.Resource, list []map[string]interface{}) []map[string]interface{} {
	if len(list) == 0 {
		return list
	}
	readable := s.getReadableFieldSet(ctx, resource)
	filtered := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		filtered = append(filtered, s.keepFields(item, readable))
	}
	return filtered
}

func (s *ResourceService) keepFields(record map[string]interface{}, readable map[string]struct{}) map[string]interface{} {
	if record == nil {
		return nil
	}
	if len(readable) == 0 {
		clone := make(map[string]interface{}, len(record))
		for k, v := range record {
			clone[k] = v
		}
		return clone
	}
	clone := make(map[string]interface{}, len(readable))
	for k, v := range record {
		if _, ok := readable[k]; ok || strings.HasSuffix(k, "_data") {
			clone[k] = v
		}
	}
	return clone
}

func (s *ResourceService) enforceWritableFields(ctx context.Context, resource admin.Resource, data map[string]interface{}) (map[string]interface{}, error) {
	if data == nil {
		return map[string]interface{}{}, nil
	}
	allowed := s.getWritableFieldSet(ctx, resource)
	readOnly := s.getReadOnlyFieldSet(resource)
	clean := make(map[string]interface{}, len(data))
	errs := make(map[string][]string)

	for field, value := range data {
		if readOnly != nil {
			if _, ok := readOnly[field]; ok {
				errs[field] = append(errs[field], "字段只读，禁止写入")
				continue
			}
		}
		if allowed != nil {
			if _, ok := allowed[field]; !ok {
				errs[field] = append(errs[field], "没有写入该字段的权限")
				continue
			}
		}
		clean[field] = value
	}

	if len(errs) > 0 {
		return nil, &ValidationError{Errors: errs}
	}

	return clean, nil
}

// getFieldNames 提取资源字段名集合
func (s *ResourceService) getFieldNames(resource admin.Resource) []string {
	if resource == nil {
		return nil
	}
	fields := resource.GetFields()
	names := make([]string, 0, len(fields))
	for _, f := range fields {
		names = append(names, f.GetName())
	}
	return names
}

// clearResourceCache 资源相关缓存
func (s *ResourceService) clearResourceCache(ctx context.Context, resourceSlug string) {
	prefix := "resource:" + resourceSlug + ":"
	if err := s.cacheManager.DeleteByPrefix(ctx, prefix); err != nil {
		_ = s.cacheManager.Flush(ctx)
	}
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

// ErrActionNotSupported 表示资源未实现动作
var ErrActionNotSupported = errors.New("action not supported")
