package center

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"plugin-platform/internal/registry"
	"plugin-platform/internal/router"
	"plugin-platform/pkg/models"

	"github.com/capyflow/allspark-go/logx"
)

// PluginCenter 插件中心
type PluginCenter struct {
	registry *registry.Registry
	router   *router.Router
	client   *http.Client
}

// New 创建插件中心
func New(reg *registry.Registry, rt *router.Router) *PluginCenter {
	return &PluginCenter{
		registry: reg,
		router:   rt,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RegisterPlugin 注册插件
func (pc *PluginCenter) RegisterPlugin(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error) {
	// 验证端点可访问
	if err := pc.healthCheck(ctx, req.Endpoint); err != nil {
		return nil, fmt.Errorf("plugin endpoint not accessible: %v", err)
	}

	// 创建插件对象
	plugin := &models.Plugin{
		Name:        req.Name,
		Version:     req.Version,
		Description: req.Description,
		Endpoint:    req.Endpoint,
		Methods:     req.Methods,
		Metadata:    req.Metadata,
		Health: models.PluginHealth{
			Status:    "healthy",
			Latency:   0,
			CheckedAt: time.Now(),
		},
	}

	// 注册到注册表
	if err := pc.registry.Register(ctx, plugin); err != nil {
		return nil, fmt.Errorf("failed to register plugin: %v", err)
	}

	// 注册路由
	if err := pc.router.Register(plugin); err != nil {
		return nil, fmt.Errorf("failed to register routes: %v", err)
	}

	return &models.RegisterResponse{
		ID:      plugin.ID,
		Status:  "active",
		Message: "Plugin registered successfully",
	}, nil
}

// UnregisterPlugin 注销插件
func (pc *PluginCenter) UnregisterPlugin(ctx context.Context, pluginID string) error {
	if err := pc.registry.Unregister(ctx, pluginID); err != nil {
		return err
	}
	pc.router.Unregister(pluginID)
	return nil
}

// Execute 执行插件方法
func (pc *PluginCenter) Execute(ctx context.Context, pluginID, methodName string, params map[string]interface{}) (*models.ExecuteResponse, error) {
	start := time.Now()

	// 获取插件信息
	plugin, err := pc.registry.Get(ctx, pluginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin: %v", err)
	}
	if plugin == nil {
		return nil, fmt.Errorf("plugin not found: %s", pluginID)
	}

	// 检查状态
	if plugin.Status != models.PluginStatusActive {
		return nil, fmt.Errorf("plugin is not active: %s", plugin.Status)
	}

	// 路由到方法
	endpoint, path, err := pc.router.Route(pluginID, methodName)
	if err != nil {
		return nil, err
	}

	// 调用插件
	result, err := pc.callPlugin(ctx, endpoint, path, params)

	latency := time.Since(start).Milliseconds()

	// 更新健康状态
	health := models.PluginHealth{
		Status:    "healthy",
		Latency:   latency,
		CheckedAt: time.Now(),
	}
	if err != nil {
		health.Status = "unhealthy"
	}
	pc.registry.UpdateHealth(ctx, pluginID, health)

	if err != nil {
		return &models.ExecuteResponse{
			Success: false,
			Error:   err.Error(),
			Latency: latency,
		}, nil
	}

	return &models.ExecuteResponse{
		Success: true,
		Result:  result,
		Latency: latency,
	}, nil
}

// ListPlugins 列出插件
func (pc *PluginCenter) ListPlugins(ctx context.Context, keyword, status string, page, limit int) (*models.PluginListResponse, error) {
	plugins, total, err := pc.registry.List(ctx, keyword, status, page, limit)
	if err != nil {
		return nil, err
	}

	return &models.PluginListResponse{
		Total:   int(total),
		Plugins: plugins,
	}, nil
}

// GetPlugin 获取插件详情
func (pc *PluginCenter) GetPlugin(ctx context.Context, pluginID string) (*models.Plugin, error) {
	return pc.registry.Get(ctx, pluginID)
}

// HealthCheck 健康检查
func (pc *PluginCenter) HealthCheck(ctx context.Context, pluginID string) (*models.PluginHealth, error) {
	plugin, err := pc.registry.Get(ctx, pluginID)
	if err != nil {
		return nil, err
	}
	if plugin == nil {
		return nil, fmt.Errorf("plugin not found")
	}

	start := time.Now()
	err = pc.healthCheck(ctx, plugin.Endpoint)
	latency := time.Since(start).Milliseconds()

	health := models.PluginHealth{
		Status:    "healthy",
		Latency:   latency,
		CheckedAt: time.Now(),
	}
	if err != nil {
		health.Status = "unhealthy"
	}

	pc.registry.UpdateHealth(ctx, pluginID, health)
	return &health, nil
}

// callPlugin 调用插件（使用标准协议）
func (pc *PluginCenter) callPlugin(ctx context.Context, endpoint, path string, params map[string]interface{}) (interface{}, error) {
	url := endpoint + "/invoke"

	// 构建标准请求
	requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
	reqBody := map[string]interface{}{
		"request_id": requestID,
		"method":     extractMethodFromPath(path),
		"timestamp":  time.Now().UnixMilli(),
		"params":     params,
		"context":    map[string]string{},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// 设置标准 Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("X-Plugin-Name", "unknown") // 需要从上层传入
	req.Header.Set("X-Method", extractMethodFromPath(path))
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))
	req.Header.Set("X-Plugin-Platform", "v1.0")

	resp, err := pc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call plugin: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result, nil
}

// extractMethodFromPath 从路径提取方法名
func extractMethodFromPath(path string) string {
	// /add -> add
	if len(path) > 0 && path[0] == '/' {
		return path[1:]
	}
	return path
}

// healthCheck 健康检查
func (pc *PluginCenter) healthCheck(ctx context.Context, endpoint string) error {
	// 尝试访问 /health 端点
	healthURL := endpoint + "/health"
	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return err
	}

	resp, err := pc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// StartHealthCheckLoop 启动健康检查循环
func (pc *PluginCenter) StartHealthCheckLoop() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			ctx := context.Background()
			// 获取所有活跃插件
			plugins, _, err := pc.registry.List(ctx, "", string(models.PluginStatusActive), 1, 1000)
			if err != nil {
				logx.Errorf("Failed to list plugins for health check: %v", err)
				continue
			}

			for _, plugin := range plugins {
				pc.HealthCheck(ctx, plugin.ID)
			}
		}
	}()
}
