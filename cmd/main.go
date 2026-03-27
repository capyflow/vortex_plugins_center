package main

import (
	"context"
	"flag"
	"fmt"

	"plugin-platform/conf"
	"plugin-platform/internal/center"
	"plugin-platform/internal/gateway"
	"plugin-platform/internal/registry"
	"plugin-platform/internal/router"

	"github.com/capyflow/allspark-go/logx"
	"github.com/redis/go-redis/v9"
)

func main() {
	var configPath = flag.String("c", "conf/config.toml", "-c=/path/to/config.toml")
	flag.Parse()

	ctx := context.Background()

	// 加载配置
	cfg := conf.LoadConfig(*configPath)

	// 初始化 Redis 客户端
	redisClient := initRedis(cfg)

	// 初始化核心组件
	reg := registry.New(redisClient)
	rt := router.New()
	pc := center.New(reg, rt)

	// 启动健康检查循环（每分钟检查一次，3次失败清理）
	pc.StartHealthCheckLoop()

	// 启动网关服务
	gw := gateway.New(pc)

	port := cfg.Port
	if port == 0 {
		port = 8080
	}

	logx.Infof("Plugin Platform starting on port %d", port)
	if err := gw.Start(ctx, port); err != nil {
		logx.Fatalf("Failed to start gateway: %v", err)
	}
}

// initRedis 初始化 Redis 客户端
func initRedis(cfg *conf.CenterConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.DBConfig.Redis.Host, cfg.DBConfig.Redis.Port),
		Username: cfg.DBConfig.Redis.Username,
		Password: cfg.DBConfig.Redis.Password,
		DB:       0,
	})

	// 测试连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	logx.Infof("Connected to Redis at %s:%d", cfg.DBConfig.Redis.Host, cfg.DBConfig.Redis.Port)
	return rdb
}
