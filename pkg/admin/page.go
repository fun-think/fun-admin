package admin

// PageInterface 自定义页面接口
type PageInterface interface {
	// GetTitle 获取页面标题
	GetTitle() string

	// GetSlug 获取页面标识符
	GetSlug() string

	// GetPath 获取页面路径
	GetPath() string

	// GetIcon 获取页面图标
	GetIcon() string

	// IsVisible 是否在导航中显示
	IsVisible() bool

	// GetPermissions 获取页面所需权限
	GetPermissions() []string
}

// BasePage 自定义页面基类
type BasePage struct {
	Title       string
	Slug        string
	Path        string
	Icon        string
	Visible     bool
	Permissions []string
}

// GetTitle 获取页面标题
func (p *BasePage) GetTitle() string {
	return p.Title
}

// GetSlug 获取页面标识符
func (p *BasePage) GetSlug() string {
	return p.Slug
}

// GetPath 获取页面路径
func (p *BasePage) GetPath() string {
	return p.Path
}

// GetIcon 获取页面图标
func (p *BasePage) GetIcon() string {
	return p.Icon
}

// IsVisible 是否在导航中显示
func (p *BasePage) IsVisible() bool {
	return p.Visible
}

// GetPermissions 获取页面所需权限
func (p *BasePage) GetPermissions() []string {
	return p.Permissions
}

// NewBasePage 创建基础页面
func NewBasePage(title, slug, path string) *BasePage {
	return &BasePage{
		Title:   title,
		Slug:    slug,
		Path:    path,
		Visible: true,
	}
}

// SetIcon 设置页面图标
func (p *BasePage) SetIcon(icon string) *BasePage {
	p.Icon = icon
	return p
}

// SetVisible 设置是否可见
func (p *BasePage) SetVisible(visible bool) *BasePage {
	p.Visible = visible
	return p
}

// SetPermissions 设置页面权限
func (p *BasePage) SetPermissions(permissions ...string) *BasePage {
	p.Permissions = permissions
	return p
}
