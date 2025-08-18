# API 接口文档

## 基础信息

- 服务器地址: `http://localhost:8080`
- 数据格式: JSON
- 字符编码: UTF-8

## 接口列表

### 1. 健康检查接口

#### 接口说明
检查服务运行状态

#### 接口地址
```
GET /health
```

#### 请求参数
无

#### 响应结果
```json
{
  "status": "ok",
  "message": "服务运行正常"
}
```

### 2. 创建会话接口

#### 接口说明
创建一个新的对话会话

#### 接口地址
```
GET /conversations
```

#### 请求参数
无

#### 响应结果
```json
{
  "id": "会话ID"
}
```

### 3. 获取会话详情接口

#### 接口说明
获取指定会话的详细信息

#### 接口地址
```
GET /conversations/detail
```

#### 请求参数
| 参数名 | 类型   | 必填 | 说明     |
| ------ | ------ | ---- | -------- |
| id     | string | 是   | 会话ID   |

#### 响应结果
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

#### 响应示例
```json
{
  "id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "title": "人工智能发展讨论",
  "model": "deepseek-chat",
  "service": "DeepSeek",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:05:00Z"
}
```

### 4. 获取会话历史记录接口

#### 接口说明
获取指定会话的所有历史消息记录

#### 接口地址
```
GET /conversations/history
```

#### 请求参数
| 参数名 | 类型   | 必填 | 说明     |
| ------ | ------ | ---- | -------- |
| id     | string | 是   | 会话ID   |

#### 响应结果
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

#### 响应示例
```json
{
  "conversation_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "messages": [
    {
      "id": 1,
      "conversation_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
      "role": "user",
      "content": "你好，请介绍一下人工智能的发展历程",
      "created_at": "2024-01-01T10:00:00Z"
    },
    {
      "id": 2,
      "conversation_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
      "role": "assistant",
      "content": "人工智能的发展可以分为几个阶段...",
      "created_at": "2024-01-01T10:00:30Z"
    }
  ]
}
```

### 5. 更新会话接口

#### 接口说明
更新会话信息（如标题）

#### 接口地址
```
PATCH /conversations/{id}
```

#### 请求参数
| 参数名 | 类型   | 必填 | 说明     |
| ------ | ------ | ---- | -------- |
| id     | string | 是   | 会话ID   |

#### 请求体
```json
{
  "title": "会话标题"
}
```

#### 响应结果
```json
{
  "message": "会话更新成功"
}
```

### 6. 流式对话接口

#### 接口说明
与AI进行流式对话，使用Server-Sent Events (SSE) 返回结果

#### 接口地址
```
GET /conversations/stream
```

#### 请求参数
| 参数名 | 类型   | 必填 | 说明               |
| ------ | ------ | ---- | ------------------ |
| id     | string | 是   | 会话ID             |
| input  | string | 是   | 用户输入           |
| model  | string | 否   | 模型名称，格式如openai/gpt-4o |

#### 模型指定方式

支持两种模型指定方式：

1. 使用`服务商/模型`格式：
   - `qwen/qwen-turbo-latest`
   - `deepseek/deepseek-chat`
   - `openai/gpt-4o`
   - `kimi/kimi-k2-0711-preview`

2. 直接指定模型名称：
   - `qwen-turbo-latest` (自动识别为Qwen)
   - `deepseek-chat` (自动识别为DeepSeek)
   - `gpt-4o` (自动识别为OpenAI)
   - `kimi-k2-0711-preview` (自动识别为Kimi)

#### 响应结果
流式对话接口使用 Server-Sent Events (SSE) 格式返回数据，包含三种事件类型：

1. **start** 事件：流式传输开始标记

```
event: start
data: {}
```

2. **data** 事件：实际的消息内容，格式遵循 OpenAI 的流式响应格式

```
event: data
data: {
  "id": "5a69843e-ca50-4127-8228-927ffbddfdf7",
  "object": "chat.completion.chunk",
  "created": 1747828085,
  "model": "deepseek-chat",
  "choices": [
    {
      "index": 0,
      "delta": {
        "content": "SS"
      },
      "logprobs": null,
      "finish_reason": null
    }
  ]
}
```

其中 `delta` 部分在不同阶段可能包含不同的内容：
当生成第一个消息时，可能包含 `role` 字段，如 `{"role": "assistant", "content": ""}`
在生成过程中，主要包含 `content` 字段，如 `{"content": "具体的内容"}`
在结束时，`delta` 为空，但会包含 `finish_reason` 字段，如 `{"finish_reason": "stop"}`

3. **end** 事件：流式传输结束标记

```
event: end
data: {}
```

整个流程是：
发送 `start` 事件表示开始
发送多个 `data` 事件，每个事件包含增量内容
最后发送一个带有 `finish_reason` 的 `data` 事件
发送 `end` 事件表示结束
这种格式与 OpenAI 和 DeepSeek 的流式响应格式兼容。

支持的服务商和默认模型：
- Qwen: `qwen-turbo-latest`
- DeepSeek: `deepseek-chat`
- OpenAI: `gpt-4o`
- Kimi: `kimi-k2-0711-preview`

### 7. 网页翻译接口

#### 接口说明
专门用于网页内容翻译的接口，支持批量翻译多个文本片段

#### 接口地址
```
POST /translate
```

#### 请求参数
| 参数名    | 类型   | 必填 | 说明                           |
| --------- | ------ | ---- | ------------------------------ |
| target    | string | 是   | 目标语言，如 "en" 表示翻译为英语 |
| segments  | array  | 是   | 要翻译的文本片段列表             |
| extra_args| string | 否   | 翻译的额外要求，如风格等         |

#### segments参数说明
| 参数名 | 类型   | 必填 | 说明                                   |
| ------ | ------ | ---- | -------------------------------------- |
| id     | string | 是   | 片段ID，用于标识片段以便返回到前端相应位置 |
| text   | string | 是   | 要翻译的文本内容                        |

#### 请求体示例
```json
{
  "target": "en",
  "segments": [
    {
      "id": "segment1",
      "text": "这是要翻译的文本"
    },
    {
      "id": "segment2", 
      "text": "这是另一段要翻译的文本"
    }
  ],
  "extra_args": "翻译风格要求，如正式、口语化等"
}
```

#### 响应结果
```json
{
  "target": "en",
  "segments": [
    {
      "id": "segment1",
      "text": "This is the text to be translated"
    },
    {
      "id": "segment2",
      "text": "This is another text to be translated"
    }
  ]
}
```