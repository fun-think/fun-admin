package admin

// Field 定义字段接口
type Field interface {
	GetName() string
	GetType() string
	GetLabel() string
	IsRequired() bool
}

// FieldWithValidators 带验证器的字段接口
type FieldWithValidators interface {
	Field
	AddValidator(validator Validator) FieldWithValidators
	Validate(value interface{}) []error
}

// BaseField 是所有字段的基类
type BaseField struct {
	name      string
	fieldType string
	label     string
	required  bool
}

func (f *BaseField) GetName() string {
	return f.name
}

func (f *BaseField) GetType() string {
	return f.fieldType
}

func (f *BaseField) GetLabel() string {
	return f.label
}

func (f *BaseField) IsRequired() bool {
	return f.required
}

// TextField 文本字段
type TextField struct {
	BaseField
	validators []Validator
}

func NewTextField(name string) *TextField {
	return &TextField{
		BaseField: BaseField{
			name:      name,
			fieldType: "text",
			label:     name,
			required:  false,
		},
		validators: []Validator{},
	}
}

func (f *TextField) Label(label string) *TextField {
	f.label = label
	return f
}

func (f *TextField) Required() *TextField {
	f.required = true
	return f
}

func (f *TextField) AddValidator(validator Validator) *TextField {
	f.validators = append(f.validators, validator)
	return f
}

func (f *TextField) Validate(value interface{}) []error {
	var errors []error
	for _, validator := range f.validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// EmailField 邮箱字段
type EmailField struct {
	BaseField
	validators []Validator
}

func NewEmailField(name string) *EmailField {
	emailValidator := NewEmailValidator()

	return &EmailField{
		BaseField: BaseField{
			name:      name,
			fieldType: "email",
			label:     name,
			required:  false,
		},
		validators: []Validator{emailValidator},
	}
}

func (f *EmailField) Label(label string) *EmailField {
	f.label = label
	return f
}

func (f *EmailField) Required() *EmailField {
	f.required = true
	return f
}

func (f *EmailField) AddValidator(validator Validator) *EmailField {
	f.validators = append(f.validators, validator)
	return f
}

func (f *EmailField) Validate(value interface{}) []error {
	var errors []error
	for _, validator := range f.validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// NumberField 数字字段
type NumberField struct {
	BaseField
	validators   []Validator
	defaultValue *int
}

func NewNumberField(name string) *NumberField {
	return &NumberField{
		BaseField: BaseField{
			name:      name,
			fieldType: "number",
			label:     name,
			required:  false,
		},
		validators: []Validator{},
	}
}

func (f *NumberField) Label(label string) *NumberField {
	f.label = label
	return f
}

func (f *NumberField) Required() *NumberField {
	f.required = true
	return f
}

func (f *NumberField) SetDefault(v int) *NumberField {
	f.defaultValue = &v
	return f
}

func (f *NumberField) AddValidator(validator Validator) *NumberField {
	f.validators = append(f.validators, validator)
	return f
}

func (f *NumberField) Validate(value interface{}) []error {
	var errors []error
	for _, validator := range f.validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// SelectField 选择字段
type SelectField struct {
	BaseField
	Options      []Option
	validators   []Validator
	defaultValue *string
}

type Option struct {
	Value string
	Label string
}

func NewSelectField(name string) *SelectField {
	return &SelectField{
		BaseField: BaseField{
			name:      name,
			fieldType: "select",
			label:     name,
			required:  false,
		},
		Options:      []Option{},
		validators:   []Validator{},
		defaultValue: nil,
	}
}

func (f *SelectField) Label(label string) *SelectField {
	f.label = label
	return f
}

func (f *SelectField) Required() *SelectField {
	f.required = true
	return f
}

func (f *SelectField) SetOptions(options []Option) *SelectField {
	f.Options = options
	return f
}

func (f *SelectField) SetDefault(v string) *SelectField {
	f.defaultValue = &v
	return f
}

func (f *SelectField) AddValidator(validator Validator) *SelectField {
	f.validators = append(f.validators, validator)
	return f
}

func (f *SelectField) Validate(value interface{}) []error {
	var errors []error
	for _, validator := range f.validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// TextareaField 文本域字段
type TextareaField struct {
	BaseField
	Rows       int
	validators []Validator
}

func NewTextareaField(name string) *TextareaField {
	return &TextareaField{
		BaseField: BaseField{
			name:      name,
			fieldType: "textarea",
			label:     name,
			required:  false,
		},
		Rows:       4,
		validators: []Validator{},
	}
}

func (f *TextareaField) Label(label string) *TextareaField {
	f.label = label
	return f
}

func (f *TextareaField) Required() *TextareaField {
	f.required = true
	return f
}

func (f *TextareaField) SetRows(rows int) *TextareaField {
	f.Rows = rows
	return f
}

func (f *TextareaField) AddValidator(validator Validator) *TextareaField {
	f.validators = append(f.validators, validator)
	return f
}

func (f *TextareaField) Validate(value interface{}) []error {
	var errors []error
	for _, validator := range f.validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// BooleanField 布尔字段
type BooleanField struct {
	BaseField
	validators   []Validator
	defaultValue *bool
}

func NewBooleanField(name string) *BooleanField {
	return &BooleanField{
		BaseField: BaseField{
			name:      name,
			fieldType: "boolean",
			label:     name,
			required:  false,
		},
		validators: []Validator{},
	}
}

func (f *BooleanField) Label(label string) *BooleanField {
	f.label = label
	return f
}

func (f *BooleanField) Required() *BooleanField {
	f.required = true
	return f
}

func (f *BooleanField) SetDefault(v bool) *BooleanField {
	f.defaultValue = &v
	return f
}

func (f *BooleanField) AddValidator(validator Validator) *BooleanField {
	f.validators = append(f.validators, validator)
	return f
}

func (f *BooleanField) Validate(value interface{}) []error {
	var errors []error
	for _, validator := range f.validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// DateTimeField 日期时间字段
type DateTimeField struct {
	BaseField
	validators []Validator
}

func NewDateTimeField(name string) *DateTimeField {
	return &DateTimeField{
		BaseField: BaseField{
			name:      name,
			fieldType: "datetime",
			label:     name,
			required:  false,
		},
		validators: []Validator{},
	}
}

func (f *DateTimeField) Label(label string) *DateTimeField {
	f.label = label
	return f
}

func (f *DateTimeField) Required() *DateTimeField {
	f.required = true
	return f
}

func (f *DateTimeField) AddValidator(validator Validator) *DateTimeField {
	f.validators = append(f.validators, validator)
	return f
}

func (f *DateTimeField) Validate(value interface{}) []error {
	var errors []error
	for _, validator := range f.validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// DateField 日期字段
type DateField struct {
	BaseField
	validators []Validator
}

func NewDateField(name string) *DateField {
	return &DateField{
		BaseField: BaseField{
			name:      name,
			fieldType: "date",
			label:     name,
			required:  false,
		},
		validators: []Validator{},
	}
}

func (f *DateField) Label(label string) *DateField {
	f.label = label
	return f
}

func (f *DateField) Required() *DateField {
	f.required = true
	return f
}

func (f *DateField) AddValidator(validator Validator) *DateField {
	f.validators = append(f.validators, validator)
	return f
}

func (f *DateField) Validate(value interface{}) []error {
	var errors []error
	for _, validator := range f.validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// RelationshipField 关联字段
type RelationshipField struct {
	BaseField
	RelatedResource string
	DisplayField    string
	validators      []Validator
}

func NewRelationshipField(name string, relatedResource string) *RelationshipField {
	return &RelationshipField{
		BaseField: BaseField{
			name:      name,
			fieldType: "relationship",
			label:     name,
			required:  false,
		},
		RelatedResource: relatedResource,
		DisplayField:    "name",
		validators:      []Validator{},
	}
}

func (f *RelationshipField) Label(label string) *RelationshipField {
	f.label = label
	return f
}

func (f *RelationshipField) Required() *RelationshipField {
	f.required = true
	return f
}

func (f *RelationshipField) SetDisplayField(field string) *RelationshipField {
	f.DisplayField = field
	return f
}

func (f *RelationshipField) AddValidator(validator Validator) *RelationshipField {
	f.validators = append(f.validators, validator)
	return f
}

func (f *RelationshipField) Validate(value interface{}) []error {
	var errors []error
	for _, validator := range f.validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// IDField 内置ID字段
type IDField struct{ BaseField }

func NewIDField() *IDField {
	return &IDField{BaseField: BaseField{name: "id", fieldType: "number", label: "id"}}
}

func (f *IDField) Label(label string) *IDField {
	f.BaseField.label = label
	return f
}

// FileField 文件上传字段
type FileField struct {
	BaseField
	AllowedTypes []string
	MaxSize      int64
}

func NewFileField(name string) *FileField {
	return &FileField{
		BaseField:    BaseField{name: name, fieldType: "file", label: name},
		AllowedTypes: []string{},
		MaxSize:      10 * 1024 * 1024,
	}
}

func (f *FileField) Label(label string) *FileField {
	f.BaseField.label = label
	return f
}

func (f *FileField) SetAllowedTypes(types []string) *FileField {
	f.AllowedTypes = types
	return f
}

func (f *FileField) SetMaxSize(size int64) *FileField {
	f.MaxSize = size
	return f
}
