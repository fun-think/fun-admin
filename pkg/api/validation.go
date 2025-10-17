package api

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitCustomValidators 初始化自定义验证器
func InitCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义标签名获取函数
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// 注册自定义验证规则
		v.RegisterValidation("mobile", validateMobile)
		v.RegisterValidation("username", validateUsername)
		v.RegisterValidation("password", validatePassword)
		v.RegisterValidation("chinese", validateChinese)
		v.RegisterValidation("id_card", validateIDCard)
	}
}

// validateMobile 验证手机号码
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	if len(mobile) != 11 {
		return false
	}

	// 简单的手机号验证规则
	if mobile[0] != '1' {
		return false
	}

	for _, char := range mobile {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	// 用户名只能包含字母、数字、下划线
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return false
		}
	}

	// 不能以数字开头
	if username[0] >= '0' && username[0] <= '9' {
		return false
	}

	return true
}

// validatePassword 验证密码强度
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 || len(password) > 32 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	// 至少包含大写字母、小写字母、数字中的两种
	validTypes := 0
	if hasUpper {
		validTypes++
	}
	if hasLower {
		validTypes++
	}
	if hasNumber {
		validTypes++
	}
	if hasSpecial {
		validTypes++
	}

	return validTypes >= 2
}

// validateChinese 验证中文字符
func validateChinese(fl validator.FieldLevel) bool {
	text := fl.Field().String()
	for _, char := range text {
		if char < '\u4e00' || char > '\u9fa5' {
			return false
		}
	}
	return true
}

// validateIDCard 验证身份证号码
func validateIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	if len(idCard) != 18 {
		return false
	}

	// 前17位必须是数字
	for i := 0; i < 17; i++ {
		if idCard[i] < '0' || idCard[i] > '9' {
			return false
		}
	}

	// 最后一位可以是数字或X
	lastChar := idCard[17]
	if !((lastChar >= '0' && lastChar <= '9') || lastChar == 'X' || lastChar == 'x') {
		return false
	}

	return true
}

// ValidationRule 验证规则结构
type ValidationRule struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// GetValidationRules 获取常用验证规则
func GetValidationRules() map[string][]ValidationRule {
	return map[string][]ValidationRule{
		"user": {
			{Field: "username", Rule: "required,username", Message: "用户名必填且格式正确"},
			{Field: "email", Rule: "required,email", Message: "邮箱必填且格式正确"},
			{Field: "password", Rule: "required,password", Message: "密码必填且强度足够"},
			{Field: "mobile", Rule: "omitempty,mobile", Message: "手机号格式不正确"},
		},
		"admin": {
			{Field: "username", Rule: "required,min=3,max=20", Message: "用户名必填，长度3-20字符"},
			{Field: "email", Rule: "required,email", Message: "邮箱必填且格式正确"},
			{Field: "password", Rule: "required,min=8,max=32", Message: "密码必填，长度8-32字符"},
			{Field: "role", Rule: "required", Message: "角色必填"},
		},
		"profile": {
			{Field: "nickname", Rule: "required,min=2,max=20", Message: "昵称必填，长度2-20字符"},
			{Field: "avatar", Rule: "omitempty,url", Message: "头像必须是有效的URL"},
			{Field: "mobile", Rule: "omitempty,mobile", Message: "手机号格式不正确"},
		},
	}
}
