package i18n

import (
	"fmt"
	"sync"
)

// ResourceManager 语言资源管理器
type ResourceManager struct {
	mu              sync.RWMutex
	resources       map[string]map[string]string // language -> key -> translation
	defaultLanguage string
}

// NewResourceManager 创建语言资源管理器
func NewResourceManager(defaultLanguage string) *ResourceManager {
	return &ResourceManager{
		resources:       make(map[string]map[string]string),
		defaultLanguage: defaultLanguage,
	}
}

// AddResource 添加语言资源
func (rm *ResourceManager) AddResource(language string, resources map[string]string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.resources[language] == nil {
		rm.resources[language] = make(map[string]string)
	}

	for key, value := range resources {
		rm.resources[language][key] = value
	}
}

// Translate 翻译文本
func (rm *ResourceManager) Translate(language, key string, args ...interface{}) string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	// 首先尝试在指定语言中查找
	if translations, ok := rm.resources[language]; ok {
		if translation, ok := translations[key]; ok {
			return rm.formatTranslation(translation, args...)
		}
	}

	// 如果在指定语言中找不到，尝试在默认语言中查找
	if language != rm.defaultLanguage {
		if translations, ok := rm.resources[rm.defaultLanguage]; ok {
			if translation, ok := translations[key]; ok {
				return rm.formatTranslation(translation, args...)
			}
		}
	}

	// 如果都找不到，返回键名
	return key
}

// formatTranslation 格式化翻译文本
func (rm *ResourceManager) formatTranslation(translation string, args ...interface{}) string {
	if len(args) == 0 {
		return translation
	}

	return fmt.Sprintf(translation, args...)
}

// GetLanguages 获取支持的语言列表
func (rm *ResourceManager) GetLanguages() []string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	languages := make([]string, 0, len(rm.resources))
	for language := range rm.resources {
		languages = append(languages, language)
	}

	return languages
}

// GlobalResourceManager 全局语言资源管理器实例
var GlobalResourceManager = NewResourceManager("zh-CN")

// AddResource 添加语言资源（全局）
func AddResource(language string, resources map[string]string) {
	GlobalResourceManager.AddResource(language, resources)
}

// Translate 翻译文本（全局）
func Translate(language, key string, args ...interface{}) string {
	return GlobalResourceManager.Translate(language, key, args...)
}
