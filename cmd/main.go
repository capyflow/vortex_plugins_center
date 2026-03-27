package main

import (
	"context"
	"flag"

	"plugin-platform/conf"
	"plugin-platform/internal/center"
	"plugin-platform/internal/gateway"
	"plugin-platform/internal/registry"
	"plugin-platform/internal/router"

	"github.com/capyflow/allspark-go/ds"
	"github.com/capyflow/allspark-go/logx"
	"github.com/capyflow/allspark-go/system"
)

func main() {
	var configPath = flag.String("c", "conf/config.toml", "-c=/path/to/config.toml")
	flag.Parse()

	ctx := context.Background()

	// 加载配置
	cfg := conf.LoadConfig(*configPath)

	dsServer := ds.InitDatabaseServer(ctx, cfg.DBConfig, func(dbIdxs map[string]interface{}) {
		dbIdxs["registry"] = 1
	})
	rdb, ok := dsServer.GetRedis("registry")
	if !ok {
		panic("Failed to get Redis client for registry")
	}
	// 初始化核心组件
	reg := registry.New(rdb)
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

	system.GracefulShutdown(func(ctx context.Context) error {
		return nil
	})
}
