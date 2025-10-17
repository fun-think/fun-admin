package repository

import (
	"context"
	"errors"
	"fmt"
	"fun-admin/internal/model"
	"fun-admin/pkg/logger"
	"gorm.io/gorm"
)

// DictionaryRepository 字典仓库接口
type DictionaryRepository interface {
	CreateDictionaryType(ctx context.Context, dictionaryType *model.DictionaryType) error
	UpdateDictionaryType(ctx context.Context, dictionaryType *model.DictionaryType) error
	DeleteDictionaryType(ctx context.Context, id uint) error
	GetDictionaryType(ctx context.Context, id uint) (*model.DictionaryType, error)
	GetDictionaryTypeByCode(ctx context.Context, code string) (*model.DictionaryType, error)
	ListDictionaryTypes(ctx context.Context, page, pageSize int, name string) ([]model.DictionaryType, int64, error)
	GetDictionaryTypes(ctx context.Context) ([]model.DictionaryType, error)
	CreateDictionaryData(ctx context.Context, dictionaryData *model.DictionaryData) error
	UpdateDictionaryData(ctx context.Context, dictionaryData *model.DictionaryData) error
	DeleteDictionaryData(ctx context.Context, id uint) error
	GetDictionaryData(ctx context.Context, id uint) (*model.DictionaryData, error)
	ListDictionaryData(ctx context.Context, typeID uint, label string, status int) ([]model.DictionaryData, error)
	GetDictByCode(ctx context.Context, code string) ([]model.DictionaryData, error)
	GetDictValue(ctx context.Context, code, value string) (string, error)
	GetDictLabel(ctx context.Context, code, label string) (string, error)
	GetDefaultDictValue(ctx context.Context, code string) (string, error)
	ClearCacheByTypeID(ctx context.Context, typeID uint) error
}

// dictionaryRepository 字典仓库实现
type dictionaryRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

// NewDictionaryRepository 创建字典仓库
func NewDictionaryRepository(
	logger *logger.Logger,
	db *gorm.DB,
) DictionaryRepository {
	return &dictionaryRepository{
		logger: logger,
		db:     db,
	}
}

// CreateDictionaryType 创建字典类型
func (r *dictionaryRepository) CreateDictionaryType(ctx context.Context, dictionaryType *model.DictionaryType) error {
	if err := r.db.WithContext(ctx).Create(dictionaryType).Error; err != nil {
		return fmt.Errorf("创建字典类型失败: %w", err)
	}
	return nil
}

// UpdateDictionaryType 更新字典类型
func (r *dictionaryRepository) UpdateDictionaryType(ctx context.Context, dictionaryType *model.DictionaryType) error {
	if err := r.db.WithContext(ctx).Save(dictionaryType).Error; err != nil {
		return fmt.Errorf("更新字典类型失败: %w", err)
	}
	return nil
}

// DeleteDictionaryType 删除字典类型
func (r *dictionaryRepository) DeleteDictionaryType(ctx context.Context, id uint) error {
	// 检查是否有字典数据
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.DictionaryData{}).Where("type_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("检查字典数据失败: %w", err)
	}

	if count > 0 {
		return errors.New("该字典类型下存在字典数据，无法删除")
	}

	// 获取字典类型编码
	var dictionaryType model.DictionaryType
	if err := r.db.WithContext(ctx).First(&dictionaryType, id).Error; err != nil {
		return fmt.Errorf("获取字典类型失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Delete(&model.DictionaryType{}, id).Error; err != nil {
		return fmt.Errorf("删除字典类型失败: %w", err)
	}

	return nil
}

// GetDictionaryType 获取字典类型
func (r *dictionaryRepository) GetDictionaryType(ctx context.Context, id uint) (*model.DictionaryType, error) {
	var dictionaryType model.DictionaryType
	if err := r.db.WithContext(ctx).First(&dictionaryType, id).Error; err != nil {
		return nil, fmt.Errorf("获取字典类型失败: %w", err)
	}
	return &dictionaryType, nil
}

// GetDictionaryTypeByCode 根据编码获取字典类型
func (r *dictionaryRepository) GetDictionaryTypeByCode(ctx context.Context, code string) (*model.DictionaryType, error) {
	var dictionaryType model.DictionaryType
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&dictionaryType).Error; err != nil {
		return nil, fmt.Errorf("获取字典类型失败: %w", err)
	}
	return &dictionaryType, nil
}

// ListDictionaryTypes 获取字典类型列表
func (r *dictionaryRepository) ListDictionaryTypes(ctx context.Context, page, pageSize int, name string) ([]model.DictionaryType, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.DictionaryType{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取字典类型总数失败: %w", err)
	}

	var dictionaryTypes []model.DictionaryType
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("sort ASC, id ASC").Find(&dictionaryTypes).Error; err != nil {
		return nil, 0, fmt.Errorf("获取字典类型列表失败: %w", err)
	}

	return dictionaryTypes, total, nil
}

// GetDictionaryTypes 获取所有字典类型
func (r *dictionaryRepository) GetDictionaryTypes(ctx context.Context) ([]model.DictionaryType, error) {
	var dictionaryTypes []model.DictionaryType
	if err := r.db.WithContext(ctx).Model(&model.DictionaryType{}).Find(&dictionaryTypes).Error; err != nil {
		return nil, fmt.Errorf("获取字典类型列表失败: %w", err)
	}
	return dictionaryTypes, nil
}

// CreateDictionaryData 创建字典数据
func (r *dictionaryRepository) CreateDictionaryData(ctx context.Context, dictionaryData *model.DictionaryData) error {
	if err := r.db.WithContext(ctx).Create(dictionaryData).Error; err != nil {
		return fmt.Errorf("创建字典数据失败: %w", err)
	}

	// 如果是默认值，清除其他默认值
	if dictionaryData.IsDefault {
		r.db.WithContext(ctx).Model(&model.DictionaryData{}).
			Where("type_id = ? AND id != ?", dictionaryData.TypeID, dictionaryData.ID).
			Update("is_default", false)
	}

	return nil
}

// UpdateDictionaryData 更新字典数据
func (r *dictionaryRepository) UpdateDictionaryData(ctx context.Context, dictionaryData *model.DictionaryData) error {
	if err := r.db.WithContext(ctx).Save(dictionaryData).Error; err != nil {
		return fmt.Errorf("更新字典数据失败: %w", err)
	}

	// 如果是默认值，清除其他默认值
	if dictionaryData.IsDefault {
		r.db.WithContext(ctx).Model(&model.DictionaryData{}).
			Where("type_id = ? AND id != ?", dictionaryData.TypeID, dictionaryData.ID).
			Update("is_default", false)
	}

	return nil
}

// DeleteDictionaryData 删除字典数据
func (r *dictionaryRepository) DeleteDictionaryData(ctx context.Context, id uint) error {
	var dictionaryData model.DictionaryData
	if err := r.db.WithContext(ctx).First(&dictionaryData, id).Error; err != nil {
		return fmt.Errorf("获取字典数据失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Delete(&dictionaryData).Error; err != nil {
		return fmt.Errorf("删除字典数据失败: %w", err)
	}

	return nil
}

// GetDictionaryData 获取字典数据
func (r *dictionaryRepository) GetDictionaryData(ctx context.Context, id uint) (*model.DictionaryData, error) {
	var dictionaryData model.DictionaryData
	if err := r.db.WithContext(ctx).Preload("Type").First(&dictionaryData, id).Error; err != nil {
		return nil, fmt.Errorf("获取字典数据失败: %w", err)
	}
	return &dictionaryData, nil
}

// ListDictionaryData 获取字典数据列表
func (r *dictionaryRepository) ListDictionaryData(ctx context.Context, typeID uint, label string, status int) ([]model.DictionaryData, error) {
	query := r.db.WithContext(ctx).Model(&model.DictionaryData{}).Where("type_id = ?", typeID)

	if label != "" {
		query = query.Where("label LIKE ?", "%"+label+"%")
	}

	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	var dictionaryData []model.DictionaryData
	if err := query.Order("sort ASC, id ASC").Find(&dictionaryData).Error; err != nil {
		return nil, fmt.Errorf("获取字典数据列表失败: %w", err)
	}

	return dictionaryData, nil
}

// GetDictByCode 根据字典编码获取字典数据
func (r *dictionaryRepository) GetDictByCode(ctx context.Context, code string) ([]model.DictionaryData, error) {
	var dictionaryType model.DictionaryType
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&dictionaryType).Error; err != nil {
		return nil, fmt.Errorf("获取字典类型失败: %w", err)
	}

	var dictionaryData []model.DictionaryData
	if err := r.db.WithContext(ctx).Where("type_id = ? AND status = 1", dictionaryType.ID).
		Order("sort ASC, id ASC").Find(&dictionaryData).Error; err != nil {
		return nil, fmt.Errorf("获取字典数据失败: %w", err)
	}

	return dictionaryData, nil
}

// GetDictValue 获取字典值对应的标签
func (r *dictionaryRepository) GetDictValue(ctx context.Context, code, value string) (string, error) {
	dictionaryData, err := r.GetDictByCode(ctx, code)
	if err != nil {
		return "", err
	}

	for _, data := range dictionaryData {
		if data.Value == value {
			return data.Label, nil
		}
	}

	return "", fmt.Errorf("字典值不存在: %s", value)
}

// GetDictLabel 获取字典标签对应的值
func (r *dictionaryRepository) GetDictLabel(ctx context.Context, code, label string) (string, error) {
	dictionaryData, err := r.GetDictByCode(ctx, code)
	if err != nil {
		return "", err
	}

	for _, data := range dictionaryData {
		if data.Label == label {
			return data.Value, nil
		}
	}

	return "", fmt.Errorf("字典标签不存在: %s", label)
}

// GetDefaultDictValue 获取字典类型的默认值
func (r *dictionaryRepository) GetDefaultDictValue(ctx context.Context, code string) (string, error) {
	dictionaryData, err := r.GetDictByCode(ctx, code)
	if err != nil {
		return "", err
	}

	for _, data := range dictionaryData {
		if data.IsDefault {
			return data.Value, nil
		}
	}

	return "", fmt.Errorf("字典类型没有默认值: %s", code)
}

// ClearCacheByTypeID 通过TypeID清除字典数据缓存
func (r *dictionaryRepository) ClearCacheByTypeID(ctx context.Context, typeID uint) error {
	// 获取字典类型编码
	var dictionaryType model.DictionaryType
	if err := r.db.WithContext(ctx).First(&dictionaryType, typeID).Error; err != nil {
		return err
	}
	return nil
}
