package repository

import (
	"context"
	"strings"
	"time"
)

// ResourceRepository 资源数据访问层
type ResourceRepository struct {
	Repository
}

// NewResourceRepository 创建资源数据访问层
func NewResourceRepository(repo Repository) *ResourceRepository {
	return &ResourceRepository{
		Repository: repo,
	}
}

// Create 创建资源记录
func (r *ResourceRepository) Create(ctx context.Context, resourceSlug string, data map[string]interface{}) error {
	db := r.DB(ctx)
	tableName := resourceSlug
	processedData := r.processData(data)
	columns := ""
	values := ""
	vals := make([]interface{}, 0)
	i := 0
	for key, value := range processedData {
		if value == nil {
			continue
		}
		if i > 0 {
			columns += ", "
			values += ", "
		}
		columns += key
		values += "?"
		vals = append(vals, value)
		i++
	}
	if columns != "" {
		columns += ", "
		values += ", "
	}
	columns += "created_at, updated_at"
	values += "?, ?"
	vals = append(vals, time.Now(), time.Now())
	query := "INSERT INTO " + tableName + " (" + columns + ") VALUES (" + values + ")"
	return db.Exec(query, vals...).Error
}

// Update 更新资源记录
func (r *ResourceRepository) Update(ctx context.Context, resourceSlug string, id interface{}, data map[string]interface{}) error {
	db := r.DB(ctx)
	tableName := resourceSlug
	processedData := r.processData(data)
	setClause := ""
	vals := make([]interface{}, 0)
	i := 0
	for key, value := range processedData {
		if value == nil {
			continue
		}
		if i > 0 {
			setClause += ", "
		}
		setClause += key + " = ?"
		vals = append(vals, value)
		i++
	}
	if setClause != "" {
		setClause += ", "
	}
	setClause += "updated_at = ?"
	vals = append(vals, time.Now())
	vals = append(vals, id)
	query := "UPDATE " + tableName + " SET " + setClause + " WHERE id = ?"
	return db.Exec(query, vals...).Error
}

// Delete 删除资源记录（软删除）
func (r *ResourceRepository) Delete(ctx context.Context, resourceSlug string, id interface{}) error {
	db := r.DB(ctx)
	tableName := resourceSlug
	query := "UPDATE " + tableName + " SET deleted_at = ? WHERE id = ?"
	return db.Exec(query, time.Now(), id).Error
}

// Restore 恢复软删除
func (r *ResourceRepository) Restore(ctx context.Context, resourceSlug string, id interface{}) error {
	db := r.DB(ctx)
	tableName := resourceSlug
	query := "UPDATE " + tableName + " SET deleted_at = NULL WHERE id = ?"
	return db.Exec(query, id).Error
}

// ForceDelete 强制删除（硬删）
func (r *ResourceRepository) ForceDelete(ctx context.Context, resourceSlug string, id interface{}) error {
	db := r.DB(ctx)
	tableName := resourceSlug
	query := "DELETE FROM " + tableName + " WHERE id = ?"
	return db.Exec(query, id).Error
}

// DeleteBatch 批量删除资源记录
func (r *ResourceRepository) DeleteBatch(ctx context.Context, resourceSlug string, ids []interface{}) (int64, error) {
	db := r.DB(ctx)
	tableName := resourceSlug
	placeholders := make([]string, len(ids))
	vals := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		vals[i] = id
	}
	query := "UPDATE " + tableName + " SET deleted_at = ? WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	vals = append([]interface{}{time.Now()}, vals...)
	result := db.Exec(query, vals...)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// FindByID 根据 ID 查找资源记录
func (r *ResourceRepository) FindByID(ctx context.Context, resourceSlug string, id interface{}) (map[string]interface{}, error) {
	db := r.DB(ctx)
	tableName := resourceSlug
	var result map[string]interface{}
	query := "SELECT * FROM " + tableName + " WHERE id = ?"
	err := db.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List 获取资源记录列表
func (r *ResourceRepository) List(ctx context.Context, resourceSlug string, page, pageSize int) ([]map[string]interface{}, int64, error) {
	db := r.DB(ctx)
	tableName := resourceSlug
	var total int64
	countQuery := "SELECT COUNT(*) FROM " + tableName
	err := db.Raw(countQuery).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	var results []map[string]interface{}
	query := "SELECT * FROM " + tableName + " ORDER BY id DESC LIMIT ? OFFSET ?"
	err = db.Raw(query, pageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

// ListWithFilters 获取资源记录列表，支持过滤和搜索
func (r *ResourceRepository) ListWithFilters(
	ctx context.Context,
	resourceSlug string,
	page, pageSize int,
	filters map[string]interface{}, // 精确过滤条件
	search map[string]interface{}, // 模糊搜索条件
	orderBy string, // 排序字段
	orderDirection string, // 排序方向 ASC/DESC
) ([]map[string]interface{}, int64, error) {
	db := r.DB(ctx)
	tableName := resourceSlug
	whereClause := ""
	vals := make([]interface{}, 0)

	// 软删除视图
	trashedMode := "without"
	if v, ok := filters["trashed"]; ok {
		if s, ok2 := v.(string); ok2 {
			trashedMode = s
		}
		delete(filters, "trashed")
	}
	// 时间范围过滤（白名单）
	var createdFrom interface{}
	var createdTo interface{}
	if v, ok := filters["created_at_from"]; ok {
		createdFrom = v
	}
	if v, ok := filters["created_at_to"]; ok {
		createdTo = v
	}

	// 处理过滤条件（精确匹配）
	for field, value := range filters {
		if field == "created_at_from" || field == "created_at_to" {
			continue
		}
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += field + " = ?"
		vals = append(vals, value)
	}

	// 附加 created_at 范围
	if createdFrom != nil {
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += "created_at >= ?"
		vals = append(vals, createdFrom)
	}
	if createdTo != nil {
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += "created_at <= ?"
		vals = append(vals, createdTo)
	}
	// 处理搜索条件（模糊匹配）
	for field, value := range search {
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += field + " LIKE ?"
		vals = append(vals, "%"+value.(string)+"%")
	}
	// 处理 deleted_at 过滤
	switch trashedMode {
	case "only":
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += "deleted_at IS NOT NULL"
	case "with":
		// 不加条件
	default: // without
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += "deleted_at IS NULL"
	}

	wherePart := ""
	if whereClause != "" {
		wherePart = " WHERE " + whereClause
	}

	var total int64
	countQuery := "SELECT COUNT(*) FROM " + tableName + wherePart
	err := db.Raw(countQuery, vals...).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	orderClause := "id DESC"
	if orderBy != "" && (orderDirection == "ASC" || orderDirection == "DESC") {
		orderClause = orderBy + " " + orderDirection
	}
	var results []map[string]interface{}
	query := "SELECT * FROM " + tableName + wherePart + " ORDER BY " + orderClause + " LIMIT ? OFFSET ?"
	pageVals := append(vals, pageSize, offset)
	err = db.Raw(query, pageVals...).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

// ListWithRelationships 获取资源记录列表，包含关联数据
func (r *ResourceRepository) ListWithRelationships(
	ctx context.Context,
	resourceSlug string,
	page, pageSize int,
	relationships map[string]string, // 关联字段名 -> 关联资源名
) ([]map[string]interface{}, int64, error) {
	db := r.DB(ctx)
	tableName := resourceSlug
	var total int64
	countQuery := "SELECT COUNT(*) FROM " + tableName
	err := db.Raw(countQuery).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	var results []map[string]interface{}
	query := "SELECT * FROM " + tableName + " ORDER BY id DESC LIMIT ? OFFSET ?"
	err = db.Raw(query, pageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}
	if len(relationships) > 0 && len(results) > 0 {
		for _, result := range results {
			for relField, relResource := range relationships {
				if relID, ok := result[relField]; ok && relID != nil {
					relatedData, err := r.FindByID(ctx, relResource, relID)
					if err == nil {
						result[relField+"_data"] = relatedData
					}
				}
			}
		}
	}
	return results, total, nil
}

// ListWithRelationshipsAndFilters 获取资源记录列表，包含关联数据和过滤搜索功能
func (r *ResourceRepository) ListWithRelationshipsAndFilters(
	ctx context.Context,
	resourceSlug string,
	page, pageSize int,
	relationships map[string]string, // 关联字段名 -> 关联资源名
	filters map[string]interface{}, // 精确过滤条件
	search map[string]interface{}, // 模糊搜索条件
	orderBy string, // 排序字段
	orderDirection string, // 排序方向 ASC/DESC
) ([]map[string]interface{}, int64, error) {
	// 先获取数据（带过滤和搜索）
	results, total, err := r.ListWithFilters(ctx, resourceSlug, page, pageSize, filters, search, orderBy, orderDirection)
	if err != nil {
		return nil, 0, err
	}
	// 处理关联数据
	if len(relationships) > 0 && len(results) > 0 {
		for _, result := range results {
			for relField, relResource := range relationships {
				if relID, ok := result[relField]; ok && relID != nil {
					// 查询关联数据
					relatedData, err := r.FindByID(ctx, relResource, relID)
					if err == nil {
						result[relField+"_data"] = relatedData
					}
				}
			}
		}
	}
	return results, total, nil
}

// FindByIDWithRelationships 根据 ID 查找资源记录，包含关联数据
func (r *ResourceRepository) FindByIDWithRelationships(
	ctx context.Context,
	resourceSlug string,
	id interface{},
	relationships map[string]string, // 关联字段名 -> 关联资源名
) (map[string]interface{}, error) {
	// 先获取主记录
	result, err := r.FindByID(ctx, resourceSlug, id)
	if err != nil {
		return nil, err
	}
	// 处理关联数据
	if len(relationships) > 0 && result != nil {
		for relField, relResource := range relationships {
			if relID, ok := result[relField]; ok && relID != nil {
				// 查询关联数据
				relatedData, err := r.FindByID(ctx, relResource, relID)
				if err == nil {
					result[relField+"_data"] = relatedData
				}
			}
		}
	}
	return result, nil
}

// QuickSearch 在指定资源的可搜索字段中按关键字进行 OR 模糊查询，限制返回条数
func (r *ResourceRepository) QuickSearch(
	ctx context.Context,
	resourceSlug string,
	fields []string,
	keyword string,
	limit int,
) ([]map[string]interface{}, error) {
	db := r.DB(ctx)
	tableName := resourceSlug
	if limit <= 0 {
		limit = 5
	}
	whereClause := "deleted_at IS NULL"
	vals := make([]interface{}, 0)
	if len(fields) > 0 && keyword != "" {
		whereClause += " AND ("
		for i, f := range fields {
			if i > 0 {
				whereClause += " OR "
			}
			whereClause += f + " LIKE ?"
			vals = append(vals, "%"+keyword+"%")
		}
		whereClause += ")"
	}
	query := "SELECT * FROM " + tableName + " WHERE " + whereClause + " ORDER BY id DESC LIMIT ?"
	vals = append(vals, limit)
	var results []map[string]interface{}
	if err := db.Raw(query, vals...).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// processData 处理数据，转换特殊字段类型
func (r *ResourceRepository) processData(data map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})
	for key, value := range data {
		switch v := value.(type) {
		case time.Time:
			// 时间类型转换为字符串
			processed[key] = v.Format("2006-01-02 15:04:05")
		case bool:
			// 布尔类型保持不变
			processed[key] = v
		case int, int8, int16, int32, int64:
			// 整数类型保持不变
			processed[key] = v
		case float32, float64:
			// 浮点数类型保持不变
			processed[key] = v
		default:
			// 其他类型转换为字符串
			if v != nil {
				processed[key] = v
			}
		}
	}
	return processed
}
