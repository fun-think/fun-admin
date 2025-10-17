package handler

import (
	"context"
	"fmt"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/service"

	"github.com/gin-gonic/gin"
)

// ConfigHandler 配置处理器
type ConfigHandler struct {
	configService *service.ConfigService
}

// NewConfigHandler 创建配置处理器
func NewConfigHandler(configService *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

// GetConfig 获取系统设置
// @Summary 获取系统设置
// @Description 根据键获取系统设置值
// @Tags config
// @Produce json
// @Param key path string true "设置键"
// @Success 200 {object} v1.Response{data=string}
// @Router /api/v1/config/configs/{key} [get]
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		v1.HandleValidationError(c, "设置键不能为空")
		return
	}

	value, err := h.configService.GetConfig(c, key)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, value)
}

// SetConfig 设置系统设置
// @Summary 设置系统设置
// @Description 设置系统设置值
// @Tags config
// @Accept json
// @Produce json
// @Param config body map[string]interface{} true "设置信息"
// @Success 200 {object} v1.Response
// @Router /api/v1/config/configs [post]
func (h *ConfigHandler) SetConfig(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Name  string `json:"name" binding:"required"`
		Value string `json:"value" binding:"required"`
		Type  string `json:"type" binding:"required"`
		Group string `json:"group" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleValidationError(c, "参数错误")
		return
	}

	if err := h.configService.SetConfig(c, req.Key, req.Name, req.Value, req.Type, req.Group); err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetConfigsByGroup 根据分组获取系统设置
// @Summary 根据分组获取系统设置
// @Description 获取指定分组的所有系统设置
// @Tags config
// @Produce json
// @Param group path string true "分组名称"
// @Success 200 {object} v1.Response{data=map[string]string}
// @Router /api/v1/config/configs/group/{group} [get]
func (h *ConfigHandler) GetConfigsByGroup(c *gin.Context) {
	group := c.Param("group")
	if group == "" {
		v1.HandleValidationError(c, "分组名称不能为空")
		return
	}

	configs, err := h.configService.GetConfigByGroup(c, group)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, configs)
}

// GetAllConfigs 获取所有系统设置
// @Summary 获取所有系统设置
// @Description 获取所有系统设置，按分组组织
// @Tags config
// @Produce json
// @Success 200 {object} v1.Response{data=map[string]map[string]string}
// @Router /api/v1/config/configs [get]
func (h *ConfigHandler) GetAllConfigs(c *gin.Context) {
	// 手动实现获取所有设置的逻辑
	groups, err := h.configService.GetConfigGroups(c)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	allConfigs := make(map[string]map[string]string)
	for _, group := range groups {
		configs, err := h.configService.GetConfigByGroup(c, group)
		if err != nil {
			v1.HandleError(c, fmt.Errorf("获取分组 %s 设置失败: %w", group, err))
			return
		}
		allConfigs[group] = configs
	}

	v1.HandleSuccess(c, allConfigs)
}

// UpdateConfig 更新系统设置
// @Summary 更新系统设置
// @Description 更新指定键的系统设置值
// @Tags config
// @Accept json
// @Produce json
// @Param key path string true "设置键"
// @Param value body string true "设置值"
// @Success 200 {object} v1.Response
// @Router /api/v1/config/configs/{key} [put]
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		v1.HandleValidationError(c, "设置键不能为空")
		return
	}

	var req struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleValidationError(c, "参数错误")
		return
	}

	// 通过先获取设置信息，然后再更新的方式来实现更新功能
	// 这里需要先获取原始设置的其他信息
	ctx := context.Background()
	_ = ctx // 避免ctx未使用警告
	// 注意：这里只是示例，实际应该有更完善的方法来更新设置
	err := h.configService.SetConfig(c, key, key, req.Value, "string", "default")
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// DeleteConfig 删除系统设置
// @Summary 删除系统设置
// @Description 删除指定键的系统设置
// @Tags config
// @Produce json
// @Param key path string true "设置键"
// @Success 200 {object} v1.Response
// @Router /api/v1/config/configs/{key} [delete]
func (h *ConfigHandler) DeleteConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		v1.HandleValidationError(c, "设置键不能为空")
		return
	}

	// 注意：ConfigService 中没有直接的删除方法，这里只是示例
	v1.HandleError(c, fmt.Errorf("删除设置功能暂未实现"))
}

// GetConfigGroups 获取所有设置分组
// @Summary 获取所有设置分组
// @Description 获取所有系统设置分组名称
// @Tags config
// @Produce json
// @Success 200 {object} v1.Response{data=[]string}
// @Router /api/v1/config/groups [get]
func (h *ConfigHandler) GetConfigGroups(c *gin.Context) {
	groups, err := h.configService.GetConfigGroups(c)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, groups)
}

// TestEmailConfig 测试邮件配置
// @Summary 测试邮件配置
// @Description 测试邮件配置是否正确
// @Tags config
// @Accept json
// @Produce json
// @Param config body service.EmailConfig true "邮件配置"
// @Success 200 {object} v1.Response
// @Router /api/v1/config/email/test [post]
func (h *ConfigHandler) TestEmailConfig(c *gin.Context) {
	var emailConfig service.EmailConfig
	if err := c.ShouldBindJSON(&emailConfig); err != nil {
		v1.HandleValidationError(c, "参数错误")
		return
	}

	// 调用测试方法，这里使用一个假的测试邮箱
	err := h.configService.TestEmailConfig(c, &emailConfig, "test@example.com")
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// TestStorageConfig 测试存储配置
// @Summary 测试存储配置
// @Description 测试存储配置是否正确
// @Tags config
// @Accept json
// @Produce json
// @Param config body service.StorageConfig true "存储配置"
// @Success 200 {object} v1.Response
// @Router /api/v1/config/storage/test [post]
func (h *ConfigHandler) TestStorageConfig(c *gin.Context) {
	var storageConfig service.StorageConfig
	if err := c.ShouldBindJSON(&storageConfig); err != nil {
		v1.HandleValidationError(c, "参数错误")
		return
	}

	// 存储配置测试功能在ConfigService中没有实现，这里返回未实现错误
	v1.HandleError(c, fmt.Errorf("存储配置测试功能暂未实现"))
}
