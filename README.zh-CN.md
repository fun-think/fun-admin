# fun-admin

<p align="center">
  一个受 Filament 启发的 Go 管理后台框架（Gin + Gorm），并内置 Vue 3 管理后台前端示例。
</p>

<p align="center">
  <a href="README.md"><kbd>English</kbd></a>
  <a href="README.zh-CN.md"><kbd>中文</kbd></a>
</p>

## 这是什么？

`fun-admin` 是一个 Go 语言的后台管理框架 + 可直接运行的示例应用。

如果你喜欢 Filament 的开发体验（Resource / Field / Action / Page），这个项目把类似的理念带到 Go：

- 用 Go 代码定义后台 **Resource**
- 自动获得一致的 CRUD、筛选、列表列、动作等能力
- 可插拔：权限、导航元信息、国际化、导入/导出等

框架代码位于 `pkg/admin`；示例后端在 `cmd/` + `internal/`；前端 UI 在 `web/`。

## 功能特性（Features）

框架能力（Filament-like 组件化能力）：

- **Resource**：`admin.Resource` 描述模型、字段、列、筛选、动作
- **Field**：内置多种字段类型（text/email/number/select/…）
- **Action**：单条/批量动作，支持动作表单与可见性控制
- **Page**：自定义页面（`admin.PageInterface` / `admin.BasePage`）
- **导航元信息**：icon/group/sort/badge，支持在导航中隐藏
- **授权扩展点**：可选接口（create/update/delete 等校验点）
- **字段级权限**：按上下文控制 readable/writable 字段集合
- **API 生成**：根据资源生成 REST API（`pkg/admin/api_generator.go`）
- **国际化**：后端 i18n 资源（`zh-CN`、`en`）

示例应用（本仓库内置）：

- JWT 管理员登录（`/api/admin/login`）
- Casbin 权限中间件（RBAC 风格）
- 操作日志 + 仪表盘统计
- 文件上传、导入/导出接口
- 默认 SQLite 配置（`config/local.yml`）

## 环境要求

- Go `1.23+`（见 `go.mod`）
- Node.js `>= 18.15`（见 `web/package.json`）
- 可选：Docker + Docker Compose（`deploy/docker-compose/` 提供 MySQL/Redis）

## 快速开始

### 后端

1. 运行迁移（建表 + 初始化管理员）：

```bash
go run ./cmd/migration -conf config/local.yml
```

可选参数：

```bash
go run ./cmd/migration -conf config/local.yml -allow-drop=false -admin-password=你的密码
```

2. 启动服务：

```bash
go run ./cmd/server -conf config/local.yml
```

本地地址（以 `config/local.yml` 为例）：

- API：`http://127.0.0.1:8001/api/admin`
- Swagger：`http://127.0.0.1:8001/swagger/index.html`

### 前端

前端位于 `web/`：

```bash
cd web
npm install
npm run dev
```

更多前端说明见 `web/README.md` / `web/README.zh-CN.md`。

## 核心概念

### Resource（资源）

Resource 是框架核心。一个 Resource 通常会描述：

- 管理哪个模型（`GetModel`）
- 可编辑字段（`GetFields`）
- 列表展示列（`GetColumns`）
- 可用筛选（`GetFilters`）
- 可执行动作（`GetActions`）

资源注册入口见 `cmd/bootstrap/bootstrap.go`（`InitAdmin`）。

### Action（动作，含批量）

若资源实现 `admin.ActionExecutor`，即可执行动作：

- 单条或批量（`ids []interface{}`）
- 可携带动作表单参数（`params map[string]interface{}`）
- 可通过可选接口控制动作可见性

## 常用命令

```bash
# 安装开发工具（mockgen、swag）
make init

# docker compose + 迁移 + 启动服务
make bootstrap

# 构建后端二进制 + 构建前端 dist
make build

# 生成 Swagger（输出到 ./docs/）
make swag
```

说明：

- `make test` 目前引用了不存在的 `test/` 目录。
- `make admin-resource` / `make admin-page` 指向了不存在的 `cmd/admin/main.go`；建议直接使用 `cmd/make`（见下文）。

## 生成资源 / 页面脚手架

仓库内置了一个简单脚手架命令：

```bash
go run ./cmd/make -action make:resource -name UserProfile
go run ./cmd/make -action make:page -name Reports
```

生成后，需要在 `cmd/bootstrap/bootstrap.go` 中注册对应的资源/页面。

## 配置说明

- 配置文件为 `config/` 下的 YAML（例如 `config/local.yml`、`config/prod.yml`）
- 选择配置方式：
  - 命令行参数：`-conf config/local.yml`
  - 环境变量：`APP_CONF=config/local.yml`

关键配置项（见 `config/local.yml`）：

- `env`：影响 Gin 模式（`prod/production` => Release，否则 Debug）
- `http.host`、`http.port`
- `data.db.user.driver`、`data.db.user.dsn`
- `data.redis.*`（若 `data.redis.addr` 为空则使用内存缓存）

## 目录结构

- `pkg/admin`：后台框架（资源/字段/动作/页面/API 生成）
- `internal/resources`：示例资源（用户、操作日志、字典、演示 CRUD）
- `internal/handler`：Gin handler（管理端 API）
- `internal/service`：服务层（CRUD/动作执行/校验/导出等）
- `internal/repository`：数据访问层
- `cmd/server`：HTTP 服务入口
- `cmd/migration`：数据库迁移 + 初始化管理员
- `web/`：管理后台前端

## License

见 `LICENSE`。
