package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	// TxKey 事务键
	TxKey ContextKey = "tx"
	// DBKey 数据库键
	DBKey ContextKey = "db"
)

// Manager 数据库管理器
type Manager struct {
	db *gorm.DB
}

// NewManager 创建数据库管理器
func NewManager(db *gorm.DB) *Manager {
	return &Manager{db: db}
}

// GetDB 获取数据库连接
func (m *Manager) GetDB(ctx context.Context) *gorm.DB {
	// 优先从上下文获取事务
	if tx, ok := ctx.Value(TxKey).(*gorm.DB); ok && tx != nil {
		return tx
	}

	// 从上下文获取数据库连接
	if db, ok := ctx.Value(DBKey).(*gorm.DB); ok && db != nil {
		return db
	}

	// 返回默认连接
	return m.db.WithContext(ctx)
}

// WithDB 在上下文中设置数据库连接
func (m *Manager) WithDB(ctx context.Context) context.Context {
	return context.WithValue(ctx, DBKey, m.db.WithContext(ctx))
}

// WithTransaction 执行事务
func (m *Manager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, TxKey, tx)
		return fn(txCtx)
	})
}

// WithReadOnlyTransaction 执行只读事务
func (m *Manager) WithReadOnlyTransaction(ctx context.Context, fn func(context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 设置为只读事务
		tx.Exec("SET TRANSACTION READ ONLY")
		txCtx := context.WithValue(ctx, TxKey, tx)
		return fn(txCtx)
	})
}

// Health 检查数据库健康状态
func (m *Manager) Health(ctx context.Context) error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Stats 获取数据库连接池统计信息
func (m *Manager) Stats() (*sql.DBStats, error) {
	sqlDB, err := m.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	stats := sqlDB.Stats()
	return &stats, nil
}

// Close 关闭数据库连接
func (m *Manager) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// QueryBuilder 查询构建器
type QueryBuilder struct {
	db *gorm.DB
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{db: db}
}

// Paginate 分页查询
func (qb *QueryBuilder) Paginate(page, pageSize int) *QueryBuilder {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	qb.db = qb.db.Offset(offset).Limit(pageSize)
	return qb
}

// OrderBy 排序
func (qb *QueryBuilder) OrderBy(column, direction string) *QueryBuilder {
	if direction != "asc" && direction != "desc" {
		direction = "desc"
	}
	qb.db = qb.db.Order(fmt.Sprintf("%s %s", column, direction))
	return qb
}

// Where 条件查询
func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Where(query, args...)
	return qb
}

// OrWhere OR条件查询
func (qb *QueryBuilder) OrWhere(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Or(query, args...)
	return qb
}

// WhereIn IN查询
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	if len(values) > 0 {
		qb.db = qb.db.Where(fmt.Sprintf("%s IN ?", column), values)
	}
	return qb
}

// WhereBetween BETWEEN查询
func (qb *QueryBuilder) WhereBetween(column string, min, max interface{}) *QueryBuilder {
	qb.db = qb.db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", column), min, max)
	return qb
}

// WhereLike LIKE查询
func (qb *QueryBuilder) WhereLike(column, value string) *QueryBuilder {
	qb.db = qb.db.Where(fmt.Sprintf("%s LIKE ?", column), "%"+value+"%")
	return qb
}

// Join JOIN查询
func (qb *QueryBuilder) Join(table, on string) *QueryBuilder {
	qb.db = qb.db.Joins(fmt.Sprintf("JOIN %s ON %s", table, on))
	return qb
}

// LeftJoin LEFT JOIN查询
func (qb *QueryBuilder) LeftJoin(table, on string) *QueryBuilder {
	qb.db = qb.db.Joins(fmt.Sprintf("LEFT JOIN %s ON %s", table, on))
	return qb
}

// Select 选择字段
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.db = qb.db.Select(columns)
	return qb
}

// Group 分组
func (qb *QueryBuilder) Group(columns ...string) *QueryBuilder {
	qb.db = qb.db.Group(fmt.Sprintf("%s", columns[0]))
	for i := 1; i < len(columns); i++ {
		qb.db = qb.db.Group(columns[i])
	}
	return qb
}

// Having HAVING条件
func (qb *QueryBuilder) Having(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Having(query, args...)
	return qb
}

// GetDB 获取GORM数据库实例
func (qb *QueryBuilder) GetDB() *gorm.DB {
	return qb.db
}

// Find 查询记录
func (qb *QueryBuilder) Find(dest interface{}) error {
	return qb.db.Find(dest).Error
}

// First 查询第一条记录
func (qb *QueryBuilder) First(dest interface{}) error {
	return qb.db.First(dest).Error
}

// Count 统计记录数
func (qb *QueryBuilder) Count(count *int64) error {
	return qb.db.Count(count).Error
}

// Pluck 查询单个字段
func (qb *QueryBuilder) Pluck(column string, dest interface{}) error {
	return qb.db.Pluck(column, dest).Error
}

// Transaction 事务处理器
type Transaction struct {
	db *gorm.DB
}

// NewTransaction 创建事务处理器
func NewTransaction(db *gorm.DB) *Transaction {
	return &Transaction{db: db}
}

// Execute 执行事务
func (t *Transaction) Execute(ctx context.Context, fn func(*gorm.DB) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// ExecuteWithOptions 执行带选项的事务
func (t *Transaction) ExecuteWithOptions(ctx context.Context, opts *sql.TxOptions, fn func(*gorm.DB) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	}, opts)
}

// CustomLogger GORM自定义日志器
type CustomLogger struct {
	logger.Interface
}

// NewCustomLogger 创建自定义日志器
func NewCustomLogger(logLevel logger.LogLevel) *CustomLogger {
	return &CustomLogger{
		Interface: logger.Default.LogMode(logLevel),
	}
}
