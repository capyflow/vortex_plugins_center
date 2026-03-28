package main

import (
	"context"
	"fmt"
	"log"

	"plugin-platform/sdk/go"
)

// ImageResizerPlugin 图片缩放插件
func main() {
	// 创建插件
	plugin := sdk.NewPlugin("image-resizer", "1.0.0").
		SetEndpoint("http://localhost:8004").
		SetSummary("图片尺寸调整工具，支持按比例缩放和指定尺寸").
		SetOutputs("object").
		AddFileParam("image", true, ".jpg,.jpeg,.png").
		AddParam("width", sdk.ParamDef{
			Type:     "number",
			Required: false,
			Min:      1,
			Max:      4096,
			Default:  800,
		}).
		AddParam("height", sdk.ParamDef{
			Type:     "number",
			Required: false,
			Min:      1,
			Max:      4096,
			Default:  600,
		}).
		AddBooleanParam("keep_ratio", true).
		AddSelectParam("quality", []string{"low", "medium", "high"}, "medium").
		SetMetadata("author", "vortex").
		SetMetadata("license", "MIT")

	// 注册到平台
	ctx := context.Background()
	resp, err := plugin.Register(ctx, "http://localhost:8080")
	if err != nil {
		log.Fatalf("Failed to register: %v", err)
	}
	log.Printf("Registered successfully: %s", resp.ID)

	// 启动服务
	log.Println("Starting image resizer plugin on port 8004...")
	if err := plugin.Serve("8004", handleResize); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// handleResize 处理图片缩放
func handleResize(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	image, _ := params["image"].(string)
	width, _ := params["width"].(float64)
	height, _ := params["height"].(float64)
	keepRatio, _ := params["keep_ratio"].(bool)
	quality, _ := params["quality"].(string)

	if image == "" {
		return nil, fmt.Errorf("image is required")
	}

	// 实际实现中会调用图像处理库
	// 这里返回模拟结果
	return map[string]interface{}{
		"original": map[string]interface{}{
			"name":   image,
			"width":  1920,
			"height": 1080,
			"size":   2048000,
		},
		"resized": map[string]interface{}{
			"url":    fmt.Sprintf("https://cdn.example.com/resized/%s", image),
			"width":  int(width),
			"height": int(height),
			"size":   512000,
		},
		"keep_ratio": keepRatio,
		"quality":    quality,
	}, nil
}
