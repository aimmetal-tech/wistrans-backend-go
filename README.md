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

2. 在 .env 文件中填写您要使用的模型API密钥:
   ```
   DEEPSEEK_API_KEY=your_deepseek_api_key
   OPENAI_API_KEY=your_openai_api_key
   KIMI_API_KEY=your_kimi_api_key
   QWEN_API_KEY=your_qwen_api_key
   DATABASE_URL=postgresql://username:password@localhost:5432/database_name?sslmode=disable
   ```

   至少需要填写一个API密钥和数据库URL。

## 数据库设置

系统会自动创建所需的表结构：
- `conversations` 表：存储会话信息
- `messages` 表：存储对话消息

## 运行项目

```bash
go run main.go
```

## API接口说明

详情见docs/api.md

### 健康检查
```
GET /health
```
检查服务运行状态

### 创建会话
```
GET /conversations
```
创建一个新的对话会话，返回会话ID

### 获取会话详情
```
GET /conversations/detail?id=会话ID
```
获取指定会话的详细信息

返回结果示例：
```json
{
  "id": "会话ID",
  "title": "会话标题",
  "model": "使用的模型",
  "service": "服务商",
  "created_at": "创建时间",
  "updated_at": "更新时间"
}
```

### 获取会话历史记录
```
GET /conversations/history?id=会话ID
```
获取指定会话的所有历史消息记录

返回结果示例：
```json
{
  "conversation_id": "会话ID",
  "messages": [
    {
      "id": 1,
      "conversation_id": "会话ID",
      "role": "user",
      "content": "用户消息内容",
      "created_at": "消息创建时间"
    },
    {
      "id": 2,
      "conversation_id": "会话ID",
      "role": "assistant",
      "content": "助手回复内容",
      "created_at": "消息创建时间"
    }
  ]
}
```

### 更新会话信息
```
PATCH /conversations/{id}
```
更新会话信息（如标题）

请求体：
```json
{
  "title": "会话标题"
}
```

### 流式对话
```
GET /conversations/stream?id=会话ID&input=用户输入内容&model=模型名称
```
与AI进行流式对话，使用Server-Sent Events (SSE) 返回结果

参数说明：
- `id` (必填): 会话ID
- `input` (必填): 用户输入内容
- `model` (可选): 模型名称，可以是`服务商/模型`格式（如`openai/gpt-4o`），也可以直接指定模型

支持的服务商和默认模型：
- Qwen: `qwen-turbo-latest`
- DeepSeek: `deepseek-chat`
- OpenAI: `gpt-4o`
- Kimi: `kimi-k2-0711-preview`

## API优先级规则

1. 程序优先从系统环境变量获取API密钥
2. 如果系统环境变量中没有配置，则从.env文件中获取
3. 至少需要配置一个API密钥

## 服务商与默认模型

| 服务商 | 默认模型 |
|--------|---------|
| Qwen | qwen-turbo-latest |
| DeepSeek | deepseek-chat |
| OpenAI | gpt-4o |
| Kimi | kimi-k2-0711-preview |

## 模型指定方式

支持多种模型指定方式：

1. 使用`服务商/模型`格式：
   - `qwen/qwen-turbo-latest`
   - `deepseek/deepseek-chat`
   - `openai/gpt-4o`
   - `kimi/kimi-k2-0711-preview`

2. 直接指定模型名称：
   - `qwen-turbo-latest` 
   - `deepseek-chat` 
   - `gpt-4o` 
   - `kimi-k2-0711-preview` 

## 会话管理

系统使用PostgreSQL数据库存储会话记录，每个会话都有唯一的ID。会话中保存了对话历史，以实现连续对话功能。

## 服务商识别

程序支持多种方式指定服务商，目前有：
- Qwen、通义千问、通义
- DeepSeek、深度求索
- OpenAI、Open AI、GPT
- Kimi、Moonshot、月之暗面