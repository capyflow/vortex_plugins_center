package router

import (
	"context"
	"fmt"
	"sync"
	"time"

	"plugin-platform/pkg/models"

	"github.com/redis/go-redis/v9"
)

// Router 插件路由器（简化版）
type Router struct {
	redis    *redis.Client
	routes   map[string]*RouteInfo // pluginID -> route info
	mu       sync.RWMutex
}

// RouteInfo 路由信息（简化后：不再需要 methods 映射）
type RouteInfo struct {
	PluginID  string
	Endpoint  string
	UpdatedAt time.Time
}

// New 创建路由器
func New(redisClient *redis.Client) *Router {
	r := &Router{
		redis:  redisClient,
		routes: make(map[string]*RouteInfo),
	}

	// 启动清理任务
	go r.cleanupLoop()

	return r
}

// Register 注册路由
func (r *Router) Register(plugin *models.Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routes[plugin.ID] = &RouteInfo{
		PluginID:  plugin.ID,
		Endpoint:  plugin.Endpoint,
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

// GetEndpoint 获取插件端点（简化后：不再需要 methodName）
func (r *Router) GetEndpoint(pluginID string) (endpoint string, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	route, exists := r.routes[pluginID]
	if !exists {
		return "", fmt.Errorf("plugin not found: %s", pluginID)
	}

	return route.Endpoint, nil
}

// GetAllRoutes 获取所有路由
func (r *Router) GetAllRoutes() map[string]*RouteInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]*RouteInfo)
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
