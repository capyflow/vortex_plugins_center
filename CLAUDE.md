# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

这是一个基于 VortexV3 HTTP 框架的插件平台，支持插件通过 HTTP 接口注册和调用。平台不关心插件的运行位置（本地、远程、Docker、Serverless），只通过 HTTP 调用获取结果。

## 常用命令

```bash
# 启动平台（需要 MongoDB 和 Redis）
go run cmd/main.go

# 启动示例插件
cd plugins/example-calculator
go run main.go

# 构建
go build -o plugin-platform cmd/main.go

# 依赖管理
go mod tidy
go mod download
```

## 架构概览

```
用户 → API Gateway → Plugin Center → Registry / Router → Plugin (HTTP)
```

### 核心组件

| 组件 | 路径 | 职责 |
|------|------|------|
| Gateway | `internal/gateway/` | HTTP API 入口，使用 VortexV3 框架 |
| PluginCenter | `internal/center/` | 业务逻辑，协调注册、路由、执行 |
| Registry | `internal/registry/` | 插件数据持久化（MongoDB + Redis 缓存） |
| Router | `internal/router/` | 方法路由映射（内存 + Redis） |

### 数据流

1. **插件注册**: Plugin → HTTP POST `/api/v1/plugins/register` → Registry 存储到 MongoDB，Router 建立路由映射
2. **方法执行**: User → Gateway → Center → Router 查询端点 → HTTP 调用 Plugin → 返回结果
3. **健康检查**: Center 每 30 秒轮询活跃插件的 `/health` 端点

## 依赖的私有仓库

```
github.com/capyflow/vortexv3      # HTTP 框架
github.com/capyflow/allspark-go   # 数据库和日志工具
```

配置 GOPRIVATE 和 Git SSH 以访问：
```bash
go env -w GOPRIVATE="github.com/capyflow/*"
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

## 插件协议

插件必须实现以下 HTTP 端点：

1. **健康检查**: `GET /health` → `{ "status": "healthy" }`
2. **方法调用**: `POST /invoke`（标准协议见 `docs/PLUGIN_PROTOCOL.md`）

请求 Headers: `X-Request-ID`, `X-Plugin-Name`, `X-Method`, `X-Timestamp`

## 关键类型

```go
// pkg/models/plugin.go
Plugin       // 插件定义，存储在 MongoDB
PluginMethod // 方法定义，包含 HTTP Path 和 Method
PluginStatus // active/inactive/error
PluginHealth // 健康状态（状态、延迟、检查时间）
```

## SDK 和示例

- **Go SDK**: `sdk/go/plugin.go` - 提供 `plugin.New()` 和注册方法
- **示例插件**: `plugins/example-calculator/main.go` - 展示如何实现一个完整插件

## 项目目录结构

```
cmd/                    # 平台入口 main.go
internal/
  gateway/              # HTTP API 层（VortexV3）
  center/               # 业务逻辑层
  registry/             # MongoDB + Redis 存储
  router/               # 路由映射
pkg/models/             # 数据模型
sdk/go/                 # Go SDK
plugins/                # 示例插件
webui/                  # Web UI (index.html)
docs/                   # 规范文档
```

## 环境要求

- Go 1.25+
- MongoDB（数据存储）
- Redis（缓存和消息）

## 注意事项

- VortexV3 的 HTTP Context 在 `github.com/capyflow/vortexv3/server/http`，使用 `vhttp.Codes` 和 `vpkg.VMsgResponse`
- 插件注册时会验证端点可访问性（health check）
- 同一插件多版本通过 name 查询最新版本更新，而非创建多条记录
- Router 使用内存 map + RWMutex 管理路由表，每 5 分钟清理过期路由
