package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// CalculatorPlugin 计算器插件（简化版）
// 示例：展示如何实现一个单一功能的插件

func main() {
	// 插件配置
	port := "8001"
	platformURL := "http://localhost:8080/api/v1/plugins/register"

	// 启动 HTTP 服务
	go func() {
		mux := http.NewServeMux()

		// 健康检查端点
		mux.HandleFunc("/health", handleHealth)

		// 执行端点（简化后：单一 /execute 端点）
		mux.HandleFunc("/execute", handleExecute)

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

// registerToPlatform 注册到平台（简化版格式）
func registerToPlatform(platformURL, port string) error {
	reqBody := map[string]interface{}{
		"name":     "calculator",
		"version":  "1.0.0",
		"endpoint": fmt.Sprintf("http://localhost:%s", port),
		"summary":  "简单计算器，支持加减乘除",
		"params": map[string]interface{}{
			"a": map[string]interface{}{
				"type":     "number",
				"required": true,
				"description": "第一个数字",
			},
			"b": map[string]interface{}{
				"type":     "number",
				"required": true,
				"description": "第二个数字",
			},
			"operator": map[string]interface{}{
				"type":    "select",
				"options": []string{"+", "-", "*", "/"},
				"default": "+",
				"description": "运算符",
			},
		},
		"outputs": "number",
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

// handleExecute 执行计算（简化后：单一入口）
func handleExecute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		A        float64 `json:"a"`
		B        float64 `json:"b"`
		Operator string  `json:"operator"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 默认运算符
	if req.Operator == "" {
		req.Operator = "+"
	}

	var result float64
	switch req.Operator {
	case "+":
		result = req.A + req.B
	case "-":
		result = req.A - req.B
	case "*":
		result = req.A * req.B
	case "/":
		if req.B == 0 {
			http.Error(w, "division by zero", http.StatusBadRequest)
			return
		}
		result = req.A / req.B
	default:
		http.Error(w, "invalid operator", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}
