package models

import (
	"time"
)

// Plugin 插件定义
type Plugin struct {
	ID          string            `json:"id" bson:"_id"`
	Name        string            `json:"name" bson:"name"`
	Version     string            `json:"version" bson:"version"`
	Description string            `json:"description" bson:"description"`
	Endpoint    string            `json:"endpoint" bson:"endpoint"` // HTTP 端点
	Methods     []PluginMethod    `json:"methods" bson:"methods"`
	Metadata    map[string]string `json:"metadata" bson:"metadata"`
	Status      PluginStatus      `json:"status" bson:"status"`
	Health      PluginHealth      `json:"health" bson:"health"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
	LastSeen    time.Time         `json:"last_seen" bson:"last_seen"`
}

// PluginMethod 插件方法定义
type PluginMethod struct {
	Name        string            `json:"name" bson:"name"`
	Description string            `json:"description" bson:"description"`
	Path        string            `json:"path" bson:"path"` // HTTP 路径
	Method      string            `json:"method" bson:"method"` // GET/POST/PUT/DELETE
	Parameters  []Parameter       `json:"parameters" bson:"parameters"`
	Returns     ReturnType        `json:"returns" bson:"returns"`
}

// Parameter 参数定义
type Parameter struct {
	Name     string `json:"name" bson:"name"`
	Type     string `json:"type" bson:"type"`
	Required bool   `json:"required" bson:"required"`
	Default  string `json:"default" bson:"default"`
}

// ReturnType 返回类型
type ReturnType struct {
	Type        string `json:"type" bson:"type"`
	Description string `json:"description" bson:"description"`
}

// PluginStatus 插件状态
type PluginStatus string

const (
	PluginStatusActive   PluginStatus = "active"
	PluginStatusInactive PluginStatus = "inactive"
	PluginStatusError    PluginStatus = "error"
)

// PluginHealth 插件健康状态
type PluginHealth struct {
	Status    string    `json:"status" bson:"status"`
	Latency   int64     `json:"latency" bson:"latency"` // ms
	CheckedAt time.Time `json:"checked_at" bson:"checked_at"`
}

// RegisterRequest 插件注册请求
type RegisterRequest struct {
	Name        string            `json:"name" validate:"required"`
	Version     string            `json:"version" validate:"required"`
	Description string            `json:"description"`
	Endpoint    string            `json:"endpoint" validate:"required,url"`
	Methods     []PluginMethod    `json:"methods" validate:"required,min=1"`
	Metadata    map[string]string `json:"metadata"`
}

// RegisterResponse 插件注册响应
type RegisterResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ExecuteRequest 执行请求
type ExecuteRequest struct {
	PluginID string            `json:"plugin_id" validate:"required"`
	Method   string            `json:"method" validate:"required"`
	Params   map[string]interface{} `json:"params"`
}

// ExecuteResponse 执行响应
type ExecuteResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	Error   string      `json:"error,omitempty"`
	Latency int64       `json:"latency"` // ms
}

// PluginListRequest 插件列表请求
type PluginListRequest struct {
	Keyword string `query:"keyword"`
	Status  string `query:"status"`
	Page    int    `query:"page" default:"1"`
	Limit   int    `query:"limit" default:"20"`
}

// PluginListResponse 插件列表响应
type PluginListResponse struct {
	Total   int      `json:"total"`
	Plugins []Plugin `json:"plugins"`
}
