package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Plugin SDK 提供插件注册和管理功能

type SDK struct {
	platformURL string
	client      *http.Client
	info        *PluginInfo
}

// PluginInfo 插件信息
type PluginInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Endpoint    string            `json:"endpoint"`
	Methods     []Method          `json:"methods"`
	Metadata    map[string]string `json:"metadata"`
	DocURL      string            `json:"doc_url"`
}

// Method 方法定义
type Method struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Path        string      `json:"path"`
	HTTPMethod  string      `json:"http_method"`
	Parameters  []Parameter `json:"parameters"`
	Returns     ReturnInfo  `json:"returns"`
}

// Parameter 参数定义
type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // string, number, boolean, object, array
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
}

// ReturnInfo 返回信息
type ReturnInfo struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// RegisterOptions 注册选项
type RegisterOptions struct {
	PlatformURL string // 平台地址，如 http://localhost:8080
	Name        string // 插件名称
	Version     string // 版本号
	Description string // 插件描述
	Endpoint    string // 插件HTTP端点
	DocURL      string // 文档地址
}

// MethodOptions 方法选项
type MethodOptions struct {
	Name        string
	Description string
	Path        string
	HTTPMethod  string // GET, POST, PUT, DELETE
}

// New 创建 SDK 实例
func New(opts RegisterOptions) *SDK {
	return &SDK{
		platformURL: opts.PlatformURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		info: &PluginInfo{
			Name:        opts.Name,
			Version:     opts.Version,
			Description: opts.Description,
			Endpoint:    opts.Endpoint,
			DocURL:      opts.DocURL,
			Methods:     []Method{},
			Metadata:    make(map[string]string),
		},
	}
}

// AddMethod 添加方法
func (s *SDK) AddMethod(opts MethodOptions, params []Parameter, returns ReturnInfo) *SDK {
	method := Method{
		Name:        opts.Name,
		Description: opts.Description,
		Path:        opts.Path,
		HTTPMethod:  opts.HTTPMethod,
		Parameters:  params,
		Returns:     returns,
	}
	s.info.Methods = append(s.info.Methods, method)
	return s
}

// AddStringParam 添加字符串参数
func StringParam(name, description string, required bool) Parameter {
	return Parameter{
		Name:        name,
		Type:        "string",
		Required:    required,
		Description: description,
	}
}

// AddNumberParam 添加数字参数
func NumberParam(name, description string, required bool) Parameter {
	return Parameter{
		Name:        name,
		Type:        "number",
		Required:    required,
		Description: description,
	}
}

// AddBoolParam 添加布尔参数
func BoolParam(name, description string, required bool) Parameter {
	return Parameter{
		Name:        name,
		Type:        "boolean",
		Required:    required,
		Description: description,
	}
}

// AddObjectParam 添加对象参数
func ObjectParam(name, description string, required bool) Parameter {
	return Parameter{
		Name:        name,
		Type:        "object",
		Required:    required,
		Description: description,
	}
}

// AddArrayParam 添加数组参数
func ArrayParam(name, description string, required bool) Parameter {
	return Parameter{
		Name:        name,
		Type:        "array",
		Required:    required,
		Description: description,
	}
}

// SetMetadata 设置元数据
func (s *SDK) SetMetadata(key, value string) *SDK {
	s.info.Metadata[key] = value
	return s
}

// Register 注册到平台
func (s *SDK) Register(ctx context.Context) (*RegisterResponse, error) {
	if len(s.info.Methods) == 0 {
		return nil, fmt.Errorf("at least one method is required")
	}

	url := s.platformURL + "/api/v1/plugins/register"
	
	jsonData, err := json.Marshal(s.info)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal plugin info: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registration failed with status %d", resp.StatusCode)
	}

	var result RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	s.info.ID = result.ID
	return &result, nil
}

// Unregister 从平台注销
func (s *SDK) Unregister(ctx context.Context) error {
	if s.info.ID == "" {
		return fmt.Errorf("plugin not registered")
	}

	url := fmt.Sprintf("%s/api/v1/plugins/%s", s.platformURL, s.info.ID)
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unregister failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetInfo 获取插件信息
func (s *SDK) GetInfo() *PluginInfo {
	return s.info
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
