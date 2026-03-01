package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      any             `json:"id"`
}

type Response struct {
	JSONRPC string `json:"jsonrpc"`
	Result  any    `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
	ID      any    `json:"id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	// MCP 服务器通过标准输入、输出进行通信，所以需要一个扫描器来读取 stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		log.Printf("Received: %s", string(line))
		// 读取每一行，通常是一个 JSON-RPC 请求，并尝试解析
		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			log.Printf("Error unmarshaling request:%v", err)
			continue
		}
		// 根据请求方法，路由到不同处理函数
		switch req.Method {
		case "initialize":
			handleInitialize(req)
		case "tools/list":
			handleToolsList(req)
		case "tools/call":
			handleToolCall(req)
		case "notifications/initialized":
			log.Println("Received initialized notification")
			// 客户端发送的初始化完成通知，无需响应
			continue

		default:
			log.Printf("Unknown method: %s", req.Method)
			sendError(req.ID, -32061, "Method not found")

		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
}

// handleInitialize 负责向 claude code 自我介绍
func handleInitialize(req Request) {
	// 符合MCP协议的initialize响应
	initialize := map[string]any{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]any{
			"tools": map[string]any{}, // 生命支持工具能力

		},
		"serverInfo": map[string]any{
			"name":    "hello-server",
			"version": "1.0.0",
		},
	}
	sendResult(req.ID, initialize)
}

// handleToolsList 返回可用工具列表
func handleToolsList(req Request) {
	toolsListResult := map[string]any{
		"tools": []map[string]any{
			{
				"name":        "greet",
				"description": "A simple tool that returns a greeting.",
				"inputSchema": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name": map[string]any{
							"type":        "string",
							"description": "The name of the person to greet.",
						},
					},
					"required": []string{"name"},
				},
			},
		},
	}
	sendResult(req.ID, toolsListResult)
}

func handleToolCall(req Request) {
	var params map[string]any
	if err := json.Unmarshal(req.Params, &params); err != nil {
		sendError(req.ID, -32602, "Invalid params")
		return
	}
	toolName, _ := params["name"].(string)
	if toolName != "greet" {
		sendError(req.ID, -32601, "Tool not found")
		return
	}
	toolArguments, _ := params["arguments"].(map[string]any)
	name, _ := toolArguments["name"].(string)
	// 工具
	greeting := fmt.Sprintf("Hello,%s Welcome to the world of MCP in go.", name)
	// MCP期望的格式
	toolResult := map[string]any{
		"content": []map[string]any{
			{
				"type": "text",
				"text": greeting,
			},
		},
	}
	sendResult(req.ID, toolResult)
}

func sendResult(id any, result any) {
	resp := Response{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	sendJSON(resp)
}

func sendError(id any, code int, msg string) {
	resp := Response{
		JSONRPC: "2.0",
		Error: &Error{
			Code:    code,
			Message: msg,
		},
		ID: id,
	}
	sendJSON(resp)
}

func sendJSON(v any) {
	encoded, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling response:%v", err)
		return
	}
	// MCP协议要求每个JSON对象后都有一个换行符
	fmt.Println(string(encoded))
}
