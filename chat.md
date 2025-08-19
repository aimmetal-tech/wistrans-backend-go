# 智慧译-后端项目 模型输入记录

## 后端项目初始化
- 后端项目初始化——搭建大模型连续对话接口。
```刘睿琪
现在让我们来实现**智慧译-后端项目**的一些功能。
## 实现功能：接入几个主流的大模型来实现连续对话功能。
## 描述：
1. 用户根据项目下的**.env.example**模板文件来创建**.env**来填写他们的API，程序先从.env本地文件中获取所需的API，如果用户没填则从系统变量中获取，也就是XXX_API_KEY，用户最少填写一个API，否则程序提示“未填写API”。

2. 要接入的大模型分别为**Qwen、DeepSeek、Kimi、OpenAI**。如果只有一个API填写，base_url就用那个大模型服务的url；如果都填写了则使用Qwen的url作为默认

3. 暂时不需要连接数据库，使用数据类型dict来存储会话记录。

```

- 大模型连续对话服务优化改进。

```刘睿琪
接下来我们还是来进行大模型对话功能的优化改进。
## 内容
1. Qwen默认模型使用qwen-turbo-latest，DeepSeek默认模型使用deepseek-chat，OpenAI默认使用gpt-5，Kimi默认使用kimi-k2-0711-preview。

2. 新增输入Service选择服务商（大小写均可）：Qwen、DeepSeek、OpenAI、Moonshot，如果用户输入模型名称也可以自动匹配上。

3. 关于API的调用机制：先从环境变量中获取，再从.env中获取。
```

## 后端项目-Gin框架

- 搭建Gin初始框架。刘睿琪

```刘睿琪
现在我们来引入新功能：Gin框架
概述：将连续对话功能放入Gin框架内
具体内容：
1. 设计三个接口：一个`/health`健康测试接口、一个`/conversations`用于创建新会话或替换已有会话标题、一个`/conversations/{id}/stream?input=XXX`
2. `/conversations`这个接口的信息：一、Method：GET，返回值为创建会话的ID；二、Method:：PATCH
3. `/conversations/{id}/stream?input=XXX`是SSE、Method：GET、返回的结果要符合OpenAI的消息聚合格式
4. 同时接入Postgre数据库，将会话数据存入数据库，数据库URL已存在.env文件中。在项目启动时，先自动检测是否连接上数据库、数据库XXX是否存在、如果存在是否存在对应的表：conversations、messages。如果不存在则进行创建操作，如果存在则直接使用。
```

### 大模型流式对话接口

- 优化SSE-流式对话接口

```刘睿琪
`/conversations/{id}/stream`接口修改为`/conversations/stream`，Method依然是GET，Query有id、input、model。model设计为XXX/AAA格式，如openai/gpt5；model也可以不指定服务商直接指定模型。
```

- 把大模型真正接入，替换之前的模拟对话功能

```刘睿琪
现在请完善大模型API调用，实现真正的与大模型流式对话。
内容：把之前的Qwen3、Kimi、DeepSeek、OpenAI接入现有的项目中，实现流式对话功能，统一采用OpenAI消息聚合格式。
```

- 优化返回消息格式

```刘睿琪
{
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

参考这份DeepSeek的消息格式，将端口`/conversations/stream`的返回内容改为上面的格式。
```

### 翻译接口

- 翻译接口初始化

```刘睿琪
现在我们来设计一个翻译接口/translate专门用于翻译网页 接口信息： Method POST Body示例：{target:"en", "segments": [{"id": "xxx-用于标识片段以便后续返回到前端相应位置", "text": "这是要翻译的文本"}], "extra_args": "这是额外参数，例如翻译的风格"}
```

```
对于"翻译风格要求: {extra_args}"这一步进行修改 新的实例Body如下： { "target": "en", "segments": [ { "id": "segment1", "text": "这是要翻译的文本" }, { "id": "segment2", "text": "这是另一段要翻译的文本" } ], "extra_args": { "style": "专业" } }
```

### 大模型对话标题自动生成

```刘睿琪
现在让我们来实现对话标题内容自动生成。 使用更新对话接口`PATCH /conversations/{id}` 内容：当用户创建对话并有了第一条聊天内容后，调用大模型总结并预测对话内容，生成25字以内的标题，然后再调用"更新对话接口`PATCH /conversations/{id}`"来把对话标题数据更新。
```

### 对话详情接口和历史记录

```
再次新增功能
接口1`/conversations/detail`，Query为id，返回值为该对话的详细信息
接口2`/conversations/history`, Query为id，返回值为该对话的历史记录
将这两个接口实现后记得更新docs/api.md，遵循routes-rule.md规则
```

### MCP支持
```
采用"扩展现有模型参数格式"，来添加对MCP的支持。现在开始进行代码的编写。
```

```
{
    "mcpServers": {
        "github.com/github/github-mcp-server": {
            "command": "docker",
            "args": [
                "run",
                "-i",
                "--rm",
                "-e",
                "GITHUB_PERSONAL_ACCESS_TOKEN",
                "ghcr.io/github/github-mcp-server"
            ],
            "env": {
                "GITHUB_PERSONAL_ACCESS_TOKEN": "XXX"
            },
            "disabled": false,
            "autoApprove": []
        },
        "Fetch": {
            "command": "docker",
            "args": [
                "run",
                "-i",
                "--rm",
                "mcp/fetch"
            ]
        }
    }
}

参考这个文件里的Fetch MCP，现在来给我们的MCP服务添加这个默认的MCP服务
```
