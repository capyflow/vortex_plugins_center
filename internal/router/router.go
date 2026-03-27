package router

import (
	"fmt"
	"sync"
	"time"

	"plugin-platform/pkg/models"

	"github.com/redis/go-redis/v9"
)

// Router 插件路由器
type Router struct {
	redis  *redis.Client
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
func New(redisClient *redis.Client) *Router {
	r := &Router{
		redis:  redisClient,
		routes: make(map[string]*RouteTable),
	}

	// 启动清理任务
	go r.cleanupLoop()

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

// cleanupLoop 清理过期路由
func (r *Router) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		now := time.Now()
		for id, route := range r.routes {
			if now.Sub(route.UpdatedAt) > 10*time.Minute {
				delete(r.routes, id)
			}
		}
		r.mu.Unlock()
	}
}
