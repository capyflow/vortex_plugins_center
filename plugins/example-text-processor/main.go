package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"plugin-platform/sdk/go"
)

// TextProcessorPlugin 文本处理插件
func main() {
	// 创建插件
	plugin := sdk.NewPlugin("text-processor", "1.0.0").
		SetEndpoint("http://localhost:8003").
		SetSummary("文本处理工具，支持大小写转换、统计、格式化").
		SetOutputs("object").
		AddParam("text", sdk.ParamDef{
			Type:        "textarea",
			Required:    true,
			Rows:        5,
			MaxLength:   10000,
			Placeholder: "请输入要处理的文本",
		}).
		AddSelectParam("operation", []string{"upper", "lower", "count", "reverse", "trim"}, "count").
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
	log.Println("Starting text processor plugin on port 8003...")
	if err := plugin.Serve("8003", handleProcess); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// handleProcess 处理文本
func handleProcess(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	text, _ := params["text"].(string)
	operation, _ := params["operation"].(string)

	if text == "" {
		return nil, fmt.Errorf("text is required")
	}

	switch operation {
	case "upper":
		return map[string]interface{}{
			"result": strings.ToUpper(text),
			"length": len(text),
		}, nil
	case "lower":
		return map[string]interface{}{
			"result": strings.ToLower(text),
			"length": len(text),
		}, nil
	case "count":
		words := strings.Fields(text)
		return map[string]interface{}{
			"char_count":  len(text),
			"word_count":  len(words),
			"line_count":  strings.Count(text, "\n") + 1,
		}, nil
	case "reverse":
		runes := []rune(text)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return map[string]interface{}{
			"result": string(runes),
		}, nil
	case "trim":
		return map[string]interface{}{
			"result": strings.TrimSpace(text),
		}, nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}
}
