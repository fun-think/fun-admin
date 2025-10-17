package docs

import (
	"encoding/json"
	"fmt"
	"fun-admin/internal/handler"
	"fun-admin/internal/model"
	"fun-admin/internal/service"
	"fun-admin/pkg/admin"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// DocumentGenerator API文档生成器
type DocumentGenerator struct {
	swaggerDoc *SwaggerDocument
}

// SwaggerDocument Swagger文档结构
type SwaggerDocument struct {
	Swagger             string                 `json:"swagger"`
	Info                Info                   `json:"info"`
	Host                string                 `json:"host"`
	BasePath            string                 `json:"basePath,omitempty"`
	Paths               map[string]interface{} `json:"paths"`
	Definitions         map[string]interface{} `json:"definitions"`
	SecurityDefinitions map[string]interface{} `json:"securityDefinitions"`
}

// Info API信息
type Info struct {
	Description    string  `json:"description"`
	Title          string  `json:"title"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
	Version        string  `json:"version"`
}

// Contact 联系信息
type Contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

// License 许可证信息
type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// NewDocumentGenerator 创建文档生成器
func NewDocumentGenerator() *DocumentGenerator {
	return &DocumentGenerator{
		swaggerDoc: &SwaggerDocument{
			Swagger: "2.0",
			Info: Info{
				Description:    "Fun-Admin 现代化管理后台 API",
				Title:          "Fun-Admin API",
				TermsOfService: "https://github.com/your-username/fun-admin",
				Contact: Contact{
					Name:  "Fun-Admin Team",
					URL:   "https://github.com/your-username/fun-admin",
					Email: "support@fun-admin.com",
				},
				License: License{
					Name: "MIT",
					URL:  "https://opensource.org/licenses/MIT",
				},
				Version: "1.0.0",
			},
			Host:        "localhost:8000",
			BasePath:    "/api/v1",
			Paths:       make(map[string]interface{}),
			Definitions: make(map[string]interface{}),
			SecurityDefinitions: map[string]interface{}{
				"Bearer": map[string]interface{}{
					"type": "apiKey",
					"name": "Authorization",
					"in":   "header",
				},
			},
		},
	}
}

// Generate 生成API文档
func (dg *DocumentGenerator) Generate(handlers map[string]interface{}, resources map[string]admin.Resource) error {
	// 生成基础模块文档
	dg.generateBasicHandlers()

	// 生成资源模块文档
	dg.generateResourceHandlers(resources)

	// 生成数据字典文档
	dg.generateDictHandlers()

	// 生成配置管理文档
	dg.generateConfigHandlers()

	// 生成文件管理文档
	dg.generateFileHandlers()

	// 生成数据导入文档
	dg.generateImportHandlers()

	// 生成通用模型定义
	dg.generateCommonDefinitions()

	return nil
}

// generateBasicHandlers 生成基础处理器文档
func (dg *DocumentGenerator) generateBasicHandlers() {
	// 用户管理
	dg.addHandlerPath("/users", "get", "用户模块", "获取用户列表", []Parameter{
		{Name: "page", Type: "integer", Description: "页码", Required: true},
		{Name: "page_size", Type: "integer", Description: "每页数量", Required: true},
		{Name: "username", Type: "string", Description: "用户名"},
		{Name: "nickname", Type: "string", Description: "昵称"},
	})

	dg.addHandlerPath("/users/{id}", "get", "用户模块", "获取用户详情", []Parameter{
		{Name: "id", Type: "integer", Description: "用户ID", Required: true, In: "path"},
	})

	// 角色管理
	dg.addHandlerPath("/roles", "get", "角色模块", "获取角色列表", []Parameter{
		{Name: "page", Type: "integer", Description: "页码", Required: true},
		{Name: "page_size", Type: "integer", Description: "每页数量", Required: true},
		{Name: "name", Type: "string", Description: "角色名称"},
	})

	// 菜单管理
	dg.addHandlerPath("/menus", "get", "菜单模块", "获取菜单列表", []Parameter{
		{Name: "type", Type: "string", Description: "菜单类型"},
	})

	// API管理
	dg.addHandlerPath("/apis", "get", "API模块", "获取API列表", []Parameter{
		{Name: "page", Type: "integer", Description: "页码", Required: true},
		{Name: "page_size", Type: "integer", Description: "每页数量", Required: true},
		{Name: "group", Type: "string", Description: "API分组"},
	})
}

// generateResourceHandlers 生成资源处理器文档
func (dg *DocumentGenerator) generateResourceHandlers(resources map[string]admin.Resource) {
	for slug, resource := range resources {
		// 资源列表
		dg.addHandlerPath("/"+slug, "get", resource.GetTitle(), "获取"+resource.GetTitle()+"列表", []Parameter{
			{Name: "page", Type: "integer", Description: "页码", Required: true},
			{Name: "page_size", Type: "integer", Description: "每页数量", Required: true},
		})

		// 资源详情
		dg.addHandlerPath("/"+slug+"/{id}", "get", resource.GetTitle(), "获取"+resource.GetTitle()+"详情", []Parameter{
			{Name: "id", Type: "integer", Description: "ID", Required: true, In: "path"},
		})

		// 创建资源
		dg.addHandlerPath("/"+slug, "post", resource.GetTitle(), "创建"+resource.GetTitle(), []Parameter{
			{Name: "data", Type: "object", Description: "资源数据", Required: true, In: "body", Schema: "#/definitions/" + dg.getModelName(resource.GetModel())},
		})

		// 更新资源
		dg.addHandlerPath("/"+slug+"/{id}", "put", resource.GetTitle(), "更新"+resource.GetTitle(), []Parameter{
			{Name: "id", Type: "integer", Description: "ID", Required: true, In: "path"},
			{Name: "data", Type: "object", Description: "资源数据", Required: true, In: "body", Schema: "#/definitions/" + dg.getModelName(resource.GetModel())},
		})

		// 删除资源
		dg.addHandlerPath("/"+slug+"/{id}", "delete", resource.GetTitle(), "删除"+resource.GetTitle(), []Parameter{
			{Name: "id", Type: "integer", Description: "ID", Required: true, In: "path"},
		})

		// 导出资源
		dg.addHandlerPath("/"+slug+"/export", "get", resource.GetTitle(), "导出"+resource.GetTitle(), []Parameter{
			{Name: "format", Type: "string", Description: "导出格式 (excel/csv)"},
		})
	}
}

// generateDictHandlers 生成字典处理器文档
func (dg *DocumentGenerator) generateDictHandlers() {
	// 字典类型
	dg.addHandlerPath("/dict/types", "get", "字典管理", "获取字典类型列表", []Parameter{
		{Name: "page", Type: "integer", Description: "页码", Required: true},
		{Name: "page_size", Type: "integer", Description: "每页数量", Required: true},
		{Name: "name", Type: "string", Description: "字典类型名称"},
	})

	dg.addHandlerPath("/dict/types", "post", "字典管理", "创建字典类型", []Parameter{
		{Name: "data", Type: "object", Description: "字典类型数据", Required: true, In: "body"},
	})

	dg.addHandlerPath("/dict/types/{id}", "put", "字典管理", "更新字典类型", []Parameter{
		{Name: "id", Type: "integer", Description: "字典类型ID", Required: true, In: "path"},
		{Name: "data", Type: "object", Description: "字典类型数据", Required: true, In: "body"},
	})

	dg.addHandlerPath("/dict/types/{id}", "delete", "字典管理", "删除字典类型", []Parameter{
		{Name: "id", Type: "integer", Description: "字典类型ID", Required: true, In: "path"},
	})

	// 字典数据
	dg.addHandlerPath("/dict/data", "get", "字典管理", "获取字典数据列表", []Parameter{
		{Name: "type_id", Type: "integer", Description: "字典类型ID", Required: true},
		{Name: "label", Type: "string", Description: "字典标签"},
		{Name: "status", Type: "integer", Description: "状态"},
	})

	dg.addHandlerPath("/dict/data", "post", "字典管理", "创建字典数据", []Parameter{
		{Name: "data", Type: "object", Description: "字典数据", Required: true, In: "body"},
	})

	dg.addHandlerPath("/dict/data/{id}", "put", "字典管理", "更新字典数据", []Parameter{
		{Name: "id", Type: "integer", Description: "字典数据ID", Required: true, In: "path"},
		{Name: "data", Type: "object", Description: "字典数据", Required: true, In: "body"},
	})

	dg.addHandlerPath("/dict/data/{id}", "delete", "字典管理", "删除字典数据", []Parameter{
		{Name: "id", Type: "integer", Description: "字典数据ID", Required: true, In: "path"},
	})

	dg.addHandlerPath("/dict/code/{code}", "get", "字典管理", "根据编码获取字典数据", []Parameter{
		{Name: "code", Type: "string", Description: "字典类型编码", Required: true, In: "path"},
	})
}

// generateConfigHandlers 生成配置处理器文档
func (dg *DocumentGenerator) generateConfigHandlers() {
	// 系统设置
	dg.addHandlerPath("/config/settings/{key}", "get", "配置管理", "获取系统设置", []Parameter{
		{Name: "key", Type: "string", Description: "设置键", Required: true, In: "path"},
	})

	dg.addHandlerPath("/config/settings", "post", "配置管理", "设置系统设置", []Parameter{
		{Name: "data", Type: "object", Description: "设置数据", Required: true, In: "body"},
	})

	dg.addHandlerPath("/config/settings/group/{group}", "get", "配置管理", "根据分组获取系统设置", []Parameter{
		{Name: "group", Type: "string", Description: "分组名称", Required: true, In: "path"},
	})

	// 邮件配置
	dg.addHandlerPath("/config/email", "get", "配置管理", "获取邮件配置", []Parameter{})

	dg.addHandlerPath("/config/email", "post", "配置管理", "设置邮件配置", []Parameter{
		{Name: "config", Type: "object", Description: "邮件配置", Required: true, In: "body"},
	})

	dg.addHandlerPath("/config/email/test", "post", "配置管理", "测试邮件配置", []Parameter{
		{Name: "config", Type: "object", Description: "邮件配置", Required: true, In: "body"},
		{Name: "test_email", Type: "string", Description: "测试邮箱", Required: true, In: "body"},
	})

	// 存储配置
	dg.addHandlerPath("/config/storage", "get", "配置管理", "获取存储配置", []Parameter{})

	dg.addHandlerPath("/config/storage", "post", "配置管理", "设置存储配置", []Parameter{
		{Name: "config", Type: "object", Description: "存储配置", Required: true, In: "body"},
	})

	dg.addHandlerPath("/config/storage/test", "post", "配置管理", "测试存储配置", []Parameter{
		{Name: "config", Type: "object", Description: "存储配置", Required: true, In: "body"},
	})
}

// generateFileHandlers 生成文件管理文档
func (dg *DocumentGenerator) generateFileHandlers() {
	// 文件上传
	dg.addHandlerPath("/files/upload", "post", "文件管理", "上传文件", []Parameter{
		{Name: "file", Type: "file", Description: "上传文件", Required: true, In: "formData"},
		{Name: "type", Type: "string", Description: "文件类型", In: "formData"},
	})

	// 文件删除
	dg.addHandlerPath("/files/{key}", "delete", "文件管理", "删除文件", []Parameter{
		{Name: "key", Type: "string", Description: "文件键", Required: true, In: "path"},
	})

	// 获取文件URL
	dg.addHandlerPath("/files/{key}/url", "get", "文件管理", "获取文件访问URL", []Parameter{
		{Name: "key", Type: "string", Description: "文件键", Required: true, In: "path"},
		{Name: "expire", Type: "integer", Description: "过期时间(秒)", In: "query"},
	})
}

// generateImportHandlers 生成数据导入文档
func (dg *DocumentGenerator) generateImportHandlers() {
	// 数据导入
	dg.addHandlerPath("/import/data", "post", "数据导入", "导入数据", []Parameter{
		{Name: "file", Type: "file", Description: "导入文件", Required: true, In: "formData"},
		{Name: "type", Type: "string", Description: "文件类型 (excel/csv/json)", Required: true, In: "formData"},
		{Name: "resource", Type: "string", Description: "资源名称", Required: true, In: "formData"},
		{Name: "has_header", Type: "boolean", Description: "是否包含表头", In: "formData"},
		{Name: "sheet_name", Type: "string", Description: "工作表名称", In: "formData"},
		{Name: "start_row", Type: "integer", Description: "开始行", In: "formData"},
	})

	// 获取Excel工作表
	dg.addHandlerPath("/import/excel/sheets", "post", "数据导入", "获取Excel工作表列表", []Parameter{
		{Name: "file", Type: "file", Description: "Excel文件", Required: true, In: "formData"},
	})

	// 获取Excel列信息
	dg.addHandlerPath("/import/excel/columns", "post", "数据导入", "获取Excel列信息", []Parameter{
		{Name: "file", Type: "file", Description: "Excel文件", Required: true, In: "formData"},
		{Name: "sheet_name", Type: "string", Description: "工作表名称", Required: true, In: "formData"},
		{Name: "has_header", Type: "boolean", Description: "是否包含表头", In: "formData"},
	})
}

// generateCommonDefinitions 生成通用模型定义
func (dg *DocumentGenerator) generateCommonDefinitions() {
	// 响应结构
	dg.swaggerDoc.Definitions["Response"] = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"code": map[string]interface{}{
				"type":        "integer",
				"description": "响应码",
			},
			"message": map[string]interface{}{
				"type":        "string",
				"description": "响应消息",
			},
			"data": map[string]interface{}{
				"description": "响应数据",
			},
		},
	}

	// 分页响应
	dg.swaggerDoc.Definitions["PageResponse"] = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"code": map[string]interface{}{
				"type":        "integer",
				"description": "响应码",
			},
			"message": map[string]interface{}{
				"type":        "string",
				"description": "响应消息",
			},
			"data": map[string]interface{}{
				"$ref": "#/definitions/PageData",
			},
		},
	}

	// 分页数据
	dg.swaggerDoc.Definitions["PageData"] = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"items": map[string]interface{}{
				"type":        "array",
				"description": "数据列表",
			},
			"total": map[string]interface{}{
				"type":        "integer",
				"description": "总数量",
			},
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "当前页",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "每页数量",
			},
			"total_pages": map[string]interface{}{
				"type":        "integer",
				"description": "总页数",
			},
		},
	}

	// 文件信息
	dg.swaggerDoc.Definitions["FileInfo"] = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "文件名",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "文件路径",
			},
			"url": map[string]interface{}{
				"type":        "string",
				"description": "访问URL",
			},
			"size": map[string]interface{}{
				"type":        "integer",
				"description": "文件大小",
			},
			"ext": map[string]interface{}{
				"type":        "string",
				"description": "文件扩展名",
			},
			"storage_type": map[string]interface{}{
				"type":        "string",
				"description": "存储类型",
			},
			"content_type": map[string]interface{}{
				"type":        "string",
				"description": "内容类型",
			},
			"etag": map[string]interface{}{
				"type":        "string",
				"description": "文件标识",
			},
		},
	}

	// 导入结果
	dg.swaggerDoc.Definitions["ImportResult"] = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"total_rows": map[string]interface{}{
				"type":        "integer",
				"description": "总行数",
			},
			"success_rows": map[string]interface{}{
				"type":        "integer",
				"description": "成功行数",
			},
			"failed_rows": map[string]interface{}{
				"type":        "integer",
				"description": "失败行数",
			},
			"errors": map[string]interface{}{
				"type":        "array",
				"description": "错误信息",
			},
			"processed_at": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"description": "处理时间",
			},
		},
	}

	// 从 model 包中直接引用核心业务模型
	dg.addModelDefinition("User", &model.SysUser{})
	dg.addModelDefinition("Role", &model.SysRole{})
	dg.addModelDefinition("Menu", &model.SysMenu{})
	dg.addModelDefinition("Api", &model.SysApi{})
	dg.addModelDefinition("OperationLog", &model.SysOperationLog{})
	dg.addModelDefinition("DictionaryType", &model.AdminDictionaryType{})
	dg.addModelDefinition("DictionaryData", &model.AdminDictionaryData{})
	dg.addModelDefinition("Config", &model.AdminSetting{})
}

// addModelDefinition 根据 Go 结构体自动生成 Swagger 模型定义
func (dg *DocumentGenerator) addModelDefinition(name string, model interface{}) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	properties := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		// 解析 json tag, 获取字段名 (忽略 ,omitempty 等选项)
		fieldName := strings.Split(jsonTag, ",")[0]

		property := map[string]interface{}{
			"description": field.Tag.Get("comment"),
		}

		switch field.Type.Kind() {
		case reflect.String:
			property["type"] = "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			property["type"] = "integer"
		case reflect.Float32, reflect.Float64:
			property["type"] = "number"
		case reflect.Bool:
			property["type"] = "boolean"
		default:
			property["type"] = "string"
		}

		properties[fieldName] = property
	}

	dg.swaggerDoc.Definitions[name] = map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}
}

// Parameter API参数
type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	In          string `json:"in,omitempty"`
	Schema      string `json:"schema,omitempty"`
}

// addHandlerPath 添加处理器路径
func (dg *DocumentGenerator) addHandlerPath(path, method, tag, summary string, params []Parameter) {
	fullPath := dg.swaggerDoc.BasePath + path

	if dg.swaggerDoc.Paths[fullPath] == nil {
		dg.swaggerDoc.Paths[fullPath] = make(map[string]interface{})
	}

	pathItem := dg.swaggerDoc.Paths[fullPath].(map[string]interface{})

	operation := map[string]interface{}{
		"tags":        []string{tag},
		"summary":     summary,
		"description": summary,
		"consumes":    []string{"application/json"},
		"produces":    []string{"application/json"},
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"description": "OK",
				"schema": map[string]interface{}{
					"$ref": "#/definitions/Response",
				},
			},
		},
	}

	// 添加参数
	if len(params) > 0 {
		var parameters []map[string]interface{}
		for _, param := range params {
			p := map[string]interface{}{
				"name":        param.Name,
				"type":        param.Type,
				"description": param.Description,
				"required":    param.Required,
			}
			if param.In != "" {
				p["in"] = param.In
			}
			if param.Schema != "" {
				p["schema"] = map[string]interface{}{
					"$ref": param.Schema,
				}
			}
			parameters = append(parameters, p)
		}
		operation["parameters"] = parameters
	}

	// 添加安全认证（除了登录接口）
	if path != "/login" {
		operation["security"] = []map[string][]string{
			{"Bearer": {}},
		}
	}

	pathItem[method] = operation
}

// getModelName 获取模型名称
func (dg *DocumentGenerator) getModelName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 获取类型名称
	typeName := t.Name()

	// 处理包名
	if idx := strings.LastIndex(typeName, "."); idx != -1 {
		typeName = typeName[idx+1:]
	}

	return typeName
}

// SaveToFile 保存到文件
func (dg *DocumentGenerator) SaveToFile(filename string) error {
	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(dg.swaggerDoc, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}

	// 写入文件
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// GenerateSwagger 生成Swagger文档
func GenerateSwagger(outputFile string) error {
	generator := NewDocumentGenerator()

	// 这里可以传入实际的handlers和resources
	// 目前使用空的map生成基础文档
	handlers := make(map[string]interface{})
	resources := make(map[string]admin.Resource)

	if err := generator.Generate(handlers, resources); err != nil {
		return fmt.Errorf("生成API文档失败: %w", err)
	}

	if err := generator.SaveToFile(outputFile); err != nil {
		return fmt.Errorf("保存API文档失败: %w", err)
	}

	return nil
}
