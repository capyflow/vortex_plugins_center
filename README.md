# Plugin Platform - 插件平台

基于 VortexV3 实现的 HTTP 插件平台，支持插件通过 HTTP 接口注册和调用。

## 架构

```
用户 → API Gateway → Plugin Center → Registry / Router → Plugin (HTTP)
```

## 核心特性

- ✅ **HTTP 注册** - 插件通过 HTTP 接口注册到平台
- ✅ **位置无关** - 插件可以运行在本地、远程、Docker、Serverless
- ✅ **动态路由** - 平台自动路由请求到正确的插件端点
- ✅ **健康检查** - 自动检测插件健康状态
- ✅ **负载均衡** - 支持同一插件多实例部署

## 快速开始

### 1. 启动平台

```bash
cd plugin-platform
go run cmd/main.go
```

平台将启动在 `:8080`

### 2. 启动示例插件

```bash
cd plugins/example-calculator
go run main.go
```

插件将启动在 `:8001`，并自动注册到平台

### 3. 调用插件

```bash
# 调用计算器插件的 add 方法
curl -X POST http://localhost:8080/api/v1/plugins/calculator/execute/add \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 20}'

# 响应
{"success": true, "result": {"result": 30}, "latency": 5}
```

## API 接口

### 插件管理

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/v1/plugins/register` | POST | 注册插件 |
| `/api/v1/plugins/:id` | DELETE | 注销插件 |
| `/api/v1/plugins/list` | GET | 获取插件列表 |
| `/api/v1/plugins/:id` | GET | 获取插件详情 |
| `/api/v1/plugins/:id/execute/:method` | POST | 执行插件方法 |
| `/api/v1/plugins/:id/health` | GET | 健康检查 |

### 注册插件

```bash
POST /api/v1/plugins/register
Content-Type: application/json

{
  "name": "calculator",
  "version": "1.0.0",
  "description": "Simple calculator",
  "endpoint": "http://localhost:8001",
  "methods": [
    {
      "name": "add",
      "description": "Add two numbers",
      "path": "/add",
      "method": "POST",
      "parameters": [
        {"name": "a", "type": "number", "required": true},
        {"name": "b", "type": "number", "required": true}
      ]
    }
  ]
}
```

## 开发插件

插件需要实现以下 HTTP 端点：

1. **健康检查** `GET /health`
2. **业务方法** 自定义路径

示例代码见 `plugins/example-calculator/main.go`

## 项目结构

```
plugin-platform/
├── cmd/
│   └── main.go              # 入口
├── internal/
│   ├── gateway/             # API 网关
│   ├── center/              # 插件中心
│   ├── registry/            # 注册表
│   └── router/              # 路由器
├── pkg/
│   └── models/              # 数据模型
├── plugins/                 # 示例插件
│   └── example-calculator/  # 计算器插件
├── go.mod
└── README.md
```

## 技术栈

- Go 1.21+
- VortexV3 (HTTP 框架)
- MongoDB (数据存储)
- Redis (缓存)
