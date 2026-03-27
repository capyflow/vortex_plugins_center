# Plugin SDK for Go

Go 语言插件 SDK，用于快速开发插件并注册到 Vortex Plugin Platform。

## 特性

- ✅ **链式调用** - 优雅的 API 设计
- ✅ **类型安全** - 支持多种参数类型
- ✅ **自动注册** - 服务启动后自动注册到平台
- ✅ **完整文档** - 支持文档地址配置

## 安装

```bash
go get github.com/capyflow/vortex_plugins_center/sdk/go
```

## 快速开始

### 1. 创建插件

```go
package main

import (
    "context"
    "log"
    "github.com/capyflow/vortex_plugins_center/sdk/go"
)

func main() {
    // 创建 SDK
    sdk := plugin.New(plugin.RegisterOptions{
        PlatformURL: "http://localhost:8080",
        Name:        "my-plugin",
        Version:     "1.0.0",
        Description: "我的第一个插件",
        Endpoint:    "http://localhost:9001",
        DocURL:      "https://docs.example.com/my-plugin",
    })

    // 添加方法
    sdk.AddMethod(
        plugin.MethodOptions{
            Name:        "hello",
            Description: "打招呼",
            Path:        "/hello",
            HTTPMethod:  "POST",
        },
        []plugin.Parameter{
            plugin.StringParam("name", "你的名字", true),
        },
        plugin.ReturnInfo{
            Type:        "string",
            Description: "问候语",
        },
    )

    // 注册到平台
    ctx := context.Background()
    resp, err := sdk.Register(ctx)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("注册成功: %s", resp.ID)
}
```

### 2. 参数类型

```go
// 字符串参数
plugin.StringParam("name", "描述", true)

// 数字参数
plugin.NumberParam("age", "年龄", true)

// 布尔参数
plugin.BoolParam("active", "是否激活", false)

// 对象参数
plugin.ObjectParam("config", "配置对象", true)

// 数组参数
plugin.ArrayParam("items", "项目列表", true)
```

### 3. 完整示例

见 `example/main.go`

## API 参考

### RegisterOptions

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| PlatformURL | string | ✓ | 平台地址 |
| Name | string | ✓ | 插件名称 |
| Version | string | ✓ | 版本号 |
| Description | string | ✓ | 插件描述 |
| Endpoint | string | ✓ | 插件HTTP端点 |
| DocURL | string | - | 文档地址 |

### MethodOptions

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Name | string | ✓ | 方法名 |
| Description | string | ✓ | 方法描述 |
| Path | string | ✓ | HTTP路径 |
| HTTPMethod | string | ✓ | GET/POST/PUT/DELETE |

### Parameter

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Name | string | ✓ | 参数名 |
| Type | string | ✓ | string/number/boolean/object/array |
| Required | bool | ✓ | 是否必填 |
| Default | string | - | 默认值 |
| Description | string | - | 参数描述 |

## 完整示例

```go
package main

import (
    "context"
    "log"
    "net/http"
    "encoding/json"
    "github.com/capyflow/vortex_plugins_center/sdk/go"
)

func main() {
    // 创建 SDK
    sdk := plugin.New(plugin.RegisterOptions{
        PlatformURL: "http://localhost:8080",
        Name:        "calculator",
        Version:     "1.0.0",
        Description: "计算器插件",
        Endpoint:    "http://localhost:9001",
        DocURL:      "https://docs.example.com/calculator",
    })

    // 添加方法
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

    // 启动 HTTP 服务
    go func() {
        http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
            var req struct{ A, B float64 }
            json.NewDecoder(r.Body).Decode(&req)
            json.NewEncoder(w).Encode(map[string]float64{
                "result": req.A + req.B,
            })
        })
        http.ListenAndServe(":9001", nil)
    }()

    // 注册
    resp, err := sdk.Register(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("注册成功: %s", resp.ID)
}
```

## 运行

```bash
# 启动平台
cd ../..
go run cmd/main.go

# 启动插件
cd example
go run main.go
```

## 许可证

MIT
