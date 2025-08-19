package models

import "time"

// MCPServer MCP服务器配置结构体
type MCPServer struct {
	Name        string            `json:"name" binding:"required"`    // 服务器名称
	Command     string            `json:"command" binding:"required"` // 执行命令
	Args        []string          `json:"args,omitempty"`             // 命令参数
	Env         map[string]string `json:"env,omitempty"`              // 环境变量
	Disabled    bool              `json:"disabled,omitempty"`         // 是否禁用
	AutoApprove []string          `json:"autoApprove,omitempty"`      // 自动批准的操作
}

// MCPRequest MCP请求结构体
type MCPRequest struct {
	Model      string                 `json:"model" binding:"required"`      // 模型名称
	MCPServers []string               `json:"mcpServers" binding:"required"` // MCP服务器名称列表
	Query      string                 `json:"query,omitempty"`               // 查询内容（用于工具调用）
	Tool       string                 `json:"tool,omitempty"`                // 要调用的工具名称
	Params     map[string]interface{} `json:"params,omitempty"`              // 工具调用参数
}

// MCPResponse MCP响应结构体
type MCPResponse struct {
	Model      string               `json:"model"`                 // 模型名称
	MCPServers map[string]MCPServer `json:"mcpServers"`            // 可用的MCP服务器配置
	Timestamp  time.Time            `json:"timestamp"`             // 响应时间戳
	ToolResult interface{}          `json:"tool_result,omitempty"` // 工具调用结果
	Query      string               `json:"query,omitempty"`       // 查询内容
	Tool       string               `json:"tool,omitempty"`        // 调用的工具
}

// GetDefaultMCPServers 获取默认的MCP服务器配置
func GetDefaultMCPServers() map[string]MCPServer {
	return map[string]MCPServer{
		"Fetch": {
			Name:    "Fetch",
			Command: "docker",
			Args: []string{
				"run",
				"-i",
				"--rm",
				"mcp/fetch",
			},
			Disabled:    false,
			AutoApprove: []string{},
		},
		"阿里云百炼_联网搜索": {
			Name:    "阿里云百炼_联网搜索",
			Command: "npx",
			Args: []string{
				"mcp-remote",
				"https://dashscope.aliyuncs.com/api/v1/mcps/WebSearch/sse",
				"--header",
				"Authorization:${AUTH_HEADER}",
			},
			Env: map[string]string{
				"AUTH_HEADER": "Bearer ${QWEN_API_KEY}",
			},
			Disabled:    false,
			AutoApprove: []string{},
		},
	}
}
