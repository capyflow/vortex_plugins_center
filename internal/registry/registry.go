package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"plugin-platform/pkg/models"

	"github.com/capyflow/allspark-go/logx"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionName = "plugins"
	RedisKeyPrefix = "plugin:"
)

// Registry 插件注册表
type Registry struct {
	mongo   *mongo.Database
	redis   *redis.Client
}

// New 创建注册表
func New(mongoDB *mongo.Database, redisClient *redis.Client) *Registry {
	return &Registry{
		mongo: mongoDB,
		redis: redisClient,
	}
}

// Register 注册插件
func (r *Registry) Register(ctx context.Context, plugin *models.Plugin) error {
	// 检查是否已存在
	existing, err := r.GetByName(ctx, plugin.Name)
	if err != nil {
		return err
	}

	now := time.Now()
	if existing != nil {
		// 更新现有插件
		plugin.ID = existing.ID
		plugin.CreatedAt = existing.CreatedAt
		plugin.UpdatedAt = now
		
		filter := bson.M{"_id": plugin.ID}
		update := bson.M{"$set": plugin}
		
		_, err := r.mongo.Collection(CollectionName).UpdateOne(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("failed to update plugin: %v", err)
		}
		
		logx.Infof("Plugin updated: %s@%s", plugin.Name, plugin.Version)
	} else {
		// 创建新插件
		plugin.ID = fmt.Sprintf("plugin-%d", time.Now().UnixNano())
		plugin.CreatedAt = now
		plugin.UpdatedAt = now
		plugin.Status = models.PluginStatusActive
		
		_, err := r.mongo.Collection(CollectionName).InsertOne(ctx, plugin)
		if err != nil {
			return fmt.Errorf("failed to create plugin: %v", err)
		}
		
		logx.Infof("Plugin registered: %s@%s", plugin.Name, plugin.Version)
	}

	// 缓存到 Redis
	if err := r.cachePlugin(ctx, plugin); err != nil {
		logx.Warnf("Failed to cache plugin: %v", err)
	}

	return nil
}

// Unregister 注销插件
func (r *Registry) Unregister(ctx context.Context, pluginID string) error {
	filter := bson.M{"_id": pluginID}
	update := bson.M{
		"$set": bson.M{
			"status":     models.PluginStatusInactive,
			"updated_at": time.Now(),
		},
	}

	_, err := r.mongo.Collection(CollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to unregister plugin: %v", err)
	}

	// 删除缓存
	r.redis.Del(ctx, RedisKeyPrefix+pluginID)
	r.redis.Del(ctx, RedisKeyPrefix+"name:"+pluginID)

	logx.Infof("Plugin unregistered: %s", pluginID)
	return nil
}

// Get 获取插件
func (r *Registry) Get(ctx context.Context, pluginID string) (*models.Plugin, error) {
	// 先查缓存
	cached, err := r.getCachedPlugin(ctx, pluginID)
	if err == nil && cached != nil {
		return cached, nil
	}

	// 查数据库
	filter := bson.M{"_id": pluginID}
	var plugin models.Plugin
	err = r.mongo.Collection(CollectionName).FindOne(ctx, filter).Decode(&plugin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get plugin: %v", err)
	}

	// 缓存
	r.cachePlugin(ctx, &plugin)

	return &plugin, nil
}

// GetByName 通过名称获取插件
func (r *Registry) GetByName(ctx context.Context, name string) (*models.Plugin, error) {
	filter := bson.M{"name": name, "status": models.PluginStatusActive}
	opts := options.FindOne().SetSort(bson.M{"updated_at": -1})

	var plugin models.Plugin
	err := r.mongo.Collection(CollectionName).FindOne(ctx, filter, opts).Decode(&plugin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get plugin by name: %v", err)
	}

	return &plugin, nil
}

// List 列出插件
func (r *Registry) List(ctx context.Context, keyword string, status string, page, limit int) ([]models.Plugin, int64, error) {
	filter := bson.M{}
	
	if status != "" {
		filter["status"] = status
	}
	
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"description": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	// 获取总数
	total, err := r.mongo.Collection(CollectionName).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count plugins: %v", err)
	}

	// 分页查询
	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"updated_at": -1})

	cursor, err := r.mongo.Collection(CollectionName).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list plugins: %v", err)
	}
	defer cursor.Close(ctx)

	var plugins []models.Plugin
	if err := cursor.All(ctx, &plugins); err != nil {
		return nil, 0, fmt.Errorf("failed to decode plugins: %v", err)
	}

	return plugins, total, nil
}

// UpdateHealth 更新健康状态
func (r *Registry) UpdateHealth(ctx context.Context, pluginID string, health models.PluginHealth) error {
	filter := bson.M{"_id": pluginID}
	update := bson.M{
		"$set": bson.M{
			"health":     health,
			"last_seen":  time.Now(),
			"updated_at": time.Now(),
		},
	}

	_, err := r.mongo.Collection(CollectionName).UpdateOne(ctx, filter, update)
	return err
}

// cachePlugin 缓存插件
func (r *Registry) cachePlugin(ctx context.Context, plugin *models.Plugin) error {
	data, err := json.Marshal(plugin)
	if err != nil {
		return err
	}

	// 缓存插件信息
	if err := r.redis.Set(ctx, RedisKeyPrefix+plugin.ID, data, 5*time.Minute).Err(); err != nil {
		return err
	}

	// 缓存名称映射
	return r.redis.Set(ctx, RedisKeyPrefix+"name:"+plugin.Name, plugin.ID, 5*time.Minute).Err()
}

// getCachedPlugin 获取缓存的插件
func (r *Registry) getCachedPlugin(ctx context.Context, pluginID string) (*models.Plugin, error) {
	data, err := r.redis.Get(ctx, RedisKeyPrefix+pluginID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var plugin models.Plugin
	if err := json.Unmarshal(data, &plugin); err != nil {
		return nil, err
	}

	return &plugin, nil
}
