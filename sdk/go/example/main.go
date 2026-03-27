package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	plugin "plugin-platform/sdk/go"
)

func main() {
	// 创建 SDK 实例
	sdk := plugin.New(plugin.RegisterOptions{
		PlatformURL: "http://localhost:8080",
		Name:        "advanced-calculator",
		Version:     "2.0.0",
		Description: "高级计算器插件，支持复杂数学运算",
		Endpoint:    "http://localhost:9001",
		DocURL:      "https://docs.example.com/calculator",
	})

	// 添加 add 方法
	sdk.AddMethod(
		plugin.MethodOptions{
			Name:        "add",
			Description: "两个数相加",
			Path:        "/add",
			HTTPMethod:  "POST",
		},
		[]plugin.Parameter{
			plugin.NumberParam("a", "第一个加数", true),
			plugin.NumberParam("b", "第二个加数", true),
		},
		plugin.ReturnInfo{
			Type:        "number",
			Description: "两数之和",
		},
	)

	// 添加 subtract 方法
	sdk.AddMethod(
		plugin.MethodOptions{
			Name:        "subtract",
			Description: "两个数相减",
			Path:        "/subtract",
			HTTPMethod:  "POST",
		},
		[]plugin.Parameter{
			plugin.NumberParam("a", "被减数", true),
			plugin.NumberParam("b", "减数", true),
		},
		plugin.ReturnInfo{
			Type:        "number",
			Description: "两数之差",
		},
	)

	// 添加 multiply 方法
	sdk.AddMethod(
		plugin.MethodOptions{
			Name:        "multiply",
			Description: "两个数相乘",
			Path:        "/multiply",
			HTTPMethod:  "POST",
		},
		[]plugin.Parameter{
			plugin.NumberParam("a", "第一个乘数", true),
			plugin.NumberParam("b", "第二个乘数", true),
		},
		plugin.ReturnInfo{
			Type:        "number",
			Description: "两数之积",
		},
	)

	// 添加 divide 方法
	sdk.AddMethod(
		plugin.MethodOptions{
			Name:        "divide",
			Description: "两个数相除",
			Path:        "/divide",
			HTTPMethod:  "POST",
		},
		[]plugin.Parameter{
			plugin.NumberParam("a", "被除数", true),
			plugin.NumberParam("b", "除数", true),
		},
		plugin.ReturnInfo{
			Type:        "number",
			Description: "两数之商",
		},
	)

	// 添加复杂计算方法（带对象参数）
	sdk.AddMethod(
		plugin.MethodOptions{
			Name:        "calculate",
			Description: "复杂计算，支持表达式",
			Path:        "/calculate",
			HTTPMethod:  "POST",
		},
		[]plugin.Parameter{
			plugin.StringParam("expression", "数学表达式，如: 2 + 3 * 4", true),
			plugin.ObjectParam("variables", "变量定义", false),
		},
		plugin.ReturnInfo{
			Type:        "number",
			Description: "计算结果",
		},
	)

	// 设置元数据
	sdk.SetMetadata("author", "Vortex Team").
		SetMetadata("license", "MIT").
		SetMetadata("category", "math")

	// 启动 HTTP 服务
	go startServer()

	// 等待服务启动
	fmt.Println("等待服务启动...")
	select {}

	// 注册到平台
	ctx := context.Background()
	resp, err := sdk.Register(ctx)
	if err != nil {
		log.Fatalf("注册失败: %v", err)
	}

	log.Printf("插件注册成功!")
	log.Printf("ID: %s", resp.ID)
	log.Printf("Status: %s", resp.Status)
	log.Printf("Message: %s", resp.Message)

	// 保持运行
	select {}
}

// startServer 启动 HTTP 服务
func startServer() {
	mux := http.NewServeMux()

	// 健康检查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// 加法
	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ A, B float64 }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]float64{"result": req.A + req.B})
	})

	// 减法
	mux.HandleFunc("/subtract", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ A, B float64 }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]float64{"result": req.A - req.B})
	})

	// 乘法
	mux.HandleFunc("/multiply", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ A, B float64 }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]float64{"result": req.A * req.B})
	})

	// 除法
	mux.HandleFunc("/divide", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ A, B float64 }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.B == 0 {
			http.Error(w, "除数不能为零", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]float64{"result": req.A / req.B})
	})

	// 复杂计算
	mux.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Expression string                 `json:"expression"`
			Variables  map[string]interface{} `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// 简化处理，实际应该解析表达式
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":     42,
			"expression": req.Expression,
		})
	})

	log.Println("插件 HTTP 服务启动在 :9001")
	if err := http.ListenAndServe(":9001", mux); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
