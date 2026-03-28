# Plugin Platform - 插件平台（简化版）

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
- ✅ **前端友好** - 简化后的数据格式，支持动态 UI 渲染

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
# 调用计算器插件（简化后：不再需要 method 参数）
curl -X POST http://localhost:8080/api/v1/plugins/calculator/execute \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 20, "operator": "+"}'

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
| `/api/v1/plugins/:id/execute` | POST | 执行插件（简化后） |
| `/api/v1/plugins/:id/health` | GET | 健康检查 |

### 注册插件（简化版格式）

```bash
POST /api/v1/plugins/register
Content-Type: application/json

{
  "name": "calculator",
  "version": "1.0.0",
  "endpoint": "http://localhost:8001",
  "summary": "简单计算器",
  "params": {
    "a": {"type": "number", "required": true},
    "b": {"type": "number", "required": true},
    "operator": {"type": "select", "options": ["+", "-", "*", "/"], "default": "+"}
  },
  "outputs": "number"
}
```

## 参数类型（支持前端动态渲染）

| 类型 | 前端渲染组件 | 额外属性 |
|------|-------------|----------|
| `string` | 文本输入框 | `maxLength`, `placeholder` |
| `number` | 数字输入框 | `min`, `max`, `step` |
| `file` | 文件上传 | `accept`, `multiple`, `maxSize` |
| `boolean` | 开关/复选框 | - |
| `select` | 下拉选择 | `options`, `multiple` |
| `textarea` | 多行文本 | `rows`, `maxLength` |
| `password` | 密码输入 | - |
| `date` | 日期选择器 | `min`, `max` |
| `object` | JSON 编辑器 | `schema` |
| `array` | 动态列表 | `itemType` |

### 参数定义示例

```json
{
  "params": {
    "username": {
      "type": "string",
      "required": true,
      "maxLength": 50,
      "placeholder": "请输入用户名"
    },
    "avatar": {
      "type": "file",
      "required": true,
      "accept": ".jpg,.png",
      "maxSize": 5242880
    },
    "category": {
      "type": "select",
      "options": ["all", "image", "doc"],
      "default": "all"
    },
    "compress": {
      "type": "boolean",
      "default": true
    }
  }
}
```

## 开发插件

插件需要实现以下 HTTP 端点：

1. **健康检查** `GET /health`
2. **执行** `POST /execute`（简化后：单一入口）

### 简化后的插件结构

```go
package main

import (
    "encoding/json"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", handleHealth)
    mux.HandleFunc("/execute", handleExecute)  // 单一入口
    http.ListenAndServe(":8001", mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func handleExecute(w http.ResponseWriter, r *http.Request) {
    var params map[string]interface{}
    json.NewDecoder(r.Body).Decode(&params)
    
    // 处理逻辑...
    result := process(params)
    
    json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
}
```

完整示例代码见 `plugins/example-calculator/main.go`

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

## 变更记录

### v2.0.0 - 格式简化版

- **BREAKING CHANGE**: 移除 `methods` 层级，一个插件对应单一功能
- **BREAKING CHANGE**: 执行接口从 `/:id/execute/:method` 改为 `/:id/execute`
- **BREAKING CHANGE**: `returns` 字段重命名为 `outputs`
- **新增**: `params` 支持类型特定属性（`accept`, `options`, `maxSize` 等）
- **优化**: 注册数据格式从 ~60 行减少到 ~15 行
