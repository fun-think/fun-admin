package handler

import (
	"fmt"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ImportHandler 导入处理器
type ImportHandler struct {
	importService *service.ImportService
}

// NewImportHandler 创建导入处理器
func NewImportHandler(importService *service.ImportService) *ImportHandler {
	return &ImportHandler{
		importService: importService,
	}
}

// ImportData 导入数据
// @Summary 导入数据
// @Description 导入Excel、CSV或JSON格式的数据
// @Tags import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "导入文件"
// @Param type formData string true "文件类型 (excel/csv/json)"
// @Param resource formData string true "资源名称"
// @Param has_header formData bool false "是否包含表头" default(true)
// @Param sheet_name formData string false "工作表名称 (仅Excel)"
// @Param start_row formData int false "开始行" default(0)
// @Success 200 {object} v1.Response{data=service.ImportResult}
// @Router /api/v1/import/data [post]
func (h *ImportHandler) ImportData(c *gin.Context) {
	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		v1.HandleValidationError(c, "获取文件失败")
		return
	}

	// 获取参数
	fileType := c.PostForm("type")
	resource := c.PostForm("resource")
	hasHeader := c.DefaultPostForm("has_header", "true") == "true"
	sheetName := c.PostForm("sheet_name")
	startRow := 0
	if sr := c.PostForm("start_row"); sr != "" {
		startRow, _ = strconv.Atoi(sr)
	}

	// 验证参数
	if fileType == "" {
		v1.HandleValidationError(c, "文件类型不能为空")
		return
	}

	if resource == "" {
		v1.HandleValidationError(c, "资源名称不能为空")
		return
	}

	// 根据资源类型获取导入配置
	option, err := h.getImportOption(resource, hasHeader, sheetName, startRow)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	// 执行导入
	var result *service.ImportResult
	switch fileType {
	case "excel":
		result, err = h.importService.ImportExcel(c, file, *option)
	case "csv":
		result, err = h.importService.ImportCSV(c, file, *option)
	case "json":
		result, err = h.importService.ImportJSON(c, file, *option)
	default:
		v1.HandleValidationError(c, "不支持的文件类型")
		return
	}

	if err != nil {
		v1.HandleError(c, fmt.Errorf("导入失败: %w", err))
		return
	}

	v1.HandleSuccess(c, result)
}

// GetExcelSheets 获取Excel工作表列表
// @Summary 获取Excel工作表列表
// @Description 获取Excel文件中的所有工作表名称
// @Tags import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel文件"
// @Success 200 {object} v1.Response{data=[]string}
// @Router /api/v1/import/excel/sheets [post]
func (h *ImportHandler) GetExcelSheets(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		v1.HandleValidationError(c, "获取文件失败")
		return
	}

	sheets, err := h.importService.GetExcelSheets(file)
	if err != nil {
		v1.HandleError(c, fmt.Errorf("获取工作表列表失败: %w", err))
		return
	}

	v1.HandleSuccess(c, sheets)
}

// GetExcelColumns 获取Excel列信息
// @Summary 获取Excel列信息
// @Description 获取Excel文件指定工作表的列信息
// @Tags import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel文件"
// @Param sheet_name formData string true "工作表名称"
// @Param has_header formData bool false "是否包含表头" default(true)
// @Success 200 {object} v1.Response{data=[]string}
// @Router /api/v1/import/excel/columns [post]
func (h *ImportHandler) GetExcelColumns(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		v1.HandleValidationError(c, "获取文件失败")
		return
	}

	sheetName := c.PostForm("sheet_name")
	if sheetName == "" {
		v1.HandleValidationError(c, "工作表名称不能为空")
		return
	}

	hasHeader := c.DefaultPostForm("has_header", "true") == "true"

	columns, err := h.importService.GetExcelColumns(file, sheetName, hasHeader)
	if err != nil {
		v1.HandleError(c, fmt.Errorf("获取列信息失败: %w", err))
		return
	}

	v1.HandleSuccess(c, columns)
}

// getImportOption 获取导入配置
func (h *ImportHandler) getImportOption(resource string, hasHeader bool, sheetName string, startRow int) (*service.ImportOption, error) {
	// TODO: 根据资源类型从数据库或配置中获取导入配置
	// 这里提供一些示例配置

	switch resource {
	case "user":
		return &service.ImportOption{
			HasHeader: hasHeader,
			SheetName: sheetName,
			StartRow:  startRow,
			FieldMapping: map[string]string{
				"用户名": "username",
				"邮箱":  "email",
				"手机号": "phone",
				"姓名":  "real_name",
				"状态":  "status",
				"角色":  "role",
			},
			ValidateFunc: h.validateUserData,
			BatchSize:    100,
		}, nil
	case "product":
		return &service.ImportOption{
			HasHeader: hasHeader,
			SheetName: sheetName,
			StartRow:  startRow,
			FieldMapping: map[string]string{
				"商品名称": "name",
				"商品编码": "code",
				"价格":   "price",
				"库存":   "stock",
				"分类":   "category",
				"状态":   "status",
			},
			ValidateFunc: h.validateProductData,
			BatchSize:    100,
		}, nil
	default:
		return nil, fmt.Errorf("不支持的资源类型: %s", resource)
	}
}

// validateUserData 验证用户数据
func (h *ImportHandler) validateUserData(data map[string]interface{}) error {
	// 验证用户名
	if username, ok := data["username"].(string); ok {
		if username == "" {
			return fmt.Errorf("用户名不能为空")
		}
	}

	// 验证邮箱
	if email, ok := data["email"].(string); ok {
		if email == "" {
			return fmt.Errorf("邮箱不能为空")
		}
		// TODO: 添加邮箱格式验证
	}

	return nil
}

// validateProductData 验证商品数据
func (h *ImportHandler) validateProductData(data map[string]interface{}) error {
	// 验证商品名称
	if name, ok := data["name"].(string); ok {
		if name == "" {
			return fmt.Errorf("商品名称不能为空")
		}
	}

	// 验证价格
	if price, ok := data["price"].(float64); ok {
		if price < 0 {
			return fmt.Errorf("价格不能为负数")
		}
	}

	// 验证库存
	if stock, ok := data["stock"].(int); ok {
		if stock < 0 {
			return fmt.Errorf("库存不能为负数")
		}
	}

	return nil
}
