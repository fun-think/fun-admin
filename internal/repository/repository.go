package repository

import (
	"context"
	"fmt"
	"fun-admin/pkg/database"
	"fun-admin/pkg/logger"
	"fun-admin/pkg/zapgorm2"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const ctxTxKey = "TxKey"

type Repository struct {
	db     *gorm.DB
	e      *casbin.SyncedEnforcer
	dbMgr  *database.Manager
	logger *logger.Logger
}

func NewRepository(
	logger *logger.Logger,
	db *gorm.DB,
	e *casbin.SyncedEnforcer,
) *Repository {
	return &Repository{
		db:     db,
		e:      e,
		dbMgr:  database.NewManager(db),
		logger: logger,
	}
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransaction(r *Repository) Transaction {
	return r
}

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	return r.dbMgr.GetDB(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.dbMgr.WithTransaction(ctx, fn)
}

// WithReadOnlyTransaction 只读事务
func (r *Repository) WithReadOnlyTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.dbMgr.WithReadOnlyTransaction(ctx, fn)
}

// GetQueryBuilder 获取查询构建器
func (r *Repository) GetQueryBuilder(ctx context.Context) *database.QueryBuilder {
	db := r.DB(ctx)
	return database.NewQueryBuilder(db)
}

// Health 检查数据库健康状态
func (r *Repository) Health(ctx context.Context) error {
	return r.dbMgr.Health(ctx)
}

func NewDB(conf *viper.Viper, l *logger.Logger) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	// 设置日志级别
	logLevel := gormlogger.Silent
	if conf.GetBool("data.db.debug") {
		logLevel = gormlogger.Info
	}
	gormLogger := zapgorm2.New(l.Logger).LogMode(logLevel)

	driver := conf.GetString("data.db.user.driver")
	dsn := conf.GetString("data.db.user.dsn")

	config := &gorm.Config{
		Logger: gormLogger,
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 跳过默认事务
		SkipDefaultTransaction: true,
		// 预编译语句缓存
		PrepareStmt: true,
	}

	// GORM doc: https://gorm.io/docs/connecting_to_the_database.html
	switch driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), config)
	case "postgres":
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), config)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), config)
	default:
		panic("unknown db driver")
	}
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	// 获取通用数据库对象sql.DB来设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get sql.DB: %v", err))
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(conf.GetInt("data.db.max_open_conns"))        // 最大打开连接数
	sqlDB.SetMaxIdleConns(conf.GetInt("data.db.max_idle_conns"))        // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(conf.GetDuration("data.db.max_lifetime"))  // 连接最大存活时间
	sqlDB.SetConnMaxIdleTime(conf.GetDuration("data.db.max_idle_time")) // 连接最大空闲时间

	// Additional connection pool settings
	sqlDB.SetMaxIdleConns(conf.GetInt("data.db.pool.max_idle_conns"))
	sqlDB.SetMaxOpenConns(conf.GetInt("data.db.pool.max_open_conns"))
	sqlDB.SetConnMaxLifetime(conf.GetDuration("data.db.pool.max_lifetime"))
	return db
}
func NewCasbinEnforcer(conf *viper.Viper, l *logger.Logger, db *gorm.DB) *casbin.SyncedEnforcer {
	a, _ := gormadapter.NewAdapterByDB(db)
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)

	if err != nil {
		panic(err)
	}
	e, _ := casbin.NewSyncedEnforcer(m, a)

	// 每10秒自动加载策略，防止启动多服务进程策略不一致
	// 如果不想用轮询DB的方式，你也可以使用Casbin Watchers来同步策略，该方式需要基于Redis、Etcd等存储中间件
	// Watchers相关文档：https://casbin.org/zh/docs/watchers
	e.StartAutoLoadPolicy(10 * time.Second)

	// Enable Logger, decide whether to show it in terminal
	//e.EnableLog(true)

	// Save the policy back to DB.
	e.EnableAutoSave(true)

	return e
}
func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}
