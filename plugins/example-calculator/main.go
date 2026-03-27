package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// CalculatorPlugin 计算器插件
// 示例：展示如何实现一个插件

func main() {
	// 插件配置
	port := "8001"
	platformURL := "http://localhost:19090/v1/api/plugins/register"

	// 启动 HTTP 服务
	go func() {
		mux := http.NewServeMux()

		// 健康检查端点
		mux.HandleFunc("/health", handleHealth)

		// 计算方法
		mux.HandleFunc("/add", handleAdd)
		mux.HandleFunc("/sub", handleSub)
		mux.HandleFunc("/mul", handleMul)
		mux.HandleFunc("/div", handleDiv)

		log.Printf("Calculator plugin starting on port %s", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待服务启动
	time.Sleep(1 * time.Second)

	// 注册到平台
	if err := registerToPlatform(platformURL, port); err != nil {
		log.Printf("Failed to register: %v", err)
	}

	// 保持运行
	select {}
}

// registerToPlatform 注册到平台
func registerToPlatform(platformURL, port string) error {
	reqBody := map[string]interface{}{
		"name":        "calculator",
		"version":     "1.0.0",
		"description": "Simple calculator plugin",
		"endpoint":    fmt.Sprintf("http://localhost:%s", port),
		"methods": []map[string]interface{}{
			{
				"name":        "add",
				"description": "Add two numbers",
				"path":        "/add",
				"method":      "POST",
				"parameters": []map[string]interface{}{
					{"name": "a", "type": "number", "required": true},
					{"name": "b", "type": "number", "required": true},
				},
				"returns": map[string]string{
					"type":        "number",
					"description": "Sum of a and b",
				},
			},
			{
				"name":        "sub",
				"description": "Subtract two numbers",
				"path":        "/sub",
				"method":      "POST",
				"parameters": []map[string]interface{}{
					{"name": "a", "type": "number", "required": true},
					{"name": "b", "type": "number", "required": true},
				},
				"returns": map[string]string{
					"type":        "number",
					"description": "Difference of a and b",
				},
			},
			{
				"name":        "mul",
				"description": "Multiply two numbers",
				"path":        "/mul",
				"method":      "POST",
				"parameters": []map[string]interface{}{
					{"name": "a", "type": "number", "required": true},
					{"name": "b", "type": "number", "required": true},
				},
				"returns": map[string]string{
					"type":        "number",
					"description": "Product of a and b",
				},
			},
			{
				"name":        "div",
				"description": "Divide two numbers",
				"path":        "/div",
				"method":      "POST",
				"parameters": []map[string]interface{}{
					{"name": "a", "type": "number", "required": true},
					{"name": "b", "type": "number", "required": true},
				},
				"returns": map[string]string{
					"type":        "number",
					"description": "Quotient of a and b",
				},
			},
		},
		"metadata": map[string]string{
			"author":  "vortex",
			"license": "MIT",
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(platformURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed with status %d", resp.StatusCode)
	}

	log.Println("Registered to platform successfully")
	return nil
}

// handleHealth 健康检查
func handleHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// handleAdd 加法
func handleAdd(w http.ResponseWriter, r *http.Request) {
	var req struct {
		A float64 `json:"a"`
		B float64 `json:"b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := req.A + req.B
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}

// handleSub 减法
func handleSub(w http.ResponseWriter, r *http.Request) {
	var req struct {
		A float64 `json:"a"`
		B float64 `json:"b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := req.A - req.B
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}

// handleMul 乘法
func handleMul(w http.ResponseWriter, r *http.Request) {
	var req struct {
		A float64 `json:"a"`
		B float64 `json:"b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := req.A * req.B
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}

// handleDiv 除法
func handleDiv(w http.ResponseWriter, r *http.Request) {
	var req struct {
		A float64 `json:"a"`
		B float64 `json:"b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.B == 0 {
		http.Error(w, "division by zero", http.StatusBadRequest)
		return
	}

	result := req.A / req.B
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}
