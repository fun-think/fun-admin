package handler

import (
	v1 "fun-admin/api/v1"
	"fun-admin/internal/model"
	"fun-admin/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DictionaryHandler 字典处理器
type DictionaryHandler struct {
	dictionaryService *service.DictionaryService
}

// NewDictionaryHandler 创建字典处理器
func NewDictionaryHandler(dictionaryService *service.DictionaryService) *DictionaryHandler {
	return &DictionaryHandler{
		dictionaryService: dictionaryService,
	}
}

// CreateDictionaryType 创建字典类型
// @Summary 创建字典类型
// @Description 创建新的字典类型
// @Tags dictionary
// @Accept json
// @Produce json
// @Param dict body model.DictionaryType true "字典类型信息"
// @Success 200 {object} v1.Response{data=model.DictionaryType}
// @Router /api/v1/dict/types [post]
func (h *DictionaryHandler) CreateDictionaryType(c *gin.Context) {
	var dictionaryType model.DictionaryType
	if err := c.ShouldBindJSON(&dictionaryType); err != nil {
		v1.HandleValidationError(c, err.Error())
		return
	}

	if err := h.dictionaryService.CreateDictionaryType(c, &dictionaryType); err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, dictionaryType)
}

// UpdateDictionaryType 更新字典类型
// @Summary 更新字典类型
// @Description 更新指定ID的字典类型
// @Tags dictionary
// @Accept json
// @Produce json
// @Param id path int true "字典类型ID"
// @Param dict body model.DictionaryType true "字典类型信息"
// @Success 200 {object} v1.Response{data=model.DictionaryType}
// @Router /api/v1/dict/types/{id} [put]
func (h *DictionaryHandler) UpdateDictionaryType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		v1.HandleValidationError(c, "无效的ID参数")
		return
	}

	var dictionaryType model.DictionaryType
	if err := c.ShouldBindJSON(&dictionaryType); err != nil {
		v1.HandleValidationError(c, err.Error())
		return
	}
	dictionaryType.ID = uint(id)

	if err := h.dictionaryService.UpdateDictionaryType(c, &dictionaryType); err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, dictionaryType)
}

// DeleteDictionaryType 删除字典类型
// @Summary 删除字典类型
// @Description 删除指定ID的字典类型
// @Tags dictionary
// @Produce json
// @Param id path int true "字典类型ID"
// @Success 200 {object} v1.Response
// @Router /api/v1/dict/types/{id} [delete]
func (h *DictionaryHandler) DeleteDictionaryType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		v1.HandleValidationError(c, "无效的ID参数")
		return
	}

	if err := h.dictionaryService.DeleteDictionaryType(c, uint(id)); err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetDictionaryType 获取字典类型详情
// @Summary 获取字典类型详情
// @Description 根据ID获取字典类型详情
// @Tags dictionary
// @Produce json
// @Param id path int true "字典类型ID"
// @Success 200 {object} v1.Response{data=model.DictionaryType}
// @Router /api/v1/dict/types/{id} [get]
func (h *DictionaryHandler) GetDictionaryType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		v1.HandleValidationError(c, "无效的ID参数")
		return
	}

	dictionaryType, err := h.dictionaryService.GetDictionaryType(c, uint(id))
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, dictionaryType)
}

// ListDictionaryTypes 获取字典类型列表
// @Summary 获取字典类型列表
// @Description 获取字典类型列表，支持分页和搜索
// @Tags dictionary
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param name query string false "字典类型名称"
// @Success 200 {object} v1.Response{data=v1.GetAdminUsersResponseData}
// @Router /api/v1/dict/types [get]
func (h *DictionaryHandler) ListDictionaryTypes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	name := c.Query("name")

	dictionaryTypes, total, err := h.dictionaryService.ListDictionaryTypes(c, page, pageSize, name)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, map[string]interface{}{
		"list":  dictionaryTypes,
		"total": total,
	})
}

// CreateDictionaryData 创建字典数据
// @Summary 创建字典数据
// @Description 创建新的字典数据
// @Tags dictionary
// @Accept json
// @Produce json
// @Param dict_data body model.DictionaryData true "字典数据信息"
// @Success 200 {object} v1.Response{data=model.DictionaryData}
// @Router /api/v1/dict/data [post]
func (h *DictionaryHandler) CreateDictionaryData(c *gin.Context) {
	var dictionaryData model.DictionaryData
	if err := c.ShouldBindJSON(&dictionaryData); err != nil {
		v1.HandleValidationError(c, err.Error())
		return
	}

	if err := h.dictionaryService.CreateDictionaryData(c, &dictionaryData); err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, dictionaryData)
}

// UpdateDictionaryData 更新字典数据
// @Summary 更新字典数据
// @Description 更新指定ID的字典数据
// @Tags dictionary
// @Accept json
// @Produce json
// @Param id path int true "字典数据ID"
// @Param dict_data body model.DictionaryData true "字典数据信息"
// @Success 200 {object} v1.Response{data=model.DictionaryData}
// @Router /api/v1/dict/data/{id} [put]
func (h *DictionaryHandler) UpdateDictionaryData(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		v1.HandleValidationError(c, "无效的ID参数")
		return
	}

	var dictionaryData model.DictionaryData
	if err := c.ShouldBindJSON(&dictionaryData); err != nil {
		v1.HandleValidationError(c, err.Error())
		return
	}
	dictionaryData.ID = uint(id)

	if err := h.dictionaryService.UpdateDictionaryData(c, &dictionaryData); err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, dictionaryData)
}

// DeleteDictionaryData 删除字典数据
// @Summary 删除字典数据
// @Description 删除指定ID的字典数据
// @Tags dictionary
// @Produce json
// @Param id path int true "字典数据ID"
// @Success 200 {object} v1.Response
// @Router /api/v1/dict/data/{id} [delete]
func (h *DictionaryHandler) DeleteDictionaryData(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		v1.HandleValidationError(c, "无效的ID参数")
		return
	}

	if err := h.dictionaryService.DeleteDictionaryData(c, uint(id)); err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// GetDictionaryData 获取字典数据详情
// @Summary 获取字典数据详情
// @Description 根据ID获取字典数据详情
// @Tags dictionary
// @Produce json
// @Param id path int true "字典数据ID"
// @Success 200 {object} v1.Response{data=model.DictionaryData}
// @Router /api/v1/dict/data/{id} [get]
func (h *DictionaryHandler) GetDictionaryData(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		v1.HandleValidationError(c, "无效的ID参数")
		return
	}

	dictionaryData, err := h.dictionaryService.GetDictionaryData(c, uint(id))
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, dictionaryData)
}

// ListDictionaryData 获取字典数据列表
// @Summary 获取字典数据列表
// @Description 获取字典数据列表，支持分页和搜索
// @Tags dictionary
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param type_id query int false "字典类型ID"
// @Param label query string false "字典数据标签"
// @Success 200 {object} v1.Response{data=v1.GetAdminUsersResponseData}
// @Router /api/v1/dict/data [get]
func (h *DictionaryHandler) ListDictionaryData(c *gin.Context) {
	typeID, _ := strconv.Atoi(c.Query("type_id"))
	label := c.Query("label")

	dictionaryData, err := h.dictionaryService.ListDictionaryData(c, uint(typeID), label, -1)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, map[string]interface{}{
		"list":  dictionaryData,
		"total": len(dictionaryData),
	})
}

// GetDictionaryDataByCode 根据字典类型编码获取字典数据列表
// @Summary 根据字典类型编码获取字典数据列表
// @Description 根据字典类型编码获取字典数据列表
// @Tags dictionary
// @Produce json
// @Param code path string true "字典类型编码"
// @Success 200 {object} v1.Response{data=[]model.DictionaryData}
// @Router /api/v1/dict/data/type/{code} [get]
func (h *DictionaryHandler) GetDictionaryDataByCode(c *gin.Context) {
	code := c.Param("code")
	dictionaryData, err := h.dictionaryService.GetDictByCode(c, code)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, dictionaryData)
}
