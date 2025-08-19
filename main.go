package main

import (
	"log"

	"github.com/aimmetal-tech/wistrans-backend/api"
	"github.com/aimmetal-tech/wistrans-backend/db"
	"github.com/aimmetal-tech/wistrans-backend/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	err := db.InitDB()
	if err != nil {
		log.Fatal("数据库初始化失败: ", err)
	}
	defer db.DB.Close()

	// 创建会话存储实例
	sessionStore := store.NewSessionStore(db.DB)

	// 创建API处理函数实例
	handlers, err := api.NewHandlers(sessionStore)
	if err != nil {
		log.Fatal("API处理器初始化失败: ", err)
	}

	// 设置Gin路由
	app := gin.Default()

	// 健康检查接口
	app.GET("/health", handlers.HealthCheck)

	// 会话相关接口
	app.GET("/conversations", handlers.CreateConversation)             // 创建新会话
	app.PATCH("/conversations/:id", handlers.UpdateConversation)       // 更新会话标题
	app.GET("/conversations/detail", handlers.GetConversationDetail)   // 获取会话详情
	app.GET("/conversations/history", handlers.GetConversationHistory) // 获取会话历史记录
	app.GET("/conversations/stream", handlers.StreamConversation)      // 流式对话接口

	// 翻译接口
	app.POST("/translate", handlers.Translate) // 网页翻译接口

	// MCP接口
	app.POST("/mcp", handlers.MCP) // MCP服务接口

	// Fetch接口
	app.POST("/fetch", handlers.Fetch) // 网页内容抓取接口

	// Web Search接口
	app.POST("/web-search", handlers.WebSearch) // 联网搜索接口

	// 启动服务
	log.Println("服务器启动在端口 8080")
	err = app.Run(":8080")
	if err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}
