package handler

import (
	"fmt"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/service"
	"fun-admin/pkg/admin"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ImportHandler 导入处理器
type ImportHandler struct {
	importService   *service.ImportService
	resourceService *service.ResourceService
}

// NewImportHandler 创建导入处理器
func NewImportHandler(importService *service.ImportService, resourceService *service.ResourceService) *ImportHandler {
	return &ImportHandler{
		importService:   importService,
		resourceService: resourceService,
	}
}

// ImportData 导入数据
func (h *ImportHandler) ImportData(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		v1.HandleValidationError(c, "获取文件失败")
		return
	}

	fileType := c.PostForm("type")
	resourceSlug := c.Param("resource")
	if resourceSlug == "" {
		resourceSlug = c.PostForm("resource")
	}
	hasHeader := c.DefaultPostForm("has_header", "true") == "true"
	sheetName := c.PostForm("sheet_name")
	startRow := 0
	if sr := c.PostForm("start_row"); sr != "" {
		startRow, _ = strconv.Atoi(sr)
	}

	if fileType == "" {
		v1.HandleValidationError(c, "文件类型不能为空")
		return
	}
	if resourceSlug == "" {
		v1.HandleValidationError(c, "资源标识不能为空")
		return
	}

	option, err := h.getImportOption(c, resourceSlug, hasHeader, sheetName, startRow)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

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

// GetExcelSheets 返回 Excel 工作表列表
func (h *ImportHandler) GetExcelSheets(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		v1.HandleValidationError(c, "获取文件失败")
		return
	}

	sheets, err := h.importService.GetExcelSheets(file)
	if err != nil {
		v1.HandleError(c, fmt.Errorf("获取工作表失败: %w", err))
		return
	}

	v1.HandleSuccess(c, sheets)
}

// GetExcelColumns 返回 Excel 列名
func (h *ImportHandler) GetExcelColumns(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		v1.HandleValidationError(c, "获取文件失败")
		return
	}
	sheetName := c.PostForm("sheet_name")
	if sheetName == "" {
		v1.HandleValidationError(c, "sheet_name 不能为空")
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

func (h *ImportHandler) getImportOption(ctx *gin.Context, resourceSlug string, hasHeader bool, sheetName string, startRow int) (*service.ImportOption, error) {
	res := admin.GlobalResourceManager.GetResourceBySlug(resourceSlug)
	if res == nil {
		return nil, fmt.Errorf("资源 %s 不存在", resourceSlug)
	}

	fieldMapping := make(map[string]string)
	for _, field := range res.GetFields() {
		fieldMapping[field.GetName()] = field.GetName()
		label := strings.TrimSpace(field.GetLabel())
		if label != "" {
			fieldMapping[label] = field.GetName()
		}
	}

	option := &service.ImportOption{
		HasHeader:    hasHeader,
		SheetName:    sheetName,
		StartRow:     startRow,
		FieldMapping: fieldMapping,
		BatchSize:    100,
	}

	option.ValidateFunc = func(row map[string]interface{}) error {
		errors := admin.ValidateResourceData(res, row)
		if len(errors) == 0 {
			return nil
		}
		var parts []string
		for field, msgs := range errors {
			parts = append(parts, fmt.Sprintf("%s: %s", field, strings.Join(msgs, ",")))
		}
		return fmt.Errorf("校验失败: %s", strings.Join(parts, "; "))
	}

	option.DataHandler = func(batch []map[string]interface{}) error {
		for _, data := range batch {
			if _, err := h.resourceService.Create(ctx, resourceSlug, data); err != nil {
				return err
			}
		}
		return nil
	}

	return option, nil
}
