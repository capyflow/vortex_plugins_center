package router

import (
	"fmt"
	"sync"
	"time"

	"plugin-platform/pkg/models"
)

// Router 插件路由器
type Router struct {
	routes map[string]*RouteTable // pluginID -> route table
	mu     sync.RWMutex
}

// RouteTable 路由表
type RouteTable struct {
	PluginID  string
	Endpoint  string
	Methods   map[string]string // method name -> http path
	UpdatedAt time.Time
}

// New 创建路由器
func New() *Router {
	r := &Router{
		routes: make(map[string]*RouteTable),
	}
	return r
}

// Register 注册路由
func (r *Router) Register(plugin *models.Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	methods := make(map[string]string)
	for _, m := range plugin.Methods {
		methods[m.Name] = m.Path
	}

	r.routes[plugin.ID] = &RouteTable{
		PluginID:  plugin.ID,
		Endpoint:  plugin.Endpoint,
		Methods:   methods,
		UpdatedAt: time.Now(),
	}

	return nil
}

// Unregister 注销路由
func (r *Router) Unregister(pluginID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.routes, pluginID)
}

// Route 路由请求
func (r *Router) Route(pluginID, methodName string) (endpoint, path string, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	route, exists := r.routes[pluginID]
	if !exists {
		return "", "", fmt.Errorf("plugin not found: %s", pluginID)
	}

	path, exists = route.Methods[methodName]
	if !exists {
		return "", "", fmt.Errorf("method not found: %s.%s", pluginID, methodName)
	}

	return route.Endpoint, path, nil
}

// RouteByName 通过名称路由
func (r *Router) RouteByName(pluginName, methodName string) (endpoint, path string, err error) {
	// TODO: 从 Redis 查询 pluginName -> pluginID 映射
	return "", "", fmt.Errorf("not implemented")
}

// GetAllRoutes 获取所有路由
func (r *Router) GetAllRoutes() map[string]*RouteTable {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]*RouteTable)
	for k, v := range r.routes {
		result[k] = v
	}
	return result
}

// UpdateTimestamp 更新路由时间戳
func (r *Router) UpdateTimestamp(pluginID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if route, exists := r.routes[pluginID]; exists {
		route.UpdatedAt = time.Now()
	}
}
