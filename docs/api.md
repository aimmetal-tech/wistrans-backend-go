# API 接口文档

## 基础信息

- 服务器地址: `http://localhost:8080`
- 数据格式: JSON
- 字符编码: UTF-8

## 环境配置

### 必需的环境变量

在项目根目录创建 `.env` 文件，并配置以下环境变量：

```bash
# 阿里云百炼API密钥 (必填，用于联网搜索功能)
QWEN_API_KEY=your_qwen_api_key_here

# 其他可选的大模型API密钥
DEEPSEEK_API_KEY=your_deepseek_api_key_here
OPENAI_API_KEY=your_openai_api_key_here
KIMI_API_KEY=your_kimi_api_key_here
```

### 获取API密钥

1. 访问 [阿里云百炼控制台](https://dashscope.console.aliyun.com/)
2. 登录您的阿里云账号
3. 在控制台中创建API密钥
4. 复制生成的API密钥

### 系统环境变量配置（可选）

您也可以直接在系统环境变量中设置：

```bash
# Windows PowerShell
$env:QWEN_API_KEY="your_qwen_api_key_here"

# Windows CMD
set QWEN_API_KEY=your_qwen_api_key_here

# Linux/macOS
export QWEN_API_KEY="your_qwen_api_key_here"
```

### 注意事项

1. **API密钥安全**: 请妥善保管您的API密钥，不要将其提交到版本控制系统
2. **使用限制**: 请查看阿里云百炼的API使用限制和计费规则
3. **网络要求**: 确保服务器能够访问阿里云百炼的API端点
4. **错误处理**: 如果遇到认证错误，请检查API密钥是否正确配置

### 故障排除

#### 常见错误
- **401 Unauthorized**: API密钥无效或未正确配置
- **403 Forbidden**: API密钥权限不足
- **429 Too Many Requests**: 请求频率超限
- **500 Internal Server Error**: 服务器内部错误

#### 调试步骤
1. 检查环境变量是否正确设置
2. 验证API密钥是否有效
3. 检查网络连接是否正常
4. 查看服务器日志获取详细错误信息

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
  "conversation_id": "会话ID",
  "title": "会话标题",
  "user_id": "用户的ID",
  "created_at": "创建时间",
  "updated_at": "更新时间"
}
```

#### 响应示例
```json
{
  "conversation_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "title": "人工智能发展讨论",
  "user_id": "1234567890",
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
      "message_id": 1,
      "conversation_id": "会话ID",
      "role": "user",
      "content": "用户消息内容",
      "created_at": "消息创建时间"
    },
    {
      "message_id": 2,
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
      "message_id": 1,
      "conversation_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
      "role": "user",
      "content": "你好，请介绍一下人工智能的发展历程",
      "created_at": "2024-01-01T10:00:00Z"
    },
    {
      "message_id": 2,
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

### 8. MCP服务接口

#### 接口说明
提供MCP（Model Context Protocol）服务器配置，支持多种MCP服务的集成，包括网页内容抓取、联网搜索等功能。该接口不仅返回MCP配置，还支持直接执行工具调用。

#### 接口地址
```
POST /mcp
```

#### 请求参数
| 参数名      | 类型     | 必填 | 说明                                    |
| ---------- | -------- | ---- | -------------------------------------- |
| model      | string   | 是   | 模型名称                                |
| mcpServers | string[] | 是   | 需要的MCP服务器名称列表                 |
| query      | string   | 否   | 查询内容（用于工具调用）                |
| tool       | string   | 否   | 要调用的工具名称（web-search/fetch）    |
| params     | object   | 否   | 工具调用参数                            |

#### 请求示例

**1. 获取MCP配置**
```json
{
  "model": "qwen-plus",
  "mcpServers": ["Fetch", "阿里云百炼_联网搜索"]
}
```

**2. 执行联网搜索**
```json
{
  "model": "qwen-plus",
  "mcpServers": ["阿里云百炼_联网搜索"],
  "query": "人工智能最新发展",
  "tool": "web-search",
  "params": {
    "max_results": 5,
    "language": "zh",
    "region": "CN",
    "time_range": "1m"
  }
}
```

**3. 执行网页抓取**
```json
{
  "model": "qwen-plus",
  "mcpServers": ["Fetch"],
  "query": "网页内容抓取",
  "tool": "fetch",
  "params": {
    "url": "https://example.com",
    "content_type": "news",
    "max_length": 3000,
    "language": "auto"
  }
}
```

#### 支持的MCP服务器
目前支持的MCP服务器包括：
- **Fetch**: 网页内容抓取服务，使用Docker运行 `mcp/fetch`
- **阿里云百炼_联网搜索**: 联网搜索服务，使用阿里云百炼MCP服务

#### 响应示例

**1. 仅获取配置**
```json
{
  "model": "qwen-plus",
  "mcpServers": {
    "Fetch": {
      "name": "Fetch",
      "command": "docker",
      "args": ["run", "-i", "--rm", "mcp/fetch"],
      "disabled": false,
      "autoApprove": []
    },
    "阿里云百炼_联网搜索": {
      "name": "阿里云百炼_联网搜索",
      "command": "npx",
      "args": ["mcp-remote", "https://dashscope.aliyuncs.com/api/v1/mcps/WebSearch/sse", "--header", "Authorization:${AUTH_HEADER}"],
      "env": {
        "AUTH_HEADER": "Bearer ${QWEN_API_KEY}"
      },
      "disabled": false,
      "autoApprove": []
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**2. 执行工具调用**
```json
{
  "model": "qwen-plus",
  "mcpServers": {
    "阿里云百炼_联网搜索": {
      "name": "阿里云百炼_联网搜索",
      "command": "npx",
      "args": ["mcp-remote", "https://dashscope.aliyuncs.com/api/v1/mcps/WebSearch/sse", "--header", "Authorization:${AUTH_HEADER}"],
      "env": {
        "AUTH_HEADER": "Bearer ${QWEN_API_KEY}"
      },
      "disabled": false,
      "autoApprove": []
    }
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "query": "人工智能最新发展",
  "tool": "web-search",
  "tool_result": {
    "status": "success",
    "total_count": 5,
    "search_time": "2024-01-15T10:30:00Z",
    "results": [
      {
        "title": "人工智能最新发展动态",
        "url": "https://example.com/ai-news",
        "snippet": "最新的人工智能技术发展...",
        "source": "科技新闻网",
        "published_date": "2024-01-15"
      }
    ]
  }
}
```

#### 错误响应
当请求不支持的MCP服务器时：
```json
{
  "error": "不支持的MCP服务器: InvalidMCP，支持的服务器: [Fetch 阿里云百炼_联网搜索]"
}
```

当工具调用失败时：
```json
{
  "error": "工具调用失败: 执行联网搜索失败: API密钥无效"
}
```

### 9. 网页内容抓取接口

#### 接口说明
使用LLM结合Fetch MCP服务抓取和分析网页内容，支持智能提取结构化信息，特别适用于新闻、文章等内容的抓取

#### 接口地址
```
POST /fetch
```

#### 请求参数
| 参数名        | 类型     | 必填 | 说明                                    |
| ------------ | -------- | ---- | -------------------------------------- |
| url          | string   | 是   | 要抓取的网页URL                         |
| content_type | string   | 否   | 内容类型，支持 news/article/blog，默认news |
| extract_fields| array   | 否   | 要提取的字段列表，不填则使用默认字段       |
| language     | string   | 否   | 内容语言，支持 zh/en/auto，默认auto      |
| max_length   | int      | 否   | 最大内容长度，默认5000                   |

#### 内容类型与默认提取字段
- **news**: title, content, summary, author, publish_date, category
- **article**: title, content, summary, author, publish_date
- **blog**: title, content, author, publish_date, tags

#### 请求体示例
```json
{
  "url": "https://english.news.cn/",
  "content_type": "news",
  "language": "en",
  "max_length": 5000,
  "extract_fields": ["title", "content", "summary", "author", "category"]
}
```

#### 响应结果
```json
{
  "url": "https://english.news.cn/",
  "title": "Xinhua – China, World, Business, Sports, Photos and Video",
  "content": "主要内容...",
  "summary": "新华网English.news.cn提供中英文新闻内容，涵盖中国、国际、商业、体育等主题...",
  "language": "en",
  "status": "success",
  "fetch_time": "2025-08-19T04:01:45Z",
  "extracted_data": {
    "author": "",
    "category": "News", 
    "publish_date": "",
    "tags": ["China", "World", "Business", "Sports", "Culture"]
  }
}
```

#### 错误响应
```json
{
  "url": "https://example.com",
  "status": "error",
  "error": "抓取网页失败: 网络连接超时",
  "fetch_time": "2025-08-19T04:01:45Z"
}
```

#### 使用场景
- 📰 新闻聚合和分析
- 📄 内容自动摘要
- 🔍 网页信息提取
- 📊 媒体监控
- 🌐 多语言内容处理

#### 技术特点
- 🤖 **LLM驱动**: 使用AI模型智能分析内容
- 🐳 **MCP集成**: 集成Fetch MCP服务进行内容抓取
- 📊 **结构化提取**: 自动提取标题、摘要、作者等信息
- 🌍 **多语言支持**: 支持中英文等多语言内容
- 🔧 **灵活配置**: 可自定义提取字段和内容类型

### 10. 联网搜索接口

#### 接口说明
集成阿里云百炼联网搜索MCP服务，提供实时网络搜索功能，支持多语言、多地区搜索，适用于信息查询、新闻搜索、知识检索等场景

#### 接口地址
```
POST /web-search
```

#### 请求参数
| 参数名        | 类型     | 必填 | 说明                                    |
| ------------ | -------- | ---- | -------------------------------------- |
| query        | string   | 是   | 搜索查询关键词                          |
| max_results  | int      | 否   | 最大结果数量，默认10                     |
| language     | string   | 否   | 搜索语言，如 zh/en，默认zh               |
| region       | string   | 否   | 搜索地区，如 CN/US，默认CN               |
| time_range   | string   | 否   | 时间范围，如 1d/1w/1m/1y，默认1y         |
| extra_params | object   | 否   | 额外搜索参数                            |

#### 请求体示例
```json
{
  "query": "人工智能最新发展",
  "max_results": 5,
  "language": "zh",
  "region": "CN",
  "time_range": "1m",
  "extra_params": {
    "safe_search": "moderate"
  }
}
```

#### 响应结果
```json
{
  "query": "人工智能最新发展",
  "results": [
    {
      "title": "2024年人工智能发展最新趋势",
      "url": "https://example.com/ai-trends-2024",
      "snippet": "本文介绍了2024年人工智能领域的最新发展趋势，包括大语言模型、多模态AI等技术的突破...",
      "source": "科技日报",
      "published_at": "2024-01-15",
      "language": "zh"
    }
  ],
  "total_count": 1250,
  "search_time": "2025-08-19T04:15:30Z",
  "status": "success"
}
```

#### 错误响应
```json
{
  "error": "执行联网搜索失败: API请求失败，状态码: 401"
}
```

#### 使用场景
- 🔍 **信息查询**: 实时搜索最新信息
- 📰 **新闻搜索**: 获取最新新闻动态
- 🎓 **知识检索**: 学术研究和知识获取
- 📊 **市场调研**: 行业趋势和竞争分析
- 🌍 **多语言搜索**: 跨语言信息获取

#### 技术特点
- 🌐 **实时搜索**: 基于阿里云百炼联网搜索服务
- 🔐 **安全认证**: 使用Bearer Token进行API鉴权
- 🌍 **多语言支持**: 支持中英文等多种语言搜索
- 📍 **地区定制**: 支持不同地区的搜索结果
- ⏰ **时间过滤**: 支持按时间范围过滤搜索结果
- 🔧 **灵活配置**: 支持自定义搜索参数和额外选项

## MCP使用示例

### 传统MCP客户端配置

在Cursor等MCP客户端中配置：

```json
{
  "mcpServers": {
    "阿里云百炼_联网搜索": {
      "command": "npx",
      "args": [
        "mcp-remote",
        "https://dashscope.aliyuncs.com/api/v1/mcps/WebSearch/sse",
        "--header",
        "Authorization:${AUTH_HEADER}"
      ],
      "env": {
        "AUTH_HEADER": "Bearer ${QWEN_API_KEY}"
      }
    }
  }
}
```

### 通过扩展的MCP接口调用（推荐）

任何大模型都可以通过HTTP请求调用：

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4",
    "mcpServers": ["阿里云百炼_联网搜索"],
    "query": "人工智能最新发展",
    "tool": "web-search",
    "params": {
      "max_results": 5,
      "language": "zh",
      "region": "CN",
      "time_range": "1m"
    }
  }'
```

### Python客户端示例

```python
import requests
import json

def call_mcp_tool(model, tool, query, params=None):
    url = "http://localhost:8080/mcp"
    payload = {
        "model": model,
        "mcpServers": ["阿里云百炼_联网搜索"] if tool == "web-search" else ["Fetch"],
        "query": query,
        "tool": tool,
        "params": params or {}
    }
    
    response = requests.post(url, json=payload)
    return response.json()

# 使用示例
result = call_mcp_tool(
    model="gpt-4",
    tool="web-search",
    query="人工智能最新发展",
    params={"max_results": 3, "language": "zh"}
)
print(json.dumps(result, indent=2, ensure_ascii=False))
```

### JavaScript客户端示例

```javascript
async function callMcpTool(model, tool, query, params = {}) {
    const url = 'http://localhost:8080/mcp';
    const payload = {
        model: model,
        mcpServers: tool === 'web-search' ? ['阿里云百炼_联网搜索'] : ['Fetch'],
        query: query,
        tool: tool,
        params: params
    };
    
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
    });
    
    return await response.json();
}

// 使用示例
const result = await callMcpTool(
    'gpt-4',
    'web-search',
    '人工智能最新发展',
    { max_results: 3, language: 'zh' }
);
console.log(JSON.stringify(result, null, 2));
```

## 使用建议

1. **模型兼容性**: 任何支持HTTP请求的大模型都可以使用此接口
2. **错误处理**: 建议在客户端实现适当的错误处理机制
3. **参数验证**: 在发送请求前验证必需参数
4. **结果缓存**: 对于频繁的搜索请求，建议实现结果缓存
5. **API限制**: 注意各API的调用频率限制