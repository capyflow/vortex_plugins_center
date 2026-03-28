// Package sdk 提供插件开发 SDK
package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Plugin SDK 插件基础结构
type Plugin struct {
	Name     string
	Version  string
	Endpoint string
	Summary  string
	Params   map[string]ParamDef
	Outputs  string
	Metadata map[string]string

	platformURL string
	client      *http.Client
}

// ParamDef 参数定义
type ParamDef struct {
	Type        string   `json:"type"`
	Required    bool     `json:"required,omitempty"`
	Default     any      `json:"default,omitempty"`
	Description string   `json:"description,omitempty"`
	MaxLength   int      `json:"maxLength,omitempty"`
	Min         float64  `json:"min,omitempty"`
	Max         float64  `json:"max,omitempty"`
	Step        float64  `json:"step,omitempty"`
	Accept      string   `json:"accept,omitempty"`
	Multiple    bool     `json:"multiple,omitempty"`
	MaxSize     int64    `json:"maxSize,omitempty"`
	Options     []string `json:"options,omitempty"`
	Rows        int      `json:"rows,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
}

// NewPlugin 创建新插件
func NewPlugin(name, version string) *Plugin {
	return &Plugin{
		Name:     name,
		Version:  version,
		Params:   make(map[string]ParamDef),
		Metadata: make(map[string]string),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetEndpoint 设置端点
func (p *Plugin) SetEndpoint(endpoint string) *Plugin {
	p.Endpoint = endpoint
	return p
}

// SetSummary 设置描述
func (p *Plugin) SetSummary(summary string) *Plugin {
	p.Summary = summary
	return p
}

// SetOutputs 设置输出类型
func (p *Plugin) SetOutputs(outputs string) *Plugin {
	p.Outputs = outputs
	return p
}

// AddParam 添加参数
func (p *Plugin) AddParam(name string, def ParamDef) *Plugin {
	p.Params[name] = def
	return p
}

// AddStringParam 添加字符串参数
func (p *Plugin) AddStringParam(name string, required bool) *Plugin {
	return p.AddParam(name, ParamDef{Type: "string", Required: required})
}

// AddNumberParam 添加数字参数
func (p *Plugin) AddNumberParam(name string, required bool) *Plugin {
	return p.AddParam(name, ParamDef{Type: "number", Required: required})
}

// AddFileParam 添加文件参数
func (p *Plugin) AddFileParam(name string, required bool, accept string) *Plugin {
	return p.AddParam(name, ParamDef{
		Type:     "file",
		Required: required,
		Accept:   accept,
	})
}

// AddSelectParam 添加选择参数
func (p *Plugin) AddSelectParam(name string, options []string, defaultVal string) *Plugin {
	return p.AddParam(name, ParamDef{
		Type:    "select",
		Options: options,
		Default: defaultVal,
	})
}

// AddBooleanParam 添加布尔参数
func (p *Plugin) AddBooleanParam(name string, defaultVal bool) *Plugin {
	return p.AddParam(name, ParamDef{
		Type:    "boolean",
		Default: defaultVal,
	})
}

// SetMetadata 设置元数据
func (p *Plugin) SetMetadata(key, value string) *Plugin {
	p.Metadata[key] = value
	return p
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	Endpoint string            `json:"endpoint"`
	Summary  string            `json:"summary,omitempty"`
	Params   map[string]ParamDef `json:"params,omitempty"`
	Outputs  string            `json:"outputs,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// RegisterResponse 注册响应结构
type RegisterResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Register 注册到平台
func (p *Plugin) Register(ctx context.Context, platformURL string) (*RegisterResponse, error) {
	p.platformURL = platformURL

	reqBody := RegisterRequest{
		Name:     p.Name,
		Version:  p.Version,
		Endpoint: p.Endpoint,
		Summary:  p.Summary,
		Params:   p.Params,
		Outputs:  p.Outputs,
		Metadata: p.Metadata,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", platformURL+"/api/v1/plugins/register", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registration failed with status %d", resp.StatusCode)
	}

	var result RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}

// Handler 插件处理器接口
type Handler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// Serve 启动 HTTP 服务
func (p *Plugin) Serve(port string, handler Handler) error {
	mux := http.NewServeMux()

	// 健康检查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// 执行端点
	mux.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		var params map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := handler(r.Context(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": result,
		})
	})

	return http.ListenAndServe("+"+port, mux)
}
