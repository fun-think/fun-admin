# N-Admin

N-Admin 是一个基于 [go-nunu](https://github.com/go-nunu/nunu) 开发的开源管理后台模板，采用 **Gin + Casbin（RBAC）+ Vue3 + AntDesignVue + AntdvPro** 技术栈，提供快速开发的基础架构。

<p align="center"><img src="https://github.com/go-nunu/fun-admin/blob/main/web/src/assets/images/preview-home.png?raw=true"></p>
<p align="center"><img src="https://github.com/go-nunu/fun-admin/blob/main/web/src/assets/images/preview-api.png?raw=true"></p>

## 要求
要运行项目，您需要在系统上安装以下软件：

* Git
* Golang 1.23 或 更高版本
* NodeJS 18 或 更高版本

## 快速开始

```
# 1. clone项目
git clone https://github.com/go-nunu/fun-admin.git

# 2. 启动项目
cd fun-admin
go run cmd/server/main.go

# 3. 访问项目
浏览器访问：http://localhost:8000


超管账号：admin
超管密码：123456

普通用户账号：user
普通用户密码：123456
```
## 📚 角色权限操作流程
当添加API接口或菜单时，需要手动添加权限策略。

1. 添加API接口（操作路径：权限模块->接口管理->添加API）
2. 添加前端菜单（操作路径：权限模块->菜单管理->添加菜单）
3. 添加权限策略（操作路径：权限模块->角色管理->添加角色/分配权限）


## 📌 功能特性
- ✅**权限管理**：基于 Casbin 实现 RBAC 角色权限控制，权限粒度支持接口和菜单控制。
- ✅**多数据库支持**：支持 MySQL、Postgres、Sqlite 等数据库。
- ✅**管理员管理**：支持管理员账号增删改查，密码加密存储。
- ✅**JWT 认证**：支持 Token 认证，提供登录、登出功能。
- ✅**前后端分离**：RESTful API 设计，支持前后端独立部署。
- ✅**支持一键打包**：整站打包为一个可执行二进制文件。
- ✅**防呆设计**：超管账号始终拥有所有菜单及API权限，防止误操作。


## 🚀 技术栈

### 后端技术栈
- **[go-nunu](https://github.com/go-nunu/nunu)** - 轻量级 Golang 脚手架
- **[Gin](https://github.com/gin-gonic/gin)** - 轻量级 Web 框架
- **[Casbin](https://github.com/casbin/casbin)** - 权限管理（RBAC）
- **[GORM](https://github.com/go-gorm/gorm)** - Golang ORM 框架
- **JWT** - 认证和授权
- **MySQL/Postgres/Sqlite** - 数据库支持

### 前端技术栈
- **[AntdvPro](https://github.com/antdv-pro/antdv-pro)** - 企业级中后台前端/设计解决方案
- **[Vue3](https://github.com/vuejs/)** - 渐进式 JavaScript 框架
- **[Vite](https://github.com/vitejs/vite)** - 极速构建工具



## 📦 安装与运行

### 1️⃣ 克隆项目
```bash
git clone https://github.com/go-nunu/fun-admin.git
cd fun-admin
```

### 2️⃣ 后端启动
#### 配置修改
编辑 `config/local.yml` 并修改必要的配置信息。

#### 执行数据迁移，初始化项目数据（仅项目首次启动时执行）
```bash
go run cmd/migration/main.go
```

#### 运行后端服务
```bash
go run cmd/server/main.go
```

或者使用 `nunu run` 进行热加载开发：
```bash
nunu run
```

### 3️⃣ 前端启动
```bash
cd web
npm install
npm run dev
```

### 4️⃣ 访问地址
后端服务运行在 `http://localhost:8000`

前端服务运行在 `http://localhost:6678`


## 🔑 权限管理（RBAC）
本项目使用 **Casbin** 进行角色权限管理。

- **角色**：管理员、普通用户等
- **权限**：增删改查等操作权限
- **模型**：采用 `RBAC` 访问控制模型
- **存储**：权限策略存储于数据库

示例策略：

API接口
```
p, admin, api:/api/user, GET
p, admin, api:/api/user, POST
p, admin, api:/api/user, PUT
p, admin, api:/api/user, DELETE
p, user, api:/api/profile, GET
```
前端菜单
```
p, admin, menu:/users, read
p, user, menu:/admin/roles, read
```

## 📦 打包部署
```
cd web
npm run build

cd ../
go build -o server cmd/server/main.go
./server

访问：http://127.0.0.1:8000/
```


注意：上面的方法会将服务端和前端的静态资源一起打包到可执行二进制程序中。

如果不需要Golang来渲染前端，可以采用Nginx等反向代理工具将前端静态资源部署到Nginx上。


## 📜 许可证
本项目基于 **MIT License** 开源，欢迎贡献！

我先全局检索核心能力（权限/资源/通知/全局搜索/多租户），再分批检索上传与媒体、仪表盘与组件、表格批量操作与筛选、关系管理、主题与暗色、审计日志、任务调度、表单系统、本地化与模拟登录。接下来开始并行搜索。
Search files...
Search files...
Search files...
Search files...
Search files...

我先做一次并行检索，确认是否已有：软删除/恢复、多租户、全局搜索、实时通知（WS/SSE）、导入（CSV/Excel）。然后基于结果给出缺失与改进建议。
Ran tool
Ran tool
Ran tool
Ran tool
Ran tool

我将并行读取几个核心文件，确认资源/表格/表单/页面能力，以及是否支持 WebSocket 等特性，然后给出与 Filament v4 的差距分析与改进建议。
Ran tool
Ran tool
Ran tool
Ran tool
Ran tool

我将并行读取资源服务、动作定义、页面/管理器、验证器和通知服务，全面对照 Filament v4 的资源/表单/表格/动作/通知能力。
Ran tool
Ran tool
Ran tool
Ran tool
Ran tool
Ran tool
Ran tool

我先查看权限与钩子接口定义，以对照 Filament 的授权与生命周期能力。随后并行打开资源定义示例与表格/过滤能力更多线索。
Ran tool


### 资源与导航
- 完善
  - 资源导航元信息：图标、分组、排序、徽章数、是否显示等
  - 资源页面元信息：默认排序/可排序字段白名单、可过滤字段白名单、可搜索字段白名单
  - 资源层面的多语言能力（目前仅在 `listResourcesHandler` 做了简单 title 映射）

### 表格（Table）
- 完善（Filament 表格常见能力）
  - 批量动作（除批量删除外）：批量导出、批量更新、批量自定义动作
  - 高级列：可见性切换、固定列、列宽、格式化器、枚举/徽章、图片/头像、链接/操作列
  - 汇总与统计、分组、树形/可重排、内联编辑
  - 筛选器表单、保存筛选器、筛选器徽章（chips）、快速搜索占位
  - 分页大小切换、密度（紧凑/宽松）切换
  - 安全：orderBy/filters/search 字段白名单校验

### 表单（Form）
- 完善（Filament 表单常见能力）
  - 高级字段：富文本、Markdown、代码编辑、颜色、级联/联动、Tags、多选、Radio/CheckboxGroup、Password（带确认）、Repeater/嵌套表单、KeyValue/Wizard 步骤
  - 表单容器/布局：Grid/Section/Tabs/Card/Divider/HelpText
  - 字段依赖与条件显示、禁用/只读回调
  - 唯一性/正则/数值区间/跨字段一致性等服务器端校验

### 动作（Actions）
- 完善
  - 动作元信息：确认对话、图标、颜色、权限键、可见性条件、是否为“批量动作”
  - 动作表单（带参数的动作）、行级动作与头部动作区分

### 关系管理（Relation Managers）
- 完善
  - 关系管理页面（hasMany/belongsToMany 的列表、Attach/Detach/Sync、内联创建关联）

### 权限/安全
- 完善
  - 前端可见性的动态控制（资源层 `editable/creatable/...` 现在固定为 true）
  - 动作/字段级权限；按字段过滤/排序/搜索白名单，防注入

### 国际化
- 完善
  - 后端系统化 i18n（资源/字段/动作/导航）；资源导航分组/排序/图标/徽章



### 全局搜索（Global Search）
- 缺失
  - 未检出全局搜索实现
- 建议
  - 在 `adminGroup` 增加 `/global-search`，聚合各资源 `GetSearchableFields()`，统一分页与高亮；`repository` 层做字段映射与白名单校验

### 通知与实时能力
- 已具备
  - 数据库存储通知、操作后中间件落库：`internal/middleware/notification.go`、`internal/service/notification_service.go`
- 缺失/可完善
  - 实时推送（WS/SSE/Broadcast）、点击通知跳转路由元信息、通知动作（标记/跳转）的一致性
- 建议
  - 新增 SSE 或 WebSocket 通道（如 `/api/admin/notifications/stream`），在 `CreateNotification` 时广播；前端订阅；保留落库与未读计数接口

### 多租户（Tenancy）
- 缺失
  - 未检出 tenant 相关中间件/字段
- 建议
  - 中间件解析租户（域名/Header/Token），在 `repository` 层自动注入 `tenant_id` 过滤；模型统一加 `TenantID`；权限/导航/全局搜索与导出均需按租户隔离；支持超级租户豁免

### 导入（Import）
- 已具备
  - 导出 CSV/Excel：`internal/service/export_service.go`（由 `resource_service.go` 调用）
- 缺失/可完善
  - 导入 CSV/Excel，字段映射、预检/回滚、错误报告，异步化（大文件）
- 建议
  - 新增 `import_service.go` + `/import` 路由，支持“解析->预检->提交”两阶段；支持任务化与通知

### 审计/操作日志
- 已具备
  - 操作日志中间件与统计接口
- 可完善
  - 字段级变更 diff；与资源/记录的关联跳转
- 建议
  - 在 `repository` 的 Update/Delete 做旧值对比，记录 diff，前端详情页展示变更历史

### 队列/任务与大文件
- 已具备
  - `internal/server/job.go`、`internal/task`、`cmd/task` 框架
- 可完善
  - 导出/导入/通知等异步化；队列状态、失败重试与进度查询
- 建议
  - 将导出/导入/报表/批量动作任务化，完成后发送通知并提供进度 API

### 文件与媒体库
- 已具备
  - 简单上传 `POST /api/admin/upload`，静态访问 `/uploads`
- 缺失/可完善
  - 多文件/分片/云存储直传/签名、图片裁剪与多尺寸转换、媒体库（目录/标签）、附件字段（多文件）
- 建议
  - 增加 `ImageField/AttachmentsField` 及处理管线；对接对象存储；生成缩略图与元数据

### 其他（可选）
- 伪装登录（Impersonate）、2FA、设备管理、速率限制/节流、表格可重排/树形/汇总、Infolists（详情展示组件）、CLI 脚手架生成资源/页面

### 优先级建议（落地路线）
- P0（立即提升使用体验与安全）
  - 字段/排序/过滤白名单与默认排序；`/global-search`；资源导航元信息与权限动态返回；导出任务化；通知 SSE
- P1（对标 Filament 表格/表单）
  - 批量动作框架与动作表单；表单容器与高级字段；关系管理器；导入
- P2（生态完善）
  - 多租户全链路、媒体库、变更 diff 审计、Impersonate/2FA、脚手架 CLI

以上建议都能在你现有的结构中平滑接入，关键是统一由 `pkg/admin/resource.go` 扩展标准元数据与白名单，然后由 `api_generator.go` 和 `repository/resource_repository.go` 承接路由与安全落地，前端即可据元数据自动渲染更丰富的 UI。