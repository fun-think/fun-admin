package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"fun-admin/pkg/logger"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ImportService 数据导入服务
type ImportService struct {
	logger *logger.Logger
}

// NewImportService 创建导入服务
func NewImportService(logger *logger.Logger) *ImportService {
	return &ImportService{
		logger: logger,
	}
}

// ImportOption 导入选项
type ImportOption struct {
	HasHeader    bool                                 // 是否包含表头
	SheetName    string                               // 工作表名称（Excel）
	StartRow     int                                  // 开始行（从0开始）
	FieldMapping map[string]string                    // 字段映射：列名 -> 字段名
	DataHandler  func([]map[string]interface{}) error // 数据处理器
	ValidateFunc func(map[string]interface{}) error   // 数据验证函数
	BatchSize    int                                  // 批量处理大小
}

// ImportResult 导入结果
type ImportResult struct {
	TotalRows   int           `json:"total_rows"`   // 总行数
	SuccessRows int           `json:"success_rows"` // 成功行数
	FailedRows  int           `json:"failed_rows"`  // 失败行数
	Errors      []ImportError `json:"errors"`       // 错误信息
	ProcessedAt time.Time     `json:"processed_at"` // 处理时间
}

// ImportError 导入错误
type ImportError struct {
	Row    int    `json:"row"`    // 行号
	Column string `json:"column"` // 列名
	Field  string `json:"field"`  // 字段名
	Error  string `json:"error"`  // 错误信息
}

// ImportExcel 导入Excel文件
func (s *ImportService) ImportExcel(ctx context.Context, file *multipart.FileHeader, option ImportOption) (*ImportResult, error) {
	return s.ImportExcelWithContext(ctx, file, option)
}

// ImportExcelWithContext 带上下文的Excel导入
func (s *ImportService) ImportExcelWithContext(ctx context.Context, file *multipart.FileHeader, option ImportOption) (*ImportResult, error) {
	// 打开Excel文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %w", err)
	}
	defer f.Close()

	// 获取工作表
	sheetName := option.SheetName
	if sheetName == "" {
		sheetName = f.GetSheetName(0) // 默认第一个工作表
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %w", err)
	}

	if len(rows) == 0 {
		return nil, errors.New("Excel文件为空")
	}

	// 处理导入
	return s.processImport(ctx, rows, option, "excel")
}

// ImportCSV 导入CSV文件
func (s *ImportService) ImportCSV(ctx context.Context, file *multipart.FileHeader, option ImportOption) (*ImportResult, error) {
	return s.ImportCSVWithContext(ctx, file, option)
}

// ImportCSVWithContext 带上下文的CSV导入
func (s *ImportService) ImportCSVWithContext(ctx context.Context, file *multipart.FileHeader, option ImportOption) (*ImportResult, error) {
	// 打开CSV文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开CSV文件失败: %w", err)
	}
	defer src.Close()

	// 读取CSV数据
	reader := csv.NewReader(src)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取CSV文件失败: %w", err)
	}

	if len(records) == 0 {
		return nil, errors.New("CSV文件为空")
	}

	// 将records转换为 [][]string 格式
	var rows [][]string
	for _, record := range records {
		rows = append(rows, record)
	}

	// 处理导入
	return s.processImport(ctx, rows, option, "csv")
}

// ImportJSON 导入JSON文件
func (s *ImportService) ImportJSON(ctx context.Context, file *multipart.FileHeader, option ImportOption) (*ImportResult, error) {
	return s.ImportJSONWithContext(ctx, file, option)
}

// ImportJSONWithContext 带上下文的JSON导入
func (s *ImportService) ImportJSONWithContext(ctx context.Context, file *multipart.FileHeader, option ImportOption) (*ImportResult, error) {
	// 打开JSON文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 读取文件内容
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 解析JSON
	var records []map[string]interface{}
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	if len(records) == 0 {
		return nil, errors.New("JSON文件为空")
	}

	// 转换为行格式
	rows := make([][]string, 0, len(records)+1)

	// 创建表头
	if option.HasHeader && len(records) > 0 {
		header := make([]string, 0)
		for key := range records[0] {
			header = append(header, key)
		}
		rows = append(rows, header)
	}

	// 添加数据行
	for _, record := range records {
		row := make([]string, 0, len(record))
		for _, key := range rows[0] { // 使用表头的顺序
			if val, exists := record[key]; exists {
				row = append(row, fmt.Sprintf("%v", val))
			} else {
				row = append(row, "")
			}
		}
		rows = append(rows, row)
	}

	// 处理导入
	return s.processImport(ctx, rows, option, "json")
}

// processImport 处理导入数据
func (s *ImportService) processImport(ctx context.Context, rows [][]string, option ImportOption, fileType string) (*ImportResult, error) {
	result := &ImportResult{
		ProcessedAt: time.Now(),
	}

	// 设置默认值
	if option.BatchSize == 0 {
		option.BatchSize = 100
	}

	// 确定开始行
	startRow := option.StartRow
	if option.HasHeader {
		startRow = 1
	}

	// 获取表头映射
	headerMapping := make(map[int]string)
	if option.HasHeader && len(rows) > 0 {
		headers := rows[0]
		for colIndex, header := range headers {
			if fieldName, exists := option.FieldMapping[header]; exists {
				headerMapping[colIndex] = fieldName
			} else {
				// 如果没有映射，直接使用表头作为字段名
				headerMapping[colIndex] = strings.TrimSpace(header)
			}
		}
	}

	// 处理数据行
	var batch []map[string]interface{}
	for rowIndex := startRow; rowIndex < len(rows); rowIndex++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		result.TotalRows++
		row := rows[rowIndex]

		// 构建数据行
		data := make(map[string]interface{})
		for colIndex, value := range row {
			fieldName := ""
			if option.HasHeader {
				// 有表头的情况
				if fieldName = headerMapping[colIndex]; fieldName == "" {
					continue
				}
			} else {
				// 无表头的情况，使用列索引
				fieldName = fmt.Sprintf("column_%d", colIndex)
			}

			// 数据类型转换
			data[fieldName] = convertValue(value)
		}

		// 数据验证
		if option.ValidateFunc != nil {
			if err := option.ValidateFunc(data); err != nil {
				result.FailedRows++
				result.Errors = append(result.Errors, ImportError{
					Row:   rowIndex + 1,
					Error: err.Error(),
				})
				continue
			}
		}

		batch = append(batch, data)

		// 批量处理
		if len(batch) >= option.BatchSize {
			if err := s.processBatch(ctx, batch, option, result); err != nil {
				return result, err
			}
			batch = batch[:0] // 清空batch
		}
	}

	// 处理剩余数据
	if len(batch) > 0 {
		if err := s.processBatch(ctx, batch, option, result); err != nil {
			return result, err
		}
	}

	s.logger.Info(fmt.Sprintf("导入完成: 文件类型=%s, 总行数=%d, 成功=%d, 失败=%d",
		fileType, result.TotalRows, result.SuccessRows, result.FailedRows))

	return result, nil
}

// processBatch 处理批量数据
func (s *ImportService) processBatch(ctx context.Context, batch []map[string]interface{}, option ImportOption, result *ImportResult) error {
	if option.DataHandler == nil {
		return errors.New("数据处理器未设置")
	}

	if err := option.DataHandler(batch); err != nil {
		// 如果批量处理失败，尝试逐条处理
		for _, data := range batch {
			if err := option.DataHandler([]map[string]interface{}{data}); err != nil {
				result.FailedRows++
				result.Errors = append(result.Errors, ImportError{
					Error: fmt.Sprintf("数据处理失败: %v", err),
				})
			} else {
				result.SuccessRows++
			}
		}
	} else {
		result.SuccessRows += len(batch)
	}

	return nil
}

// convertValue 转换值类型
func convertValue(value string) interface{} {
	// 尝试转换为整数
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}

	// 尝试转换为浮点数
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}

	// 尝试转换为布尔值
	if boolVal, err := strconv.ParseBool(value); err == nil {
		return boolVal
	}

	// 返回字符串
	return strings.TrimSpace(value)
}

// toString 将interface{}转换为string
func toString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetExcelSheets 获取Excel文件的工作表列表
func (s *ImportService) GetExcelSheets(file *multipart.FileHeader) ([]string, error) {
	// 打开Excel文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %w", err)
	}
	defer f.Close()

	return f.GetSheetList(), nil
}

// GetExcelColumns 获取Excel文件的列信息
func (s *ImportService) GetExcelColumns(file *multipart.FileHeader, sheetName string, hasHeader bool) ([]string, error) {
	// 打开Excel文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %w", err)
	}
	defer f.Close()

	if sheetName == "" {
		sheetName = f.GetSheetName(0)
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %w", err)
	}

	if len(rows) == 0 {
		return nil, errors.New("工作表为空")
	}

	if hasHeader {
		return rows[0], nil
	}

	// 如果没有表头，返回第一行的数据作为列示例
	if len(rows) > 0 {
		columns := make([]string, len(rows[0]))
		for i := range columns {
			columns[i] = fmt.Sprintf("列%d", i+1)
		}
		return columns, nil
	}

	return nil, errors.New("工作表为空")
}
