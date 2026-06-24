# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
# Run the server locally (port 8081)
make serve          # go run .

# Build for Linux deployment
make build          # CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/goadmin

# Generate CRUD table models from database
make generate       # adm generate -c adm.ini
```

## Testing

```bash
make test                    # Run both black-box and user-acceptance tests (requires admin_test.db)
make black-box-test          # go test -v -test.run=TestMainBlackBox
make user-acceptance-test    # go test -v -test.run=TestMainUserAcceptance (needs chromedriver)
```

Black-box tests use the GoAdmin test framework with `gin.NewHandler` and `httpexpect`. Acceptance tests use `agouti` + chromedriver for browser-based testing.

## Architecture

This is a game operations management (GM) web panel for an MMORPG (**Mserver**), built on the **GoAdmin** framework (`github.com/GoAdminGroup/go-admin`). The framework is replaced locally at `../go-admin@v1.2.23` via `go.mod` `replace` directive — modifications to the framework itself happen there.

### Key tiers

- **`main.go`** — Entry point. Initializes GoAdmin engine with configuration from `config.yml`, registers page generators, starts the fusion service framework, registers custom routes, and boots Gin on `:8081`.

- **`common/`** — Shared globals: `AdminEngine` (GoAdmin engine), `GinEngine` (Gin router), `GormDBList` / `ServerDBList` (DB connection pools), `MyCfg` (custom config), `Cfg_yml` (GoAdmin config). Also contains type definitions under `common/def/` for data structures (player, log, mail, server, items, etc.).

- **`fusion/`** — Custom service framework (imported as `admin/fusion`):
  - `ServiceBase` — Service lifecycle with per-second tick and timer wheel
  - `ServerMaster.go` — Signal handling, graceful shutdown, goroutine pools (`ants` with pools of 1000 and 10)
  - `Tools.go` — The "kitchen sink": DB helpers (`GetBaseGormDB`, `InitBaseDataBase`, `InitServerDataBases`), HTTP utilities (`CallToDeploy`, `CallToCenter`), config loading (`InitConfig`), table rendering helpers, pprof init, and many form/page building functions using the GoAdmin table DSL
  - `timer/` — Wheel timer implementation for scheduled tasks
  - `base/` — Panic-recovery wrapper (`SafeHandler`), simple numeric helpers

- **`mgr/`** — Central timer management (`timerMgr.go`). Registers all periodic background jobs (player online counting, auction monitoring, log archiving, data backups). Sub-package `pageMgr/` contains logic for specific management pages (server, mail, player, log managers).

- **`pages/`** — All admin panel pages:
  - `tables.go` — Maps URL prefixes to page generators (`Generators` map). Each entry corresponds to a GoAdmin CRUD table route (`/admin/info/<prefix>`).
  - `enter.go` — Registers custom HTML/data routes for non-CRUD pages
  - `AutoPages/` — GoAdmin-generated table definitions (one file per DB table). Each file defines the table schema, filters, form fields, and column rendering.
  - `CustomPages/` — Hand-built pages with custom business logic (player management, mail, server operations, charts, daily sign-in, reborn, battle grouping, item add, hotfix pipeline, etc.)

- **`hotfix/`** — Hotfix pipeline module for data table hot-reload:
  - `pipeline.go` — Full pipeline orchestration (build → export → archive → deploy → GM reload)
  - `session.go` — Pipeline session tracking with progress updates
  - `gm_commands.go` — GM command dispatch to game servers
  - `table_names.go` — Table name resolution and special table detection
  - `build.go` / `export.go` / `archive.go` / `upload.go` — Build and deployment steps
  - `config.go` — Hotfix-specific configuration

- **`adm.ini`** — CLI config for the `adm` code generator tool

- **`config.yml`** — GoAdmin framework config (DB connections, theme, language, logging). Contains **credentials** — keep out of commits.

- **`SQL/`** — Database schema dumps for initializing the system tables

### Database connections

Six MySQL databases defined in `config.yml`: `default` (go-admin system), `db_global` (mmorpg_global), `go_manager`, `go_backup`, `db_world` (mmorpg_world_newui), `db_log` (mmorpg_log), `db_login_log` (reporter). Server-specific connections are lazy-initialized per game server.

**Database naming convention**: All server databases use the `_s{serverID}` suffix uniformly (e.g., `mmorpg_log_s1`, `mmorpg_log_s2`), including server ID 1.

### GM Command API Routing

GM commands are dispatched to different server types via center APIs:
- `MapServer` → `GM2MS`
- `GameServer` → `GM2GS`
- `GateServer` → `GM2GATE`
- `SocialServer` → `GM2SOCIAL`
- `DBPServer` → `GM2DBP`
- `AdminServer` / others → `GM2S` (default)

Hotfix data table commands broadcast to all API endpoints to ensure all server types receive the update.

### Adding a new admin page

1. **CRUD page** (simple table): Create a file in `pages/AutoPages/` defining the table generator function, then add an entry in `pages/tables.go` mapping a URL prefix to the generator.
2. **Custom page** (complex logic): Add handler in `pages/CustomPages/`, register route in `pages/enter.go`, and add table entry in `pages/tables.go` if it needs a CRUD listing.

### i18n

Localization happens via `github.com/leonelquinteros/gotext`. Translation files live in `path/` (`.po`/`.mo` files). Menu translations are in `common/menuTranslate/`.

### Email

`SendEMail/SendMail.go` — Sends email notifications via SMTP to all users in the `goadmin_users` table, using credentials from `MyCfg.Email`.

---

# 维护文档（中文）

## 构建和运行

```bash
# 本地运行 (端口 8081)
make serve          # go run .

# 构建 Linux 部署包
make build          # CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/goadmin

# 从数据库生成 CRUD 表模型
make generate       # adm generate -c adm.ini
```

## 测试

```bash
make test                    # 运行黑盒测试和用户验收测试（需要 admin_test.db）
make black-box-test          # go test -v -test.run=TestMainBlackBox
make user-acceptance-test    # go test -v -test.run=TestMainUserAcceptance（需要 chromedriver）
```

## 架构概览

本项目是 **Mserver** 的 MMORPG 游戏运营管理（GM）Web 面板，基于 **GoAdmin** 框架构建。框架通过 `go.mod` 的 `replace` 指令指向本地 `../go-admin@v1.2.23`。

### 核心模块

- **`main.go`** — 入口，初始化 GoAdmin 引擎、注册页面生成器、启动融合服务框架、注册自定义路由，Gin 监听 `:8081`
- **`common/`** — 全局共享：引擎实例、数据库连接池、自定义配置、数据结构定义（玩家、日志、邮件、服务器、道具等）
- **`fusion/`** — 自定义服务框架：服务生命周期、信号处理、协程池（ants 1000/10）、数据库辅助函数、HTTP 工具、配置加载、表格渲染
- **`mgr/`** — 定时任务管理：玩家在线统计、拍卖监控、日志归档、数据备份
- **`pages/`** — 管理页面：`AutoPages/` 自动生成的 CRUD 表、`CustomPages/` 自定义业务页面（角色管理、邮件、服务器操作、图表、签到、重生、战斗分组、道具添加、热更流水线等）
- **`hotfix/`** — 热更数据表模块：构建→导出→归档→部署→GM热加载的完整流水线

### 数据库连接

6 个 MySQL 数据库：`default`、`db_global`、`go_manager`、`go_backup`、`db_world`、`db_log`、`db_login_log`。各游戏服数据库按需延迟初始化。

**数据库命名规范**：所有服务器统一使用 `_s{serverID}` 后缀（如 `mmorpg_log_s1`、`mmorpg_log_s2`），包括 s1。

### GM 命令 API 路由

| 服务器类型 | API 配置键 |
|---|---|
| MapServer | GM2MS |
| GameServer | GM2GS |
| GateServer | GM2GATE |
| SocialServer | GM2SOCIAL |
| DBPServer | GM2DBP |
| AdminServer 等 | GM2S（默认） |

热更数据表命令会向所有 API 端点广播，确保所有服务器类型都收到更新。

### 新增管理页面

1. **CRUD 页面**：在 `pages/AutoPages/` 创建表生成函数，在 `pages/tables.go` 注册 URL 映射
2. **自定义页面**：在 `pages/CustomPages/` 添加 handler，在 `pages/enter.go` 注册路由

### 国际化

使用 `github.com/leonelquinteros/gotext`，翻译文件在 `path/`（`.po`/`.mo`），菜单翻译在 `common/menuTranslate/`。

### 邮件

`SendEMail/SendMail.go` — 通过 SMTP 向 `goadmin_users` 表中所有用户发送邮件通知。
