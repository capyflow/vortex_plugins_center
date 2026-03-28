# Plugin Platform v2.0.0 变更文档

## 概述

本次更新对插件平台进行了重大简化，主要目标：
1. **简化插件注册格式** - 减少冗余字段，降低开发成本
2. **单一功能原则** - 一个插件对应一个功能，不再支持多 methods
3. **增强前端支持** - 参数类型支持更多 UI 渲染属性

---

## 破坏性变更 (Breaking Changes)

### 1. 数据模型变更

#### Plugin 结构体

```go
// 旧版本
type Plugin struct {
    Name        string
    Description string
    Methods     []PluginMethod  // 支持多个方法
    ...
}

// 新版本
type Plugin struct {
    Name     string
    Summary  string            // 替代 Description
    Params   map[string]ParamDef  // 直接定义参数
    Outputs  string            // 替代 Returns
    ...
}
```

#### 移除的字段

| 字段 | 说明 |
|------|------|
| `Description` | 使用 `Summary` 替代 |
| `Methods` | 移除，一个插件单一功能 |
| `PluginMethod` | 移除 |
| `Returns` | 使用 `Outputs` 替代 |

#### 新增/修改的字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `Summary` | string | 简短描述 |
| `Params` | map[string]ParamDef | 参数定义（key 为参数名） |
| `Outputs` | string | 输出类型 |

### 2. API 变更

#### 执行接口

```bash
# 旧版本
POST /api/v1/plugins/:id/execute/:method

# 新版本
POST /api/v1/plugins/:id/execute
```

**变更说明**：不再需要 method 参数，每个插件只有一个执行入口。

### 3. 注册格式变更

#### 旧格式（约 60 行）

```json
{
  "name": "calculator",
  "version": "1.0.0",
  "description": "Simple calculator plugin",
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
      ],
      "returns": {"type": "number", "description": "Sum"}
    },
    {
      "name": "sub",
      "description": "Subtract two numbers",
      "path": "/sub",
      "method": "POST",
      "parameters": [...]
    }
  ]
}
```

#### 新格式（约 15 行）

```json
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

---

## 参数类型系统

### 支持的类型

| 类型 | 说明 | 额外属性 |
|------|------|----------|
| `string` | 文本输入 | `maxLength`, `placeholder` |
| `number` | 数字输入 | `min`, `max`, `step` |
| `file` | 文件上传 | `accept`, `multiple`, `maxSize` |
| `boolean` | 开关 | - |
| `select` | 下拉选择 | `options`, `multiple` |
| `textarea` | 多行文本 | `rows`, `maxLength` |
| `password` | 密码输入 | - |
| `date` | 日期选择 | `min`, `max` |
| `object` | JSON 对象 | `schema` |
| `array` | 数组 | `itemType` |

### 参数定义结构

```go
type ParamDef struct {
    Type        string   // 参数类型
    Required    bool     // 是否必填
    Default     any      // 默认值
    Description string   // 描述
    
    // 类型特定属性
    MaxLength   int      // string/textarea
    Min         float64  // number
    Max         float64  // number
    Step        float64  // number
    Accept      string   // file: ".jpg,.png"
    Multiple    bool     // file/select
    MaxSize     int64    // file: bytes
    Options     []string // select
    Rows        int      // textarea
    Placeholder string   // string/textarea
}
```

---

## 前端渲染支持

### 动态表单生成

基于 `params` 定义，前端可以自动生成表单：

```typescript
// 伪代码
function renderPluginForm(plugin: Plugin) {
  return (
    <Form>
      {Object.entries(plugin.params).map(([name, def]) => (
        <Form.Item key={name} label={name}>
          {renderInputByType(name, def)}
        </Form.Item>
      ))}
    </Form>
  );
}

function renderInputByType(name: string, def: ParamDef) {
  switch (def.type) {
    case 'file':
      return <Upload accept={def.accept} maxSize={def.maxSize} />;
    case 'select':
      return <Select options={def.options} />;
    case 'number':
      return <InputNumber min={def.min} max={def.max} />;
    // ... 其他类型
  }
}
```

---

## 迁移指南

### 1. 更新模型定义

```bash
# 更新后的模型文件
pkg/models/plugin.go
```

### 2. 更新插件代码

#### 修改注册数据

```go
// 旧代码
reqBody := map[string]interface{}{
    "name": "my-plugin",
    "methods": []map[string]interface{}{
        {"name": "method1", "parameters": [...]},
        {"name": "method2", "parameters": [...]},
    },
}

// 新代码
reqBody := map[string]interface{}{
    "name": "my-plugin",
    "params": map[string]interface{}{
        "param1": {"type": "string", "required": true},
        "param2": {"type": "number"},
    },
    "outputs": "object",
}
```

#### 修改 HTTP 端点

```go
// 旧代码
mux.HandleFunc("/method1", handleMethod1)
mux.HandleFunc("/method2", handleMethod2)

// 新代码
mux.HandleFunc("/execute", handleExecute)  // 单一入口
```

### 3. 更新 API 调用

```bash
# 旧调用
POST /api/v1/plugins/my-plugin/execute/method1

# 新调用
POST /api/v1/plugins/my-plugin/execute
```

---

## 文件变更列表

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `pkg/models/plugin.go` | 修改 | 简化模型定义 |
| `internal/center/center.go` | 修改 | 适配新模型，简化 Execute 方法 |
| `internal/gateway/gateway.go` | 修改 | 移除 method 参数 |
| `internal/router/router.go` | 修改 | 简化路由结构 |
| `plugins/example-calculator/main.go` | 修改 | 更新示例代码 |
| `README.md` | 修改 | 更新文档 |
| `CHANGELOG_v2.md` | 新增 | 本变更文档 |

---

## 优势对比

| 维度 | 旧版本 | 新版本 |
|------|--------|--------|
| 注册代码行数 | ~60 行 | ~15 行 |
| 嵌套层级 | 5 层 | 2 层 |
| 前端解析复杂度 | 高 | 低 |
| 插件开发成本 | 高 | 低 |
| 概念数量 | 多（methods, parameters, returns...） | 少（params, outputs） |
| 灵活性 | 一个插件多功能 | 一个插件单一功能（更清晰） |

---

## 注意事项

1. **这是一个破坏性更新** - 旧版插件需要修改后才能在新平台运行
2. **数据库迁移** - 如果已存储旧格式数据，需要迁移脚本
3. **前端适配** - 前端需要根据新格式调整渲染逻辑
4. **版本标记** - 建议标记为 v2.0.0，与 v1.x 区分

---

## 后续计划

- [ ] 添加参数校验中间件
- [ ] 支持文件上传代理
- [ ] 插件版本管理
- [ ] 插件市场（Web UI）
