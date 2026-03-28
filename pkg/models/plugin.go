package models

import (
	"time"
)

// Plugin 插件定义（简化版）
type Plugin struct {
	ID          string              `json:"id" bson:"_id"`
	Name        string              `json:"name" bson:"name"`
	Version     string              `json:"version" bson:"version"`
	Endpoint    string              `json:"endpoint" bson:"endpoint"` // HTTP 端点
	Summary     string              `json:"summary" bson:"summary"`
	Params      map[string]ParamDef `json:"params" bson:"params"`
	Outputs     string              `json:"outputs" bson:"outputs"`
	Metadata    map[string]string   `json:"metadata" bson:"metadata"`
	Status      PluginStatus        `json:"status" bson:"status"`
	Health      PluginHealth        `json:"health" bson:"health"`
	CreatedAt   time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" bson:"updated_at"`
	LastSeen    time.Time           `json:"last_seen" bson:"last_seen"`
}

// ParamDef 参数定义
type ParamDef struct {
	Type        string   `json:"type" bson:"type"` // string/number/file/select/boolean/textarea/password/date/object/array
	Required    bool     `json:"required,omitempty" bson:"required,omitempty"`
	Default     any      `json:"default,omitempty" bson:"default,omitempty"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`

	// 类型特定属性
	MaxLength   int      `json:"maxLength,omitempty" bson:"maxLength,omitempty"`   // string/textarea
	Min         float64  `json:"min,omitempty" bson:"min,omitempty"`             // number
	Max         float64  `json:"max,omitempty" bson:"max,omitempty"`             // number
	Step        float64  `json:"step,omitempty" bson:"step,omitempty"`           // number
	Accept      string   `json:"accept,omitempty" bson:"accept,omitempty"`       // file: ".jpg,.png"
	Multiple    bool     `json:"multiple,omitempty" bson:"multiple,omitempty"`   // file/select
	MaxSize     int64    `json:"maxSize,omitempty" bson:"maxSize,omitempty"`     // file: bytes
	Options     []string `json:"options,omitempty" bson:"options,omitempty"`     // select
	Rows        int      `json:"rows,omitempty" bson:"rows,omitempty"`           // textarea
	Placeholder string   `json:"placeholder,omitempty" bson:"placeholder,omitempty"` // string/textarea
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

// RegisterRequest 插件注册请求（简化版）
type RegisterRequest struct {
	Name     string              `json:"name" validate:"required"`
	Version  string              `json:"version" validate:"required"`
	Endpoint string              `json:"endpoint" validate:"required,url"`
	Summary  string              `json:"summary,omitempty"`
	Params   map[string]ParamDef `json:"params,omitempty"`
	Outputs  string              `json:"outputs,omitempty"`
	Metadata map[string]string   `json:"metadata,omitempty"`
}

// RegisterResponse 插件注册响应
type RegisterResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ExecuteRequest 执行请求
type ExecuteRequest struct {
	PluginID string                 `json:"plugin_id" validate:"required"`
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
