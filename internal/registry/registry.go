package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"plugin-platform/pkg/models"

	"github.com/capyflow/allspark-go/logx"
	"github.com/redis/go-redis/v9"
)

const (
	KeyPluginPrefix    = "plugin:id:"
	KeyPluginName      = "plugin:name:"
	KeyPluginList      = "plugin:list"
	KeyHealthFailCount = "plugin:health_fail:"
	MaxHealthFailures  = 3
)

// Registry 插件注册表（纯 Redis 实现）
type Registry struct {
	redis *redis.Client
}

// New 创建注册表
func New(redisClient *redis.Client) *Registry {
	return &Registry{
		redis: redisClient,
	}
}

// Register 注册插件
func (r *Registry) Register(ctx context.Context, plugin *models.Plugin) error {
	now := time.Now()

	// 检查是否已存在同名插件
	existing, err := r.GetByName(ctx, plugin.Name)
	if err != nil {
		return err
	}

	if existing != nil {
		// 更新现有插件
		plugin.ID = existing.ID
		plugin.CreatedAt = existing.CreatedAt
	} else {
		// 创建新插件
		plugin.ID = fmt.Sprintf("plugin-%d", now.UnixNano())
		plugin.CreatedAt = now
	}

	plugin.UpdatedAt = now
	plugin.Status = models.PluginStatusActive
	plugin.Health = models.PluginHealth{
		Status:    "healthy",
		Latency:   0,
		CheckedAt: now,
	}

	// 保存到 Redis
	data, err := json.Marshal(plugin)
	if err != nil {
		return fmt.Errorf("failed to marshal plugin: %v", err)
	}

	// 保存插件数据
	if err := r.redis.Set(ctx, KeyPluginPrefix+plugin.ID, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to save plugin: %v", err)
	}

	// 保存名称到 ID 的映射
	if err := r.redis.Set(ctx, KeyPluginName+plugin.Name, plugin.ID, 0).Err(); err != nil {
		return fmt.Errorf("failed to save plugin name mapping: %v", err)
	}

	// 添加到插件列表
	if err := r.redis.SAdd(ctx, KeyPluginList, plugin.ID).Err(); err != nil {
		return fmt.Errorf("failed to add plugin to list: %v", err)
	}

	// 重置健康失败计数
	r.redis.Del(ctx, KeyHealthFailCount+plugin.ID)

	if existing != nil {
		logx.Infof("Plugin updated: %s@%s", plugin.Name, plugin.Version)
	} else {
		logx.Infof("Plugin registered: %s@%s", plugin.Name, plugin.Version)
	}

	return nil
}

// Unregister 注销插件
func (r *Registry) Unregister(ctx context.Context, pluginID string) error {
	// 获取插件信息
	plugin, err := r.Get(ctx, pluginID)
	if err != nil {
		return err
	}
	if plugin == nil {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}

	// 更新状态为 inactive
	plugin.Status = models.PluginStatusInactive
	plugin.UpdatedAt = time.Now()

	data, err := json.Marshal(plugin)
	if err != nil {
		return fmt.Errorf("failed to marshal plugin: %v", err)
	}

	if err := r.redis.Set(ctx, KeyPluginPrefix+pluginID, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to update plugin: %v", err)
	}

	// 从名称映射中删除
	r.redis.Del(ctx, KeyPluginName+plugin.Name)

	// 从列表中移除
	r.redis.SRem(ctx, KeyPluginList, pluginID)

	// 清理健康失败计数
	r.redis.Del(ctx, KeyHealthFailCount+pluginID)

	logx.Infof("Plugin unregistered: %s", pluginID)
	return nil
}

// Delete 彻底删除插件
func (r *Registry) Delete(ctx context.Context, pluginID string) error {
	plugin, err := r.Get(ctx, pluginID)
	if err != nil {
		return err
	}

	if plugin != nil {
		r.redis.Del(ctx, KeyPluginName+plugin.Name)
	}

	r.redis.Del(ctx, KeyPluginPrefix+pluginID)
	r.redis.SRem(ctx, KeyPluginList, pluginID)
	r.redis.Del(ctx, KeyHealthFailCount+pluginID)

	logx.Infof("Plugin deleted: %s", pluginID)
	return nil
}

// Get 获取插件
func (r *Registry) Get(ctx context.Context, pluginID string) (*models.Plugin, error) {
	data, err := r.redis.Get(ctx, KeyPluginPrefix+pluginID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get plugin: %v", err)
	}

	var plugin models.Plugin
	if err := json.Unmarshal(data, &plugin); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin: %v", err)
	}

	return &plugin, nil
}

// GetByName 通过名称获取插件
func (r *Registry) GetByName(ctx context.Context, name string) (*models.Plugin, error) {
	pluginID, err := r.redis.Get(ctx, KeyPluginName+name).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get plugin id by name: %v", err)
	}

	return r.Get(ctx, pluginID)
}

// List 列出插件
func (r *Registry) List(ctx context.Context, keyword, status string, page, limit int) ([]models.Plugin, int64, error) {
	// 获取所有插件 ID
	ids, err := r.redis.SMembers(ctx, KeyPluginList).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get plugin list: %v", err)
	}

	var plugins []models.Plugin
	for _, id := range ids {
		plugin, err := r.Get(ctx, id)
		if err != nil {
			continue
		}
		if plugin == nil {
			continue
		}

		// 状态过滤
		if status != "" && string(plugin.Status) != status {
			continue
		}

		// 关键字过滤
		if keyword != "" {
			if !strings.Contains(strings.ToLower(plugin.Name), strings.ToLower(keyword)) &&
				!strings.Contains(strings.ToLower(plugin.Description), strings.ToLower(keyword)) {
				continue
			}
		}

		plugins = append(plugins, *plugin)
	}

	total := int64(len(plugins))

	// 分页
	start := (page - 1) * limit
	end := start + limit
	if start > int(total) {
		return []models.Plugin{}, total, nil
	}
	if end > int(total) {
		end = int(total)
	}

	return plugins[start:end], total, nil
}

// UpdateHealth 更新健康状态
func (r *Registry) UpdateHealth(ctx context.Context, pluginID string, health models.PluginHealth) error {
	plugin, err := r.Get(ctx, pluginID)
	if err != nil {
		return err
	}
	if plugin == nil {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}

	plugin.Health = health
	plugin.LastSeen = time.Now()
	plugin.UpdatedAt = time.Now()

	data, err := json.Marshal(plugin)
	if err != nil {
		return fmt.Errorf("failed to marshal plugin: %v", err)
	}

	return r.redis.Set(ctx, KeyPluginPrefix+pluginID, data, 0).Err()
}

// RecordHealthCheck 记录健康检查结果，返回是否需要清理
func (r *Registry) RecordHealthCheck(ctx context.Context, pluginID string, isHealthy bool) (shouldRemove bool, err error) {
	if isHealthy {
		// 健康时重置失败计数
		r.redis.Del(ctx, KeyHealthFailCount+pluginID)
		return false, nil
	}

	// 不健康时增加失败计数
	count, err := r.redis.Incr(ctx, KeyHealthFailCount+pluginID).Result()
	if err != nil {
		return false, err
	}

	// 设置过期时间（防止残留）
	r.redis.Expire(ctx, KeyHealthFailCount+pluginID, 10*time.Minute)

	logx.Warnf("Plugin %s health check failed, failure count: %d/%d", pluginID, count, MaxHealthFailures)

	return count >= MaxHealthFailures, nil
}

// GetAllActivePlugins 获取所有活跃插件
func (r *Registry) GetAllActivePlugins(ctx context.Context) ([]*models.Plugin, error) {
	ids, err := r.redis.SMembers(ctx, KeyPluginList).Result()
	if err != nil {
		return nil, err
	}

	var plugins []*models.Plugin
	for _, id := range ids {
		plugin, err := r.Get(ctx, id)
		if err != nil || plugin == nil {
			continue
		}
		if plugin.Status == models.PluginStatusActive {
			plugins = append(plugins, plugin)
		}
	}

	return plugins, nil
}
