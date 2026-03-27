# Vortex Plugin Platform - 插件平台规范

## 目录

1. [概述](#概述)
2. [整体流程](#整体流程)
3. [注册规范](#注册规范)
4. [HTTP 接口规范](#http-接口规范)
5. [多语言 SDK](#多语言-sdk)
6. [示例代码](#示例代码)

---

## 概述

Vortex Plugin Platform 是一个基于 HTTP 的插件平台，支持插件通过 HTTP 接口注册到平台，平台不关心插件的运行位置，只通过 HTTP 调用获取结果。

### 核心特性

- ✅ **位置无关** - 插件可运行在本地、远程、Docker、Serverless
- ✅ **动态注册** - 插件启动时主动注册到平台
- ✅ **统一路由** - 平台通过插件 ID 路由到正确的端点
- ✅ **多语言支持** - Go、Java、Python、Node.js 等

---

## 整体流程

### 1. 系统架构

```
┌─────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  User   │────▶│ API Gateway │────▶│Plugin Center│────▶│   Plugin    │
│ (用户)   │     │  (网关)      │     │  (插件中心)  │     │  (插件)      │
└─────────┘     └─────────────┘     └─────────────┘     └─────────────┘
     │                                                    │
     │  POST /api/v1/plugins/calculator/execute/add      │
     │  { "a": 1, "b": 2 }                               │
     │                                                    │
     │◀───────────────────────────────────────────────────│
     │                                                    │
     │  { "success": true, "result": 3 }                 │
```

### 2. 插件生命周期

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│ 开发插件 │───▶│ 启动服务 │───▶│ 注册平台 │───▶│ 接收调用 │
└─────────┘    └─────────┘    └─────────┘    └─────────┘
     │              │              │              │
     │              │              │              │
     ▼              ▼              ▼              ▼
  编写代码      启动HTTP       POST /register   执行业务
  定义接口       服务          上报端点信息      逻辑
```

### 3. 调用流程

1. **用户** 向平台发送请求
2. **平台** 查询注册表获取插件端点
3. **平台** 转发请求到插件
4. **插件** 执行业务逻辑
5. **插件** 返回结果给平台
6. **平台** 返回结果给用户

---

## 注册规范

### 1. 注册请求

**URL:** `POST /api/v1/plugins/register`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "calculator",
  "version": "1.0.0",
  "description": "计算器插件",
  "endpoint": "http://localhost:8001",
  "doc_url": "https://docs.example.com/calculator",
  "methods": [
    {
      "name": "add",
      "description": "加法运算",
      "path": "/add",
      "http_method": "POST",
      "parameters": [
        {
          "name": "a",
          "type": "number",
          "required": true,
          "description": "第一个加数"
        }
      ],
      "returns": {
        "type": "number",
        "description": "两数之和"
      }
    }
  ],
  "metadata": {
    "author": "Vortex Team",
    "license": "MIT"
  }
}
```

### 2. 字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | ✓ | 插件名称，全局唯一 |
| version | string | ✓ | 版本号，遵循 SemVer |
| description | string | ✓ | 插件描述 |
| endpoint | string | ✓ | 插件HTTP端点 |
| doc_url | string | - | 文档地址 |
| methods | array | ✓ | 方法列表 |
| metadata | object | - | 元数据 |

### 3. 参数类型

| 类型 | 说明 |
|------|------|
| string | 字符串 |
| number | 数字 |
| boolean | 布尔 |
| object | 对象 |
| array | 数组 |

---

## HTTP 接口规范

### 1. 插件必须实现的端点

```
GET /health
Response: { "status": "healthy" }
```

### 2. 错误处理

```json
{
  "error": "Division by zero",
  "code": "CALCULATION_ERROR"
}
```

---

## 多语言 SDK

### Go SDK

```go
sdk := plugin.New(plugin.RegisterOptions{
    PlatformURL: "http://localhost:8080",
    Name:        "calculator",
    Version:     "1.0.0",
    Description: "计算器插件",
    Endpoint:    "http://localhost:8001",
    DocURL:      "https://docs.example.com/calculator",
})

sdk.AddMethod(
    plugin.MethodOptions{
        Name:        "add",
        Description: "加法",
        Path:        "/add",
        HTTPMethod:  "POST",
    },
    []plugin.Parameter{
        plugin.NumberParam("a", "第一个数", true),
        plugin.NumberParam("b", "第二个数", true),
    },
    plugin.ReturnInfo{Type: "number", Description: "和"},
)

resp, err := sdk.Register(context.Background())
```

### Java SDK

```java
SDK sdk = new SDK(new RegisterOptions()
    .setPlatformURL("http://localhost:8080")
    .setName("calculator")
    .setVersion("1.0.0")
    .setDescription("计算器插件")
    .setEndpoint("http://localhost:8001")
    .setDocURL("https://docs.example.com/calculator")
);

sdk.addMethod(new MethodOptions()
    .setName("add")
    .setDescription("加法")
    .setPath("/add")
    .setHttpMethod("POST")
    .addParam(Parameter.number("a", "第一个数", true))
    .addParam(Parameter.number("b", "第二个数", true))
    .setReturns("number", "和")
);

RegisterResponse resp = sdk.register();
```

### Python SDK

```python
from vortex_plugin import SDK, RegisterOptions, MethodOptions, Parameter

sdk = SDK(RegisterOptions(
    platform_url="http://localhost:8080",
    name="calculator",
    version="1.0.0",
    description="计算器插件",
    endpoint="http://localhost:8001",
    doc_url="https://docs.example.com/calculator"
))

sdk.add_method(MethodOptions(
    name="add",
    description="加法",
    path="/add",
    http_method="POST",
    parameters=[
        Parameter.number("a", "第一个数", required=True),
        Parameter.number("b", "第二个数", required=True)
    ],
    returns={"type": "number", "description": "和"}
))

resp = sdk.register()
```

### Node.js SDK

```javascript
const { SDK } = require('vortex-plugin-sdk');

const sdk = new SDK({
  platformURL: 'http://localhost:8080',
  name: 'calculator',
  version: '1.0.0',
  description: '计算器插件',
  endpoint: 'http://localhost:8001',
  docURL: 'https://docs.example.com/calculator'
});

sdk.addMethod({
  name: 'add',
  description: '加法',
  path: '/add',
  httpMethod: 'POST',
  parameters: [
    { name: 'a', type: 'number', required: true, description: '第一个数' },
    { name: 'b', type: 'number', required: true, description: '第二个数' }
  ],
  returns: { type: 'number', description: '和' }
});

const resp = await sdk.register();
```

---

## 示例代码

### 完整计算器插件

见 `sdk/go/example/main.go`

---

## 目录结构

```
vortex_plugins_center/
├── cmd/                    # 平台入口
├── internal/               # 内部实现
│   ├── gateway/           # API网关
│   ├── center/            # 插件中心
│   ├── registry/          # 注册表
│   └── router/            # 路由器
├── sdk/                    # SDK目录
│   ├── go/                # Go SDK
│   ├── java/              # Java SDK
│   ├── python/            # Python SDK
│   └── nodejs/            # Node.js SDK
├── webui/                  # Web UI
├── docs/                   # 文档
│   └── SPECIFICATION.md   # 本规范
└── plugins/                # 示例插件
    └── example-calculator/
```
