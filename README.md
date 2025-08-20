# 智慧译-后端

智慧译是一个集成了多种大语言模型的翻译和对话系统。

## 功能特性

- 支持多种大语言模型 (Qwen, DeepSeek, Kimi, OpenAI)
- 连续对话功能
- 环境变量配置管理
- 会话管理
- 多服务商切换支持
- 基于Gin框架的RESTful API
- PostgreSQL数据库存储会话数据
- SSE流式响应支持
- 实时大模型API调用

## 环境配置

1. 复制 [.env.example](file:///e%3A/code/Project/wistrans-backend/.env.example) 文件并重命名为 .env:
   ```bash
   cp .env.example .env
   ```

2. 在 .env 文件中填写API密钥:

阿里云百炼API、Redis和PostgreSQL，以及FireCrawl API都要填写。

PostgreSQL必须创建名为wistrans的数据库

## 运行项目

```bash
go run main.go
```
