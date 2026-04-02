# MyOpenClaw 学习进度

## 📚 项目概述

**目标**：从零开始，手写一个类似 OpenClaw 的 AI Agent 系统

**实现语言**：Go

**学习方式**：模块化、渐进式，从简单到复杂

**指导原则**：以 OpenClaw 源码为参考，理解其设计思想

---

## 🎯 整体架构

```
┌─────────────────────────────────────────────┐
│                   CLI (main.go)              │
│  - 用户输入                                   │
│  - 消息构造                                   │
│  - 响应显示                                   │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│              Gateway (gateway.go)            │
│  - 会话管理（Session Management）            │
│  - 消息路由（Message Routing）               │
│  - 并发控制（Concurrency Control）           │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│           Agent Runtime (agent.go)           │
│  - 构建 System Prompt                        │
│  - 调用 LLM                                  │
│  - 工具调用（Function Calling）              │
└──────────────────┬──────────────────────────┘
                   │
        ┌──────────┴──────────┐
        ▼                     ▼
┌──────────────┐      ┌──────────────┐
│  LLM Client  │      │    Tools     │
│ (llm/client) │      │ (tools/*)    │
│  - DeepSeek  │      │  - Echo      │
│  - OpenAI    │      │  - File      │
│  - Claude    │      │  - Shell     │
└──────────────┘      └──────────────┘
```

---

## ✅ 已完成的模块

### Stage 1: 基础架构（Mock 版本）

#### 1. 核心数据结构 (`types/message.go`)
**完成时间**：第 1 天

**内容**：
- `Message`：入站消息格式
  - `ID`：消息唯一标识
  - `Content`：消息内容
  - `SessionID`：会话 ID（由 Gateway 分配）
  - `UserID`：平台用户 ID（如 `Cli-User-Id-xxx`）
  - `Channel`：消息来源渠道（如 `iTerm`）
  - `CreatedTime`：创建时间

- `Response`：Agent 响应格式
  - `ID`：响应唯一标识
  - `Content`：响应内容
  - `SessionID`：对应的会话 ID
  - `CreatedTime`：创建时间

- `Session`：会话状态
  - `ID`：会话唯一标识（UUID）
  - `UserID`：用户 ID
  - `Channel`：渠道
  - `LastActivityTime`：最后活跃时间
  - `CreatedTime`：创建时间

**关键设计决策**：
- `UserID` 和 `Channel` 同时存在于 `Message` 和 `Session`
- `UserID` 是平台特定的，不是全局唯一的
- `Session.ID` 使用 UUID，确保唯一性
- Session Key（`UserID + "_" + Channel`）用于查找当前活跃会话

**学习要点**：
- Message 的自包含性（Self-contained）
- Session Key vs Session ID 的区别
- 为什么需要 `Channel` 字段

---

#### 2. 工具系统 (`tools/tool.go`, `tools/echo.go`)
**完成时间**：第 1 天

**内容**：
- `Tool` 接口：
  - `Name() string`：工具名称
  - `Description() string`：工具描述（含参数说明）
  - `Execute(args map[string]interface{}) (string, error)`：执行工具

- `EchoTool` 实现：
  - 简单的回显工具
  - 演示工具接口的实现
  - 安全的类型断言（`value, ok := args["key"].(type)`）

**关键设计决策**：
- 接口化设计，便于扩展
- `Description` 包含详细的参数说明和示例
- 错误处理：参数缺失、类型错误

**学习要点**：
- Go 的接口设计
- 安全的类型断言
- 错误处理的最佳实践

---

#### 3. Gateway (`gateway/gateway.go`)
**完成时间**：第 1 天

**内容**：
- `Gateway` 结构：
  - `Session map[string]*types.Session`：会话存储（Session Key → Session 指针）
  - `RunTime *agent.Runtime`：Agent Runtime 引用
  - `GlobalRw sync.RWMutex`：全局读写锁

- `NewGateway()`：构造函数，初始化 Session Map

- `getOrCreateSession(msg *types.Message) (*types.Session, error)`：
  - 根据 `UserID + "_" + Channel` 生成 Session Key
  - 如果存在则返回，不存在则创建新 Session
  - 使用 `sync.RWMutex` 保证线程安全

- `HandleMessage(ctx context.Context, msg *types.Message) (*types.Response, error)`：
  - 获取或创建 Session
  - 将 Session.ID 赋值给 msg.SessionID
  - 调用 Runtime.ProcessMessage()
  - 返回响应

**关键设计决策**：
- Session Map 存储指针（`*types.Session`），而不是值
- Session Key 用于查找，Session ID 用于标识
- 全局锁保证并发安全（Stage 1 简化版）
- Gateway 不处理业务逻辑，只做路由

**学习要点**：
- Go 的 Map 存储指针 vs 值的区别
- 并发安全：`sync.RWMutex` 的使用
- Session Key 的设计（确定性 vs 唯一性）
- `/new` 命令的检测逻辑（Gateway 层处理）

---

#### 4. Agent Runtime (`agent/agent.go`)
**完成时间**：第 1 天（Mock 版本）

**内容**：
- `Runtime` 结构：
  - `Tools map[string]tools.Tool`：工具 Map

- `NewRuntime()`：初始化 Runtime，注册 EchoTool

- `buildSystemPrompt() string`：
  - 构建 System Prompt（Markdown 格式）
  - 包含 Agent 身份、工具列表、工作原则
  - 使用 `strings.Builder` 高效拼接

- `ProcessMessage(ctx context.Context, msg *types.Message) (*types.Response, error)`：
  - 构建 System Prompt
  - 返回 Mock 响应（Stage 1）

**关键设计决策**：
- System Prompt 使用 Markdown 格式
- 工具列表动态生成
- Mock 阶段返回 System Prompt 用于验证

**学习要点**：
- `strings.Builder` 的使用
- System Prompt 的结构设计
- LLM Client 的 Lazy Loading 概念

---

#### 5. CLI 入口 (`main.go`)
**完成时间**：第 1 天

**内容**：
- 加载 `.env` 文件（`godotenv.Load()`）
- 检查 API Key 是否设置
- 初始化 Gateway 和 Runtime
- 连接 Gateway 和 Runtime（`gw.RunTime = runtime`）
- 循环读取用户输入
- 构造 Message 并发送给 Gateway
- 显示 Agent 响应

**关键设计决策**：
- SessionID 由 Gateway 分配，main.go 里留空
- UserID 使用 `Cli-User-Id` + UUID
- Channel 固定为 `iTerm`
- 支持 `exit` 命令退出

**学习要点**：
- 环境变量的安全管理
- `bufio.Scanner` 读取用户输入
- 模块间的连接方式

---

### Stage 2: 集成真实 LLM

#### 6. LLM Client (`llm/client.go`)
**完成时间**：第 2 天

**内容**：
- `LLMClient` 接口：
  - `Chat(ctx context.Context, systemPrompt string, userMessage string) (string, error)`

- `DeepSeekClient` 结构：
  - `APIKey`：从环境变量读取
  - `BaseURL`：`https://api.deepseek.com`
  - `Model`：`deepseek-chat`

- `DeepSeekRequest`：请求格式
  - `Messages`：消息数组（system + user）
  - `Stream`：是否流式返回
  - `Model`：模型名称

- `DeepSeekResponse`：响应格式
  - `Choices[0].Message.Content`：LLM 的回复

- `Chat()` 方法：
  - 构造请求体（JSON）
  - 发送 HTTP POST 请求
  - 设置 Authorization Header
  - 解析响应
  - 返回 LLM 回复

**关键设计决策**：
- 接口化设计，支持多个 LLM 提供商
- API Key 从环境变量读取（安全）
- 完整的错误处理（网络、序列化、API 错误）
- 检查响应是否为空

**学习要点**：
- Go 的 HTTP 客户端使用
- JSON 序列化/反序列化
- 环境变量的安全管理（`.env` 文件 + `.gitignore`）
- Shell `${}` vs Go 字符串拼接的区别

**遇到的问题**：
- ❌ JSON 字段名错误：`message` → `messages`（复数）
- ❌ Authorization Header 格式：`Bearer ${}` → `Bearer `
- ❌ `ioutil.ReadAll` 已废弃 → `io.ReadAll`
- ❌ 没有检查 `Choices` 是否为空
- ❌ API Key 硬编码 → 环境变量
- ✅ 全部修复

---

#### 7. 改造 Agent Runtime
**完成时间**：第 2 天

**改动**：
- 添加 `LLM llm.LLMClient` 字段
- `NewRuntime()` 里初始化 DeepSeek Client
- `ProcessMessage()` 调用真实的 LLM
- 返回 LLM 的真实响应

**测试结果**：
```
用户: 你好
Agent: 你好！我是 MyOpenClaw，很高兴为你服务。有什么我可以帮助你的吗？
```
✅ 成功调用 DeepSeek API！

---

## 🚧 正在进行的模块

### Stage 3: Function Calling（工具调用）

**目标**：让 LLM 能够调用工具（如 Echo Tool）

**已完成**：
- ✅ 扩展 Tool 接口，添加 `Parameters()` 方法
- ✅ 实现 EchoTool 的 `Parameters()` 方法（返回 JSON Schema）
- ✅ 修改 `LLMClient.Chat()` 签名，支持传入工具列表
- ✅ 验证 OpenClaw 源码，确认使用 JSON Schema 定义工具参数

**正在进行**：
- 🚧 改造 `llm/client.go` 数据结构：
  - 修改 `DeepSeekToolFunction`（`Properties` → `Parameters`）
  - 添加 `ToolCall` 结构体
  - 修改 `DeepSeekResponse`（添加 `ToolCalls` 字段）
  - 修改工具注入逻辑（使用 `tool.Parameters()`）

**待完成**：
- [ ] 改造 `Chat()` 方法，解析 tool_calls
- [ ] 改造 `agent/agent.go`，实现工具调用循环
- [ ] 测试完整的 Function Calling 流程

**核心流程**：
```
1. 用户消息 → LLM（带工具列表）
2. LLM 返回 tool_calls
3. 执行工具
4. 工具结果 → LLM
5. LLM 返回最终回复
6. 返回给用户
```

**状态**：正在改造数据结构

---

## 📋 待完成的功能

### Stage 4: 命令系统
- [ ] `/new` - 创建新会话
- [ ] `/reset` - 重置当前会话
- [ ] `/help` - 显示帮助信息

### Stage 5: 多轮对话
- [ ] 会话历史管理
- [ ] Context Window 管理
- [ ] 自动裁剪历史消息

### Stage 6: 更多工具
- [ ] File Tool（读写文件）
- [ ] Shell Tool（执行命令）
- [ ] Search Tool（搜索）

### Stage 7: 多入口支持
- [ ] HTTP API（REST）
- [ ] WebSocket
- [ ] Telegram Bot
- [ ] Discord Bot

### Stage 8: 高级特性
- [ ] 流式响应（Streaming）
- [ ] 多 Agent 协作
- [ ] 插件系统
- [ ] 持久化存储

---

## 📂 项目结构

```
myopenclaw/
├── .env                    # 环境变量（不提交到 Git）
├── .gitignore              # Git 忽略文件
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖校验
├── main.go                 # CLI 入口
├── PROGRESS.md             # 学习进度（本文件）
│
├── types/                  # 核心数据结构
│   └── message.go          # Message, Response, Session
│
├── gateway/                # Gateway（会话管理 + 路由）
│   └── gateway.go          # Gateway 实现
│
├── agent/                  # Agent Runtime（执行引擎）
│   └── agent.go            # Runtime 实现
│
├── llm/                    # LLM Client（大模型调用）
│   └── client.go           # LLMClient 接口 + DeepSeek 实现
│
└── tools/                  # 工具系统
    ├── tool.go             # Tool 接口
    └── echo.go             # Echo 工具实现
```

---

## 🎓 学到的核心概念

### 1. 架构设计
- **多入口、单内核**：多个入口（CLI、HTTP、WebSocket）共享同一个 Gateway 和 Runtime
- **分层设计**：CLI → Gateway → Runtime → LLM/Tools
- **接口化**：Tool 接口、LLM 接口，便于扩展

### 2. 会话管理
- **Session Key vs Session ID**：
  - Session Key（`UserID:Channel`）：确定性，用于查找当前会话
  - Session ID（UUID）：唯一性，用于标识会话实例
  - 一对一（同一时刻）vs 一对多（时间线上）
- **Context Window**：由用户通过 `/new`、`/reset` 命令控制
- **并发控制**：使用 `sync.RWMutex` 保证线程安全

### 3. LLM 集成
- **Lazy Loading**：LLM Client 按需创建，不在启动时加载
- **System Prompt**：包含 Agent 身份、工具列表、工作原则
- **API 调用**：HTTP POST + JSON 序列化/反序列化
- **环境变量管理**：`.env` 文件 + `godotenv` 库

### 4. Go 语言特性
- **接口（Interface）**：鸭子类型，只要实现了方法就满足接口
- **指针 vs 值**：Map 存储指针可以修改原对象
- **并发安全**：`sync.RWMutex` 的 Lock/Unlock
- **错误处理**：`fmt.Errorf` 包装错误，提供上下文
- **类型断言**：`value, ok := args["key"].(type)` 安全断言
- **环境变量**：`os.Getenv()` 读取，`godotenv.Load()` 加载 `.env` 文件
- **HTTP 客户端**：`http.NewRequest()` + `http.Client.Do()`
- **JSON 序列化**：`json.Marshal()` 和 `json.Unmarshal()`
- **包查找顺序**：标准库 → go.mod 依赖 → 当前项目（需要明确路径）

### 5. 安全原则
- **不吞异常**：所有错误必须处理或向上传递
- **API Key 安全**：使用环境变量，不硬编码
- **类型安全**：使用安全的类型断言

---

## 🐛 遇到的问题和解决方案

### 问题 1：Session ID 的生成逻辑
**问题**：最初使用 `UserID + Channel` 生成 Session.ID，导致固定 ID，无法创建新会话

**解决**：
- Session Key（`UserID + "_" + Channel`）：用于 Map 查找
- Session ID（UUID）：用于唯一标识会话实例
- 通过 `/new` 命令可以创建新的 Session ID

---

### 问题 2：Map 存储值 vs 指针
**问题**：`map[string]types.Session` 存储值，修改 Session 不会更新 Map

**解决**：改为 `map[string]*types.Session` 存储指针

---

### 问题 3：并发安全
**问题**：多个 goroutine 同时访问 Session Map 会导致竞态条件

**解决**：使用 `sync.RWMutex` 在 `getOrCreateSession` 里加锁

---

### 问题 4：LLM Client 的加载时机
**问题**：是否需要在启动时加载所有 LLM Client？

**解决**：Lazy Loading，按需创建（OpenClaw 的设计）

---

### 问题 5：DeepSeek API 调用失败（401）
**问题**：API Key 没有正确传递

**解决**：
- 使用 `godotenv.Load()` 加载 `.env` 文件
- 在 `main.go` 里检查 API Key 是否设置

---

### 问题 6：JSON 字段名错误
**问题**：`Message` vs `Messages`（单数 vs 复数）

**解决**：参考官方文档，使用 `messages`（复数）

---

### 问题 7：GoLand Debug 配置
**问题**：`package myopenclaw is not in std`

**原因**：
- GoLand 的工作目录可能是 `joke`，而不是 `myopenclaw`
- 配置里写 `myopenclaw` 会让 Go 去标准库找
- Go 的包查找顺序：先标准库 → 再 go.mod 依赖 → 最后当前项目
- 如果写 `myopenclaw`（不带路径），Go 会认为这是标准库的包

**解决**：
- 使用 `.` 表示当前目录
- 或使用相对路径 `./myopenclaw`
- 或修改工作目录为 `myopenclaw` 目录
- 或直接在终端运行

**关键理解**：
- ❌ `软件包路径: myopenclaw` → Go 去标准库找
- ✅ `软件包路径: .` → Go 编译当前目录
- ✅ `软件包路径: ./myopenclaw` → Go 编译相对路径

---

### 问题 8：环境变量加载
**问题**：`os.Getenv("DEEPSEEK_API_KEY")` 返回空

**原因**：
- Go 程序不会自动读取 `.env` 文件
- 需要使用 `godotenv` 库来加载

**解决**：
1. 安装依赖：`go get github.com/joho/godotenv`
2. 在 `main.go` 里加载：
   ```go
   import "github.com/joho/godotenv"
   
   func main() {
       err := godotenv.Load()
       if err != nil {
           fmt.Println("警告: 未找到 .env 文件")
       }
       // ...
   }
   ```

**验证**：
- 创建 `.env` 文件：`DEEPSEEK_API_KEY=sk-xxx`
- 添加到 `.gitignore`：`.env`
- 运行程序，检查 API Key 是否加载成功

---

## 📊 当前状态

```
✅ types/message.go      - 核心数据结构
✅ tools/tool.go         - 工具接口（已添加 Parameters() 方法）
✅ tools/echo.go         - Echo 工具（已实现 Parameters()）
✅ llm/client.go         - LLM 接口 + DeepSeek 实现（支持工具参数）
✅ agent/agent.go        - Agent Runtime（真实 LLM 调用）
✅ gateway/gateway.go    - Gateway（会话管理 + 路由）
✅ main.go               - CLI 入口（支持 .env 加载）
✅ .env                  - 环境变量配置
✅ .gitignore            - Git 忽略文件
✅ PROGRESS.md           - 学习进度文档

🚧 正在进行：改造 llm/client.go 数据结构（Function Calling）
```

---

## 🎯 下一步计划

### Stage 3: Function Calling

**目标**：让 LLM 能够调用 Echo Tool

**步骤**：
1. 改造 `llm/client.go`：
   - 添加 `Tools` 参数到 `Chat()` 方法
   - 添加 `ToolCalls` 解析

2. 改造 `agent/agent.go`：
   - 实现工具调用循环
   - 解析 tool_calls
   - 执行工具
   - 把结果再发给 LLM

**预期效果**：
```
用户: 帮我回显 Hello World
LLM: [调用 echo 工具]
工具: Hello World
LLM: 已为您回显：Hello World
```

---

## 📝 学习笔记

### OpenClaw 的核心设计思想

1. **单一职责**：每个模块只做一件事
   - Gateway：会话管理 + 路由
   - Runtime：业务逻辑 + LLM 调用
   - Tools：具体功能实现

2. **接口化**：便于扩展和测试
   - Tool 接口
   - LLM 接口
   - Bridge 接口（未实现）

3. **并发安全**：使用锁保护共享资源
   - Gateway 的 Session Map
   - 未来：每个 Session 的独立锁

4. **Lazy Loading**：按需加载，节省资源
   - LLM Client 按需创建
   - 不在启动时加载所有资源

---

## 🎉 里程碑

- ✅ **2026-03-24**：项目启动，完成基础架构（Mock 版本）
- ✅ **2026-03-26 上午**：集成 DeepSeek LLM，实现真实对话
- ✅ **2026-03-26 下午**：扩展 Tool 接口，添加 Parameters() 方法
- 🚧 **进行中**：实现 Function Calling（工具调用）

---

## 💡 重要的设计讨论

### 讨论 1：Channel 应该放在哪里？
**问题**：Channel 是 Message 的属性还是 Session 的属性？

**结论**：两者都需要
- Message 里的 Channel：自包含性，方便日志和追踪
- Session 里的 Channel：Session 的核心属性

**参考**：OpenClaw 的 `MsgContext` 包含 `From` 和 `SenderId`

---

### 讨论 2：UserID 是什么？
**问题**：UserID 是全局唯一的吗？

**结论**：不是
- UserID 是平台特定的（如 `tg-123456`、`discord-789`）
- 不同平台的 UserID 可能重复
- 需要 Binding 系统来关联不同平台的同一个用户

**参考**：OpenClaw 的设计

---

### 讨论 3：Session Key vs Session ID
**问题**：为什么需要两个 ID？

**结论**：
- Session Key（`UserID:Channel`）：确定性，用于查找"当前"会话
- Session ID（UUID）：唯一性，用于标识会话实例
- 一个 Session Key 可以对应多个 Session ID（时间线上）
- 通过 `/new` 命令创建新 Session ID

**类比**：
- Session Key = 你的手机号（固定，用于找到你）
- Session ID = 每次通话记录的 ID（唯一，标识每次通话）

---

### 讨论 4：命令检测
**问题**：如何检测 `/new` 命令？如果用户说"我找了个/new女友"怎么办？

**结论**：
- OpenClaw 的命令必须以 `/` 开头
- 在 Gateway 层检测（`strings.HasPrefix(msg.Content, "/")`）
- "我找了个/new女友" 不会被识别为命令

---

### 讨论 5：LLM Client 的加载时机
**问题**：每个 Session 可以配置不同的 LLM，是启动时全部加载还是按需加载？

**结论**：
- Lazy Loading（按需加载）
- 启动时 LLM Client 是空的
- 第一次使用时才创建
- 节省资源，提高启动速度

**参考**：OpenClaw 的设计

---

### 讨论 6：工具参数的定义方式
**问题**：工具是全局通用的，为什么要放在 DeepSeekRequest 里？

**结论**：
- 工具定义确实是全局通用的（在 `tools/` 目录）
- 但每次调用 LLM 时，需要把工具列表"翻译"成 JSON Schema 发给 LLM
- 这是 LLM Function Calling 的标准流程
- LLM 需要知道有哪些工具可用，以及每个工具的参数结构

**验证**：
- 查看 OpenClaw 源码 `src/agents/tools/web-search.ts`
- Tool 结构包含 `parameters` 字段（JSON Schema 格式）
- OpenClaw 使用 TypeBox 库生成 JSON Schema
- 我们使用 Go 的 `map[string]interface{}` 表示相同结构

**参考文件**：
- `openclaw/src/agents/tools/browser-tool.ts` 第 390 行
- `openclaw/src/agents/tools/browser-tool.schema.ts` 第 88 行

---

### 讨论 7：JSON Schema 的含义
**问题**：`type: "object"`, `properties`, `required` 这些都是什么意思？

**答案**：这是 JSON Schema 格式，用来描述 JSON 数据的结构

**结构解释**：
```json
{
  "type": "object",           // 参数是一个对象
  "properties": {             // 对象包含哪些字段
    "message": {              // 字段名
      "type": "string",       // 字段类型
      "description": "说明"   // 字段说明
    }
  },
  "required": ["message"]     // 必填字段列表
}
```

**为什么需要 JSON Schema？**
- LLM 需要精确理解工具的参数结构
- 纯文本描述不够精确，LLM 可能生成错误的参数
- JSON Schema 是 LLM Function Calling 的标准格式

**示例**：
- ❌ 只用文本："需要一个 message 参数，类型是字符串"
  - LLM 可能生成：`{"msg": "Hello"}`（字段名错了）
- ✅ 使用 JSON Schema：LLM 会严格按照结构生成
  - LLM 生成：`{"message": "Hello"}`（正确）

---

## 📚 参考资料

- **OpenClaw 学习笔记**：`/Users/qiyuekurong/go/cursor/joke/OpenClaw学习笔记.md`
- **DeepSeek API 文档**：https://api.deepseek.com
- **Go 官方文档**：https://go.dev/doc/

---

## 🔧 开发环境

- **Go 版本**：1.22.12
- **IDE**：Cursor + GoLand
- **操作系统**：macOS
- **依赖**：
  - `github.com/google/uuid` v1.6.0
  - `github.com/joho/godotenv` v1.5.1

---

## 🎯 当前任务

**✅ Function Calling 已完成！**

---

### 2026-03-26 的进度

#### 上午：集成真实 LLM
- ✅ 创建 `llm/client.go`（LLMClient 接口 + DeepSeek 实现）
- ✅ 实现 HTTP 调用 DeepSeek API
- ✅ 改造 `agent/agent.go`，调用真实 LLM
- ✅ 解决环境变量加载问题（godotenv）
- ✅ 解决 GoLand Debug 配置问题
- ✅ 成功实现真实对话！

**测试结果**：
```
用户: 你好
Agent: 你好！我是 MyOpenClaw，很高兴为你服务。有什么我可以帮助你的吗？
```

#### 下午：开始实现 Function Calling
- ✅ 扩展 Tool 接口，添加 `Parameters()` 方法
- ✅ 实现 EchoTool 的 `Parameters()` 方法
- ✅ 验证 OpenClaw 源码，确认设计方向
- 🚧 开始改造 `llm/client.go` 数据结构

---

### 2026-03-27 的进度

#### 架构重构：抽象层设计
**问题发现**：`LLMClient` 接口使用 `[]DeepSeekMessage` 违反了依赖倒置原则

**解决方案**：参考 OpenClaw 的设计，引入抽象层

**实现步骤**：
1. ✅ 创建 `types/llm_message.go`：
   - `LLMMessage`：通用消息格式
   - `ToolCall`：通用工具调用格式（`Arguments` 是 `map[string]interface{}`）

2. ✅ 创建 `constant/constant.go`：
   - 定义角色常量（`SysRole`, `UserRole`）

3. ✅ 拆分 LLM 层：
   - `llm/client.go`：接口定义（使用通用类型）
   - `llm/deepseekClient.go`：DeepSeek 适配器（实现转换逻辑）

4. ✅ 实现转换函数：
   - `convertToDeepSeekMessages()`：`[]types.LLMMessage` → `[]DeepSeekMessage`
   - `convertToCommonLLMTool()`：`[]DeepSeekToolCall` → `[]types.ToolCall`

5. ✅ 修改 `agent/agent.go`：
   - 使用 `[]types.LLMMessage` 构建消息
   - 实现工具调用循环（最多 5 轮）
   - 添加 `executeTool()` 方法

**关键设计决策**：
- **通用层**（`types.ToolCall`）：`Arguments` 是已解析的 `map[string]interface{}`
- **适配器层**（`DeepSeekToolCall`）：`Arguments` 是 JSON 字符串
- **转换逻辑**：适配器负责格式转换（`json.Marshal` / `json.Unmarshal`）

**参考 OpenClaw 源码**：
- `src/agents/openai-ws-stream.ts:564-575`：ToolCall 的通用格式
- `src/agents/ollama-stream.ts:265-311`：`convertToOllamaMessages()` 转换函数

---

### 2026-04-02 的进度

#### Function Calling 完整实现与测试

**最终修复**：
1. ✅ 修复 `DeepSeekRequest` 添加 `Stream` 字段
2. ✅ 修复 `convertToDeepSeekMessages()` 的 `make` 长度问题
3. ✅ 修复 `convertToCommonLLMTool()` 的 `Arguments` 赋值问题
4. ✅ 修复 `types.ToolCall` 字段（删除多余字段，`Parameters` → `Arguments`）
5. ✅ 编译通过！

**测试结果**：
```
用户: 帮我 echo 一下 hello world

[第一轮]
LLM 返回 tool_calls: 
  - name: "echo"
  - arguments: {"message": "hello world"}

[Agent 执行工具]
echo 工具返回: "hello world"

[第二轮]
LLM 收到工具结果，生成最终回复:
"已成功回显 \"hello world\"。"

✅ 成功！
```

**验证的功能**：
- ✅ LLM 正确识别需要调用工具
- ✅ Agent 正确执行工具
- ✅ 工具结果正确返回给 LLM
- ✅ LLM 生成自然语言回复
- ✅ 完整的多轮对话流程

---

## 📊 当前架构总结

```
┌─────────────────────────────────────────┐
│            main.go (CLI)                │
│         用户输入 → 消息构建              │
└──────────────────┬──────────────────────┘
                   │
                   ↓
┌─────────────────────────────────────────┐
│         Gateway (消息路由)               │
│  - Session 管理                          │
│  - 消息分发                              │
└──────────────────┬──────────────────────┘
                   │
                   ↓
┌─────────────────────────────────────────┐
│      Agent Runtime (执行引擎)            │
│  - buildSystemPrompt()                  │
│  - ProcessMessage() [工具调用循环]       │
│  - executeTool()                        │
└──────────────────┬──────────────────────┘
                   │
                   ↓
┌─────────────────────────────────────────┐
│      LLMClient (抽象层)                  │
│  使用 types.LLMMessage                   │
└──────────────────┬──────────────────────┘
                   │
                   ↓
┌─────────────────────────────────────────┐
│    DeepSeekClient (适配器层)             │
│  - convertToDeepSeekMessages()          │
│  - convertToCommonLLMTool()             │
└──────────────────┬──────────────────────┘
                   │
                   ↓
┌─────────────────────────────────────────┐
│         DeepSeek API                    │
│  Function Calling 完整支持               │
└─────────────────────────────────────────┘
```

---

## 🎓 今天学到的核心概念

### 1. 依赖倒置原则（Dependency Inversion Principle）
- **错误设计**：接口依赖具体实现（`LLMClient` 使用 `DeepSeekMessage`）
- **正确设计**：接口依赖抽象（`LLMClient` 使用 `types.LLMMessage`）
- **好处**：添加 OpenAI、Claude 等 Provider 时，不需要修改接口

### 2. 适配器模式（Adapter Pattern）
- **抽象层**：`types.LLMMessage`（业务层使用）
- **适配器层**：`DeepSeekClient`（负责格式转换）
- **转换函数**：`convertToDeepSeekMessages()` / `convertToCommonLLMTool()`

### 3. Function Calling 完整流程
```
用户消息 → LLM 分析 → 返回 tool_calls → Agent 执行工具 
→ 工具结果 → LLM 生成回复 → 返回用户
```

### 4. 多轮对话的消息历史管理
- System 消息（第一条）
- User 消息
- Assistant 消息（带 tool_calls）
- Tool 消息（带 tool_call_id）
- Assistant 消息（最终回复）

---

## 📂 新增/修改的文件

### 新增文件
1. `types/llm_message.go`：通用消息类型定义
2. `constant/constant.go`：角色常量定义

### 重构文件
1. `llm/client.go`：接口定义（抽象层）
2. `llm/deepseekClient.go`：DeepSeek 适配器（实现层）
3. `agent/agent.go`：添加工具调用循环

---

*最后更新：2026-04-02 16:30*
