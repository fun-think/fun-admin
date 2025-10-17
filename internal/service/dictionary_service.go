package service

import (
	"context"
	"fmt"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"
	"fun-admin/pkg/cache"
	"fun-admin/pkg/logger"
	"sync"
	"time"
)

// DictionaryService 字典服务
type DictionaryService struct {
	dictionaryRepo repository.DictionaryRepository
	cache          cache.CacheManager
	logger         *logger.Logger
	dictMap        map[string][]model.DictionaryData // 字典缓存
	mu             sync.RWMutex
}

// NewDictionaryService 创建字典服务
func NewDictionaryService(dictionaryRepo repository.DictionaryRepository, logger *logger.Logger, cache cache.CacheManager) *DictionaryService {
	return &DictionaryService{
		dictionaryRepo: dictionaryRepo,
		logger:         logger,
		cache:          cache,
		dictMap:        make(map[string][]model.DictionaryData),
	}
}

// getCache 获取字典缓存
func (s *DictionaryService) getCache(code string) ([]model.DictionaryData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, exists := s.dictMap[code]
	return data, exists
}

// setCache 设置字典缓存
func (s *DictionaryService) setCache(code string, data []model.DictionaryData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dictMap[code] = data
}

// CreateDictionaryType 创建字典类型
func (s *DictionaryService) CreateDictionaryType(ctx context.Context, dictionaryType *model.DictionaryType) error {
	return s.dictionaryRepo.CreateDictionaryType(ctx, dictionaryType)
}

// UpdateDictionaryType 更新字典类型
func (s *DictionaryService) UpdateDictionaryType(ctx context.Context, dictionaryType *model.DictionaryType) error {
	if err := s.dictionaryRepo.UpdateDictionaryType(ctx, dictionaryType); err != nil {
		return err
	}

	// 清除缓存
	s.clearCache(dictionaryType.Code)

	return nil
}

// DeleteDictionaryType 删除字典类型
func (s *DictionaryService) DeleteDictionaryType(ctx context.Context, id uint) error {
	if err := s.dictionaryRepo.DeleteDictionaryType(ctx, id); err != nil {
		return err
	}

	// 获取字典类型编码以清除缓存
	dictionaryType, err := s.dictionaryRepo.GetDictionaryType(ctx, id)
	if err != nil {
		return fmt.Errorf("获取字典类型失败: %w", err)
	}

	// 清除缓存
	s.clearCache(dictionaryType.Code)

	return nil
}

// GetDictionaryType 获取字典类型
func (s *DictionaryService) GetDictionaryType(ctx context.Context, id uint) (*model.DictionaryType, error) {
	return s.dictionaryRepo.GetDictionaryType(ctx, id)
}

// GetDictionaryTypeByCode 根据编码获取字典类型
func (s *DictionaryService) GetDictionaryTypeByCode(ctx context.Context, code string) (*model.DictionaryType, error) {
	return s.dictionaryRepo.GetDictionaryTypeByCode(ctx, code)
}

// ListDictionaryTypes 获取字典类型列表
func (s *DictionaryService) ListDictionaryTypes(ctx context.Context, page, pageSize int, name string) ([]model.DictionaryType, int64, error) {
	return s.dictionaryRepo.ListDictionaryTypes(ctx, page, pageSize, name)
}

// GetDictionaryTypes 获取所有字典类型
func (s *DictionaryService) GetDictionaryTypes(ctx context.Context) ([]model.DictionaryType, error) {
	return s.dictionaryRepo.GetDictionaryTypes(ctx)
}

// CreateDictionaryData 创建字典数据
func (s *DictionaryService) CreateDictionaryData(ctx context.Context, dictionaryData *model.DictionaryData) error {
	if err := s.dictionaryRepo.CreateDictionaryData(ctx, dictionaryData); err != nil {
		return err
	}

	// 如果是默认值，清除其他默认值
	if dictionaryData.IsDefault {
		s.clearCacheByTypeID(ctx, dictionaryData.TypeID)
	}

	// 清除缓存
	s.clearCacheByTypeID(ctx, dictionaryData.TypeID)

	return nil
}

// UpdateDictionaryData 更新字典数据
func (s *DictionaryService) UpdateDictionaryData(ctx context.Context, dictionaryData *model.DictionaryData) error {
	if err := s.dictionaryRepo.UpdateDictionaryData(ctx, dictionaryData); err != nil {
		return err
	}

	// 如果是默认值，清除其他默认值
	if dictionaryData.IsDefault {
		s.clearCacheByTypeID(ctx, dictionaryData.TypeID)
	}

	// 清除缓存
	s.clearCacheByTypeID(ctx, dictionaryData.TypeID)

	return nil
}

// DeleteDictionaryData 删除字典数据
func (s *DictionaryService) DeleteDictionaryData(ctx context.Context, id uint) error {
	if err := s.dictionaryRepo.DeleteDictionaryData(ctx, id); err != nil {
		return err
	}

	// 获取字典数据以清除缓存
	data, err := s.dictionaryRepo.GetDictionaryData(ctx, id)
	if err != nil {
		return err
	}

	// 清除缓存
	s.clearCacheByTypeID(ctx, data.TypeID)

	return nil
}

// GetDictionaryData 获取字典数据
func (s *DictionaryService) GetDictionaryData(ctx context.Context, id uint) (*model.DictionaryData, error) {
	return s.dictionaryRepo.GetDictionaryData(ctx, id)
}

// ListDictionaryData 获取字典数据列表
func (s *DictionaryService) ListDictionaryData(ctx context.Context, typeID uint, label string, status int) ([]model.DictionaryData, error) {
	return s.dictionaryRepo.ListDictionaryData(ctx, typeID, label, status)
}

// GetDictByCode 根据字典编码获取字典数据
func (s *DictionaryService) GetDictByCode(ctx context.Context, code string) ([]model.DictionaryData, error) {
	// 先从缓存获取
	if data, exists := s.getCache(code); exists {
		return data, nil
	}

	// 从数据库获取
	data, err := s.dictionaryRepo.GetDictByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// 设置缓存
	s.setCache(code, data)

	return data, nil
}

// GetDictValue 获取字典值对应的标签
func (s *DictionaryService) GetDictValue(ctx context.Context, code, value string) (string, error) {
	return s.dictionaryRepo.GetDictValue(ctx, code, value)
}

// GetDictLabel 获取字典标签对应的值
func (s *DictionaryService) GetDictLabel(ctx context.Context, code, label string) (string, error) {
	return s.dictionaryRepo.GetDictLabel(ctx, code, label)
}

// GetDefaultDictValue 获取字典类型的默认值
func (s *DictionaryService) GetDefaultDictValue(ctx context.Context, code string) (string, error) {
	return s.dictionaryRepo.GetDefaultDictValue(ctx, code)
}

// clearCache 清除缓存
func (s *DictionaryService) clearCache(dictCode string) {
	cacheKey := fmt.Sprintf("dict:%s", dictCode)
	s.cache.Delete(context.Background(), cacheKey)

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.dictMap, dictCode)
}

// clearCacheByTypeID 通过TypeID清除字典数据缓存
func (s *DictionaryService) clearCacheByTypeID(ctx context.Context, typeID uint) {
	// 获取字典类型编码
	if err := s.dictionaryRepo.ClearCacheByTypeID(ctx, typeID); err != nil {
		s.logger.Error("清除字典缓存失败: " + err.Error())
	}
}

// LoadDictToCache 加载字典到缓存
func (s *DictionaryService) LoadDictToCache(ctx context.Context) error {
	dictionaryTypes, err := s.dictionaryRepo.GetDictionaryTypes(ctx)
	if err != nil {
		return fmt.Errorf("获取字典类型列表失败: %w", err)
	}

	for _, dictionaryType := range dictionaryTypes {
		dictionaryData, err := s.dictionaryRepo.GetDictByCode(ctx, dictionaryType.Code)
		if err != nil {
			s.logger.Error("加载字典数据失败: " + err.Error())
			continue
		}

		cacheKey := fmt.Sprintf("dict:%s", dictionaryType.Code)
		s.cache.Set(ctx, cacheKey, dictionaryData, time.Hour)

		s.mu.Lock()
		s.dictMap[dictionaryType.Code] = dictionaryData
		s.mu.Unlock()
	}

	s.logger.Info("字典数据加载完成")
	return nil
}

// GetDictFromMemory 从内存获取字典数据
func (s *DictionaryService) GetDictFromMemory(code string) ([]model.DictionaryData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.dictMap[code]
	return data, exists
}
