package main

import (
	"context"
	"fmt"
	"log"

	"plugin-platform/sdk/go"
)

// FileUploaderPlugin 文件上传插件
func main() {
	// 创建插件
	plugin := sdk.NewPlugin("file-uploader", "1.0.0").
		SetEndpoint("http://localhost:8002").
		SetSummary("文件上传插件，支持图片压缩和水印").
		SetOutputs("object").
		AddFileParam("file", true, ".jpg,.jpeg,.png,.gif").
		AddBooleanParam("compress", true).
		AddBooleanParam("watermark", false).
		AddSelectParam("format", []string{"original", "webp", "png"}, "original").
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
	log.Println("Starting file uploader plugin on port 8002...")
	if err := plugin.Serve("8002", handleUpload); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// handleUpload 处理上传
func handleUpload(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 实际实现中这里会处理文件上传
	// 示例返回模拟结果
	fileName, _ := params["file"].(string)
	compress, _ := params["compress"].(bool)
	watermark, _ := params["watermark"].(bool)
	format, _ := params["format"].(string)

	return map[string]interface{}{
		"original_name": fileName,
		"url":           fmt.Sprintf("https://cdn.example.com/%s", fileName),
		"size":          1024000,
		"compressed":    compress,
		"watermarked":   watermark,
		"format":        format,
		"status":        "success",
	}, nil
}
