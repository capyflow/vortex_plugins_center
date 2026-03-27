package main

import (
	"context"
	"flag"
	"log"
	"plugin-platform/internal/center"
	"plugin-platform/internal/gateway"
	"plugin-platform/internal/registry"
	"plugin-platform/internal/router"

	"github.com/capyflow/allspark-go/ds"
	"github.com/capyflow/allspark-go/logx"
)

func main() {
	var configPath = flag.String("c", "", "-c=/path/to/config.toml")
	flag.Parse()

	ctx := context.Background()

	// 初始化数据库
	mongoDB, err := ds.InitMongoDB(ctx, "plugin-platform")
	if err != nil {
		logx.Fatalf("Failed to connect MongoDB: %v", err)
	}

	redisClient, err := ds.InitRedis(ctx, 0)
	if err != nil {
		logx.Fatalf("Failed to connect Redis: %v", err)
	}

	// 初始化核心组件
	reg := registry.New(mongoDB, redisClient)
	rt := router.New(redisClient)
	pc := center.New(reg, rt)

	// 启动网关服务
	gw := gateway.New(pc)
	
	logx.Infof("Plugin Platform starting on port 8080")
	if err := gw.Start(":8080"); err != nil {
		logx.Fatalf("Failed to start gateway: %v", err)
	}
}
