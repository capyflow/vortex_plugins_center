# Plugin Protocol - 插件通信协议

## 概述

定义 Center 和 Plugin 之间的标准通信格式，所有插件必须实现此协议才能被平台调用。

---

## 请求格式

### HTTP Headers

| Header | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Content-Type | string | ✓ | `application/json` |
| X-Request-ID | string | ✓ | 请求唯一标识 |
| X-Plugin-Name | string | ✓ | 插件名称 |
| X-Method | string | ✓ | 调用的方法名 |
| X-Timestamp | int64 | ✓ | 请求时间戳（毫秒） |

### Request Body

```json
{
  "request_id": "req-abc123",
  "plugin_name": "calculator",
  "method": "add",
  "timestamp": 1711500000000,
  "params": {
    "a": 10,
    "b": 20
  }
}
```

---

## 响应格式

### 成功响应

```json
{
  "success": true,
  "request_id": "req-abc123",
  "data": {
    "result": 30
  },
  "meta": {
    "duration": 5,
    "timestamp": 1711500000005
  }
}
```

### 错误响应

```json
{
  "success": false,
  "request_id": "req-abc123",
  "error": {
    "code": "CALCULATION_ERROR",
    "message": "Division by zero"
  },
  "meta": {
    "duration": 2,
    "timestamp": 1711500000002
  }
}
```

---

## 插件必须实现的接口

### 1. 健康检查

```
GET /health
Response: { "status": "healthy", "version": "1.0.0" }
```

### 2. 方法调用

```
POST /invoke
Headers:
  X-Request-ID: req-xxx
  X-Plugin-Name: calculator
  X-Method: add

Body: { "request_id": "req-xxx", "method": "add", "params": {...} }
Response: { "success": true, "data": {...}, "meta": {...} }
```

---

## 标准错误码

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| INVALID_REQUEST | 请求格式错误 | 400 |
| METHOD_NOT_FOUND | 方法不存在 | 404 |
| INVALID_PARAMS | 参数错误 | 400 |
| EXECUTION_ERROR | 执行错误 | 500 |
| TIMEOUT | 执行超时 | 504 |
