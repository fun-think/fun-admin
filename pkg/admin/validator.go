package admin

import (
	"fmt"
	"reflect"
	"strings"
)

// Validator 定义验证器接口
type Validator interface {
	Validate(value interface{}) error
	GetMessage() string
}

// BaseValidator 基础验证器
type BaseValidator struct {
	message string
}

func (v *BaseValidator) GetMessage() string {
	return v.message
}

// RequiredValidator 必填验证器
type RequiredValidator struct {
	BaseValidator
}

func NewRequiredValidator() *RequiredValidator {
	return &RequiredValidator{
		BaseValidator: BaseValidator{
			message: "此字段为必填项",
		},
	}
}

func (v *RequiredValidator) Validate(value interface{}) error {
	if value == nil {
		return fmt.Errorf(v.message)
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.String:
		if strings.TrimSpace(rv.String()) == "" {
			return fmt.Errorf(v.message)
		}
	case reflect.Ptr, reflect.Slice, reflect.Map:
		if rv.IsNil() || rv.Len() == 0 {
			return fmt.Errorf(v.message)
		}
	default:
		// 对于其他类型，只要不是零值就算通过
		if rv.IsZero() {
			return fmt.Errorf(v.message)
		}
	}

	return nil
}

// EmailValidator 邮箱验证器
type EmailValidator struct {
	BaseValidator
}

func NewEmailValidator() *EmailValidator {
	return &EmailValidator{
		BaseValidator: BaseValidator{
			message: "请输入有效的邮箱地址",
		},
	}
}

func (v *EmailValidator) Validate(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("邮箱字段必须是字符串类型")
	}

	if str != "" && !strings.Contains(str, "@") {
		return fmt.Errorf(v.message)
	}

	return nil
}

// MinLengthValidator 最小长度验证器
type MinLengthValidator struct {
	BaseValidator
	minLength int
}

func NewMinLengthValidator(minLength int) *MinLengthValidator {
	return &MinLengthValidator{
		BaseValidator: BaseValidator{
			message: fmt.Sprintf("字段长度不能少于 %d 个字符", minLength),
		},
		minLength: minLength,
	}
}

func (v *MinLengthValidator) Validate(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("该验证器只适用于字符串类型")
	}

	if len(str) < v.minLength {
		return fmt.Errorf(v.message)
	}

	return nil
}

// MaxLengthValidator 最大长度验证器
type MaxLengthValidator struct {
	BaseValidator
	maxLength int
}

func NewMaxLengthValidator(maxLength int) *MaxLengthValidator {
	return &MaxLengthValidator{
		BaseValidator: BaseValidator{
			message: fmt.Sprintf("字段长度不能超过 %d 个字符", maxLength),
		},
		maxLength: maxLength,
	}
}

func (v *MaxLengthValidator) Validate(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("该验证器只适用于字符串类型")
	}

	if len(str) > v.maxLength {
		return fmt.Errorf(v.message)
	}

	return nil
}

// ValidatedField 带验证器的包装字段
// 避免与接口 FieldWithValidators 重名
type ValidatedField struct {
	Field
	Validators []Validator
}

// NewFieldWithValidators 构造包装字段
func NewFieldWithValidators(field Field) *ValidatedField {
	return &ValidatedField{Field: field, Validators: []Validator{}}
}

// AddValidator 添加验证器
func (f *ValidatedField) AddValidator(validator Validator) *ValidatedField {
	f.Validators = append(f.Validators, validator)
	return f
}

// Validate 执行所有验证器
func (f *ValidatedField) Validate(value interface{}) []error {
	var errors []error
	for _, validator := range f.Validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// 验证资源数据
func ValidateResourceData(resource Resource, data map[string]interface{}) map[string][]string {
	errors := make(map[string][]string)

	fields := resource.GetFields()
	for _, field := range fields {
		fieldName := field.GetName()
		value, exists := data[fieldName]

		// 检查必填字段
		if field.IsRequired() && (!exists || value == nil || value == "") {
			errors[fieldName] = append(errors[fieldName], "此字段为必填项")
			continue
		}

		// 如果字段有验证器，执行验证（duck typing）
		if fieldWithValidators, ok := any(field).(interface{ Validate(interface{}) []error }); ok {
			if errs := fieldWithValidators.Validate(value); len(errs) > 0 {
				for _, err := range errs {
					errors[fieldName] = append(errors[fieldName], err.Error())
				}
			}
		}

		// 特殊字段类型验证
		switch field.GetType() {
		case "email":
			if value != nil && value != "" {
				if emailStr, ok := value.(string); ok {
					if !strings.Contains(emailStr, "@") {
						errors[fieldName] = append(errors[fieldName], "请输入有效的邮箱地址")
					}
				}
			}
		}
	}

	return errors
}
