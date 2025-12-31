package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExportService 导出服务
type ExportService struct {
}

// NewExportService 创建导出服务
func NewExportService() *ExportService {
	return &ExportService{}
}

// ExportToCSV 导出数据到 CSV 格式
func (s *ExportService) ExportToCSV(data []map[string]interface{}, headers map[string]string) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	defer writer.Flush()

	// 写入表头
	if len(headers) > 0 {
		headerRow := make([]string, 0, len(headers))
		for _, label := range headers {
			headerRow = append(headerRow, label)
		}
		if err := writer.Write(headerRow); err != nil {
			return nil, fmt.Errorf("写入CSV表头失败: %w", err)
		}
	}

	// 写入数据行
	for _, row := range data {
		dataRow := make([]string, 0, len(headers))
		for fieldName := range headers {
			value := row[fieldName]
			if value == nil {
				dataRow = append(dataRow, "")
			} else {
				// 处理布尔值
				if b, ok := value.(bool); ok {
					if b {
						dataRow = append(dataRow, "是")
					} else {
						dataRow = append(dataRow, "否")
					}
				} else {
					dataRow = append(dataRow, fmt.Sprintf("%v", value))
				}
			}
		}
		if err := writer.Write(dataRow); err != nil {
			return nil, fmt.Errorf("写入CSV数据行失败: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// ExportToExcel 导出数据到 Excel 格式
func (s *ExportService) ExportToExcel(data []map[string]interface{}, headers map[string]string) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			// 关闭Excel文件失败，忽略错误
		}
	}()

	// 创建工作表
	sheetName := "数据导出"
	f.SetSheetName("Sheet1", sheetName)

	// 写入表头
	colIndex := 0
	for _, label := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIndex+1, 1)
		f.SetCellValue(sheetName, cell, label)
		colIndex++
	}

	// 写入数据行
	for rowIndex, row := range data {
		colIndex = 0
		for fieldName := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			value := row[fieldName]
			if value == nil {
				f.SetCellValue(sheetName, cell, "")
			} else {
				// 处理布尔值
				if b, ok := value.(bool); ok {
					if b {
						f.SetCellValue(sheetName, cell, "是")
					} else {
						f.SetCellValue(sheetName, cell, "否")
					}
				} else {
					// 处理数字
					if fValue, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64); err == nil {
						f.SetCellValue(sheetName, cell, fValue)
					} else {
						f.SetCellValue(sheetName, cell, value)
					}
				}
			}
			colIndex++
		}
	}

	// 自动调整列宽
	for i := 0; i < len(headers); i++ {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, colName, colName, 20)
	}

	// 保存到字节流
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("保存Excel文件失败: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateFileName 生成导出文件名
func (s *ExportService) GenerateFileName(prefix string, format string) string {
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s.%s", prefix, timestamp, format)
}
