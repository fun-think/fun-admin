# fun-admin

<p align="center">
  A Filament-inspired admin framework for Go (Gin + Gorm) with a Vue 3 admin UI.
</p>

<p align="center">
  <a href="README.md"><kbd>English</kbd></a>
  <a href="README.zh-CN.md"><kbd>中文</kbd></a>
</p>

## What is this?

`fun-admin` is a Go admin panel framework and a runnable example app.

If you like the “Resource / Field / Action / Page” development experience of Filament, this project brings a similar idea to Go:

- Define admin **Resources** in Go code
- Get consistent CRUD, filters, columns, and actions
- Plug in authorization, navigation meta, i18n, import/export, etc.

The framework code lives in `pkg/admin`. The example backend is under `cmd/` + `internal/`. The frontend admin UI is under `web/`.

## Features

Framework (Filament-like building blocks):

- **Resources**: `admin.Resource` defines model binding, fields, columns, filters, actions
- **Fields**: typed field definitions (text/email/number/select/…)
- **Actions**: per-record and bulk actions, optional action forms, action visibility
- **Pages**: custom pages via `admin.PageInterface` / `admin.BasePage`
- **Navigation meta**: icon/group/sort/badge, and hide-from-navigation hooks
- **Authorization hooks**: optional interfaces (create/update/delete checks, etc.)
- **Field-level permissions**: readable/writable field control per context
- **API generation**: REST endpoints generated from resources (`pkg/admin/api_generator.go`)
- **i18n**: backend translation resources (`zh-CN`, `en`)

Example app capabilities (included in this repo):

- JWT admin login (`/api/admin/login`)
- Casbin permission middleware (RBAC-style)
- Operation logs + dashboard stats
- File upload, import/export endpoints
- SQLite config by default (`config/local.yml`)

## Requirements

- Go `1.23+` (see `go.mod`)
- Node.js `>= 18.15` (see `web/package.json`)
- Optional: Docker + Docker Compose (for MySQL/Redis in `deploy/docker-compose/`)

## Quick Start

### Backend

1. Run migration (creates tables + initial admin user):

```bash
go run ./cmd/migration -conf config/local.yml
```

Optional flags:

```bash
go run ./cmd/migration -conf config/local.yml -allow-drop=false -admin-password=your_password
```

2. Start the server:

```bash
go run ./cmd/server -conf config/local.yml
```

Default URLs (from `config/local.yml`):

- Admin API: `http://127.0.0.1:8001/api/admin`
- Swagger UI: `http://127.0.0.1:8001/swagger/index.html`

### Frontend

```bash
cd web
npm install
npm run dev
```

See `web/README.md` / `web/README.zh-CN.md` for the UI stack and scripts.

## Core Concepts

### Resources

Resources are the heart of the framework. A resource describes:

- Which model it manages (`GetModel`)
- Which fields are editable (`GetFields`)
- Which columns show in list views (`GetColumns`)
- Which filters are available (`GetFilters`)
- Which actions can be executed (`GetActions`)

Register resources in `cmd/bootstrap/bootstrap.go` (see `InitAdmin`).

### Actions (including bulk)

Implement `admin.ActionExecutor` on your resource to execute actions:

- Single record or bulk actions (`ids []interface{}`)
- Optional action form parameters (`params map[string]interface{}`)
- Optional visibility control per request/context

## Scaffolding

This repo includes a small scaffolding CLI:

```bash
go run ./cmd/make -action make:resource -name UserProfile
go run ./cmd/make -action make:page -name Reports
```

Then register the generated resource/page in `cmd/bootstrap/bootstrap.go`.

## Common Commands

```bash
make init      # install dev tools (mockgen, swag)
make bootstrap # docker compose + migration + server
make build     # build backend binary + frontend dist
make swag      # generate Swagger docs (outputs to ./docs/)
```

Notes:

- `make test` currently references a `test/` directory that is not present in this repo.
- `make admin-resource` / `make admin-page` currently points to `cmd/admin/main.go` which is not present; use `cmd/make` instead.

## Configuration

- Config files live in `config/` (examples: `config/local.yml`, `config/prod.yml`)
- Select config via:
  - CLI flag: `-conf config/local.yml`
  - Env var: `APP_CONF=config/local.yml`

Key settings:

- `env`: controls Gin mode (`prod/production` => Release; otherwise Debug)
- `http.host`, `http.port`
- `data.db.user.driver`, `data.db.user.dsn`
- `data.redis.*` (empty `data.redis.addr` => in-memory cache)

## Project Layout

- `pkg/admin`: the admin framework (resources/fields/actions/pages/api generator)
- `internal/resources`: example resources (users, operation logs, dictionaries, demo CRUD)
- `internal/handler`: Gin handlers for admin APIs
- `internal/service`: resource service layer (CRUD/action execution, validation, export)
- `internal/repository`: database access
- `cmd/server`: runnable HTTP server
- `cmd/migration`: migration + initial admin setup
- `web/`: admin UI

## License

See `LICENSE`.
