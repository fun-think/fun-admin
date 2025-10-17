package admin

// Column 表格列定义
// 提供基础元数据，前端可据此渲染

type Column struct {
	Name      string
	Label     string
	Type      string // text, number, boolean, date, datetime, badge, link
	Sortable  bool
	Align     string // left, center, right
	Visible   *bool
	Width     *int
	Sticky    string
	Formatter string
	EnumMap   map[string]string
	BadgeMap  map[string]string
	UrlField  string
}

func NewColumn(name, label, typ string) *Column {
	return &Column{Name: name, Label: label, Type: typ, Align: "left"}
}

func (c *Column) SetSortable(v bool) *Column              { c.Sortable = v; return c }
func (c *Column) AlignCenter() *Column                    { c.Align = "center"; return c }
func (c *Column) AlignRight() *Column                     { c.Align = "right"; return c }
func (c *Column) SetVisible(v bool) *Column               { c.Visible = &v; return c }
func (c *Column) SetWidth(px int) *Column                 { c.Width = &px; return c }
func (c *Column) SetSticky(pos string) *Column            { c.Sticky = pos; return c }
func (c *Column) SetFormatter(name string) *Column        { c.Formatter = name; return c }
func (c *Column) SetEnumMap(m map[string]string) *Column  { c.EnumMap = m; return c }
func (c *Column) SetBadgeMap(m map[string]string) *Column { c.BadgeMap = m; return c }
func (c *Column) SetUrlField(f string) *Column            { c.UrlField = f; return c }

// Filter 过滤器定义（简单版）

type Filter struct {
	Name    string
	Label   string
	Type    string // text, select, date, daterange, numberrange, boolean
	Options []Option
}

func NewFilter(name, label, typ string) *Filter    { return &Filter{Name: name, Label: label, Type: typ} }
func (f *Filter) SetOptions(opts []Option) *Filter { f.Options = opts; return f }
