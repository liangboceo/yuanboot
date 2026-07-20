# Eino 框架开发 Spec

> 基于当前仓库 `github.com/cloudwego/eino` 扫描，并结合 CloudWeGo Eino 中文官网文档维护生成，用于后续基于 Eino 做 LLM 应用、组件、Agent 与编排流程开发。

## 1. 框架定位

Eino 是 Go 语言 LLM 应用开发框架，提供：

- 统一组件抽象：`ChatModel`、`Tool`、`Retriever`、`Indexer`、`Embedding`、`ChatTemplate`、`Loader`、`Transformer`。
- 编排能力：`compose.Graph`、`compose.Chain`、`compose.Workflow` 组合为可执行 `Runnable`。
- ADK 智能体套件：`ChatModelAgent`、工具调用、Runner、事件流、中断/恢复、检查点、多 Agent、预置 Agent。
- 通用数据模型：`schema.Message`、`schema.AgenticMessage`、`schema.Document`、`schema.ToolInfo`、`schema.StreamReader`。
- 回调切面：统一注入日志、追踪、指标、调试能力。

Go 版本要求：`go 1.18+`。

## 2. 顶层模块职责

| 目录 | 职责 |
| --- | --- |
| `schema` | 核心数据结构与流式抽象：消息、文档、工具定义、序列化、厂商扩展字段。 |
| `components` | 组件接口与 Options：model、tool、prompt、retriever、indexer、embedding、document。 |
| `compose` | 编排核心：Graph、Chain、Workflow、Runnable、节点、边、分支、并行、状态、字段映射、检查点、中断/恢复。 |
| `callbacks` | 统一 Handler、RunInfo、全局/局部 callback 注入与流式 callback 约束。 |
| `adk` | Agent 开发套件：Agent、Runner、ChatModelAgent、AgentTool、TurnLoop、中断/恢复、取消、重试、故障转移。 |
| `adk/middlewares` | Agent 中间件：filesystem、skill、summarization、reduction、dynamic tool search、plan task、agents.md。 |
| `adk/prebuilt` | 预置 Agent：supervisor、planexecute、deep。 |
| `flow` | 高层经典流程封装：ReAct agent、RAG retriever/indexer、多 Agent host。 |
| `utils` | 辅助能力，主要是 callback helper。 |
| `internal` | 内部实现，业务项目与外部扩展不应依赖。 |

## 3. 核心开发模型

### 3.1 数据类型优先级

优先复用 `schema`：

- 对话：`*schema.Message`
- Agentic 对话：`*schema.AgenticMessage`
- 文档/RAG：`*schema.Document`
- 工具描述：`*schema.ToolInfo`
- 工具参数：`schema.ParameterInfo` 或 JSON Schema
- 流式数据：`*schema.StreamReader[T]` / `schema.StreamWriter[T]`

不要为消息、工具、文档重复定义平行结构，除非是业务 DTO。

### 3.2 四种执行范式

`compose.Runnable[I, O]` 支持：

| 模式 | 签名 | 语义 |
| --- | --- | --- |
| Invoke | `I -> O` | 非流式输入、非流式输出。 |
| Stream | `I -> StreamReader[O]` | 非流式输入、流式输出。 |
| Collect | `StreamReader[I] -> O` | 流式输入、非流式输出。 |
| Transform | `StreamReader[I] -> StreamReader[O]` | 流式输入、流式输出。 |

框架会自动桥接范式不匹配：非流式可包装为单 chunk fake stream，流式可 concat 为完整值。

流式约束：

- `StreamReader` 只能读取一次。
- 读取方必须 `Close()`，包括读到 `io.EOF` 后。
- 多消费者必须先 `Copy()`。
- callback 中拿到的 stream copy 也必须关闭。

## 4. 组件开发规范

### 4.1 通用约定

- 方法首参使用 `context.Context`。
- 可变配置使用对应组件包的 `Option`。
- 错误使用 `%w` 包装底层错误。
- 导出类型、接口、函数必须有 GoDoc。
- 实现 `components.Typer` 可提供 DevOps/调试展示名：`GetType() string`。
- 如果组件自行精确触发 callback，可实现 `components.Checker` 并返回 `IsCallbacksEnabled() == true`。

### 4.2 ChatModel

接口位置：`components/model/interface.go`

- 基础模型实现 `model.BaseChatModel`：`Generate` 与 `Stream`。
- Agentic 模型实现 `model.AgenticModel`，消息类型为 `*schema.AgenticMessage`。
- 支持工具调用的模型优先实现 `model.ToolCallingChatModel.WithTools`。
- 不推荐新代码依赖 `ChatModel.BindTools`，该方式原地修改实例，存在并发风险。
- 模型实例应尽量无状态或只读；每次请求参数通过 Options 或 `WithTools` 派生实例传递。

### 4.3 Tool

接口位置：`components/tool/interface.go`

工具能力分层：

- `tool.BaseTool`：只提供 `Info(ctx)`，让模型知道工具 schema。
- `tool.InvokableTool`：非流式执行，`argumentsInJSON -> string`。
- `tool.StreamableTool`：流式执行，`argumentsInJSON -> StreamReader[string]`。
- `tool.EnhancedInvokableTool`：结构化/多模态工具，`ToolArgument -> ToolResult`。
- `tool.EnhancedStreamableTool`：结构化/多模态流式工具。

建议：

- 简单工具优先用 `components/tool/utils.NewTool`、`NewStreamTool`、`NewEnhancedTool`、`NewEnhancedStreamTool`。
- 工具名稳定、语义清晰，适合 LLM 选择。
- `ToolInfo` 参数 schema 必须完整描述字段、类型、必填项、枚举和说明。
- 工具内部避免全局可变状态；需要状态时通过 config、session 或 context 传入。

### 4.4 Prompt

接口位置：`components/prompt/interface.go`

- `ChatTemplate.Format(ctx, map[string]any, ...Option) ([]*schema.Message, error)`
- `AgenticChatTemplate.Format(ctx, map[string]any, ...Option) ([]*schema.AgenticMessage, error)`

建议：模板变量名保持稳定；Graph/Chain 中常用 `compose.WithOutputKey` 将上游输出转换为模板变量 map；系统提示词包含字面量 `{}` 时注意 FString 自动替换行为。

### 4.5 RAG 组件

接口位置：`components/document`、`components/embedding`、`components/indexer`、`components/retriever`。

- `document.Loader`：从 `Source.URI` 加载原始文档。
- `document.Transformer`：切分、过滤、合并、重排文档。
- `embedding.Embedder`：批量文本向量化。
- `indexer.Indexer`：存储文档并返回 ID。
- `retriever.Retriever`：按 query 检索相关文档。

约束：Indexer 与 Retriever 必须使用同一 Embedding 模型；Transformer 应保留已有 `Document.MetaData` 并 merge 新字段；Retriever 返回结果应按相关性排序。

## 5. 编排开发规范

### 5.1 选择方式

| 场景 | 推荐 |
| --- | --- |
| 线性流程 | `compose.Chain[I, O]` |
| 复杂 DAG、分支、并行、子图组合 | `compose.Graph[I, O]` |
| 带状态、多输入输出、业务步骤流 | `compose.Workflow[I, O]` |
| 单个函数节点 | Lambda Node |
| 让 Agent 调用确定性流程 | Graph/Chain 编译或包装为 GraphTool |

### 5.2 Graph/Chain 规范

- 使用 `compose.START` 和 `compose.END` 标识入口/出口。
- `AddEdge(start, end)` 前确保节点已添加。
- 节点 key 稳定、语义化，用于日志、callback、checkpoint、resume 地址。
- 图或链 `Compile` 后不可再修改。
- 类型不匹配时优先用 FieldMapping 或 Lambda 显式转换。
- 复杂图配置 graph name，便于 callback 与调试。
- Chain 适合模板 -> 模型 -> 解析 -> 工具/检索等线性流程，也可作为子图加入 Graph。

### 5.3 中断/恢复与检查点

- 需要人机协同、长任务、工具审批时使用中断/恢复。
- Runner 或 Graph 配置 checkpoint store 后可恢复。
- resume 目标依赖稳定地址，节点名、Agent 名、工具名不要随意变更。
- 中断点保存最小必要状态，状态类型应可序列化。

## 6. ADK Agent 开发规范

### 6.1 推荐入口

典型入口：

- `adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{...})`
- `adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent})`
- `runner.Query(ctx, query)` 或 `runner.Run(ctx, messages)`

`Runner` 输出 `AsyncIterator[*adk.AgentEvent]`，消费时循环 `Next()`，处理模型输出、工具输出、自定义输出、中断与错误事件。

### 6.2 ChatModelAgent 配置

关键字段：

- `Name`：Agent 名称；作为 AgentTool 或多 Agent 时必须稳定且唯一。
- `Description`：能力说明；作为工具/子 Agent 时帮助上层选择。
- `Instruction`：系统提示词，支持 session value 占位符。
- `Model`：必填，类型为 `model.BaseModel[M]`。
- `ToolsConfig`：工具列表、工具中间件、ReturnDirectly、内部事件透传。
- `GenModelInput`：自定义模型输入构造逻辑。
- `MaxIterations`：ReAct 最大循环次数，默认 20。
- `Handlers`：推荐的 Agent middleware 扩展方式。
- `ModelRetryConfig` / `ModelFailoverConfig`：模型重试与故障转移。

### 6.3 多 Agent 推荐模式

优先：

1. 用 `adk.NewAgentTool` 将专业 Agent 包装为 Tool。
2. 用 `adk/prebuilt/deep` 处理复杂任务拆解与子 Agent 协作。
3. 对确定性流程使用 Graph/Chain，再包装为工具交给 Agent 调用。

不优先推荐依赖全上下文共享的 Agent transfer。

### 6.4 Agent Handler 规范

- 新开发优先使用 `ChatModelAgentMiddleware` / `TypedChatModelAgentMiddleware`，不要新增依赖已废弃 `AgentMiddleware`。
- Handler 顺序：先注册的是外层 wrapper。
- 修改工具列表推荐在 `BeforeModelRewriteState` 修改 `state.ToolInfos` / `state.DeferredToolInfos`。
- 不建议在 `WrapModel` 临时注入工具，因为不会持久化到 state，且可能破坏 prompt cache。

## 7. Callback 规范

回调时机：`OnStart`、`OnEnd`、`OnError`、`OnStartWithStreamInput`、`OnEndWithStreamOutput`。

使用规范：

- 全局 callback 使用 `callbacks.AppendGlobalHandlers`，只在程序初始化阶段调用。
- 单次调用 callback 通过 `compose.WithCallbacks` 等 Option 注入。
- Handler 应实现 `TimingChecker`，避免不必要的 stream copy 和 goroutine 开销。
- Handler 不得修改输入/输出对象，避免竞态。
- Stream callback 获取的是 copy，读取后必须关闭。
- Handler 基于 `RunInfo.Component`、`RunInfo.Name`、`RunInfo.Type` 过滤，不依赖不同 Handler 之间的执行顺序。

## 8. 错误处理与并发规范

- 外部 IO、模型调用、工具调用必须接收并尊重 `context.Context` 取消。
- 错误包装使用 `fmt.Errorf("...: %w", err)`。
- 不吞掉 stream 中途错误，必须通过 `StreamReader` 传播给调用者。
- 组件实现尽量无状态；如需缓存，必须并发安全。
- 避免原地修改传入的 message、document、tool info；必要时复制后修改。
- 新代码优先使用 `WithTools` 或 Options，不使用有并发风险的 `BindTools`。

## 9. 测试与代码规范

来自仓库 README 与 CONTRIBUTING：

- 使用 `golangci-lint run ./...`。
- 代码格式符合 `gofmt -s`。
- import 顺序符合 `goimports`：标准库 -> 第三方 -> 本地。
- 导出的函数、接口、package 等需要 GoDoc 注释。
- 新功能应补充对应测试；组件接口通常使用 `go:generate mockgen` 生成 mock。
- PR 前需完整 lint 与测试。

## 10. 推荐项目开发模板

### 10.1 应用层目录建议

```text
app/
  agent/          # Agent 构造、Runner 封装、AgentTool
  graph/          # Graph/Chain/Workflow 编排
  tools/          # Tool 实现与 ToolInfo
  rag/            # Loader/Transformer/Indexer/Retriever 组合
  model/          # ChatModel 初始化与配置
  callback/       # tracing/logging/metrics handler
  schema/         # 业务 DTO，不重复 Eino schema
  config/         # 配置加载
```

### 10.2 开发流程

1. 定义业务输入/输出 DTO。
2. 选定模型实现，封装 `model.BaseChatModel` 初始化。
3. 定义工具，先写 `ToolInfo`，再实现执行逻辑。
4. 需要确定性流程时先用 `Chain`，复杂后升级为 `Graph`。
5. 需要自主决策时用 `ChatModelAgent`，把工具和确定性 GraphTool 暴露给 Agent。
6. 接入 callbacks 做日志、trace、metrics。
7. 对长任务/人工审批接入 checkpoint 与 resume。
8. 补充单测、流式测试、工具参数 schema 测试和错误路径测试。

## 11. 官网技术文档维护索引

官网入口：`https://www.cloudwego.cn/zh/docs/eino/`

指定实践页：`http://www.cloudwego.cn/zh/docs/eino/overview/bytedance_eino_practice/`

### 11.1 官网文档目录

| 一级目录 | 主题 |
| --- | --- |
| 概述 | 框架定位、开源背景、ADK、Agent/Graph 路线辨析、字节实践文章。 |
| 快速开始 | 从 `ChatModel`、`Message` 到 Agent、多轮会话、工具、Middleware、Callback、Interrupt、Graph Tool、Skill、A2UI、TurnLoop。 |
| Cookbook | 面向具体场景的实践示例和可复用方案。 |
| 核心模块 | Components、Chain/Graph/Workflow、Flow、ADK、应用开发工具链。 |
| 组件集成 | OpenAI、ARK、Document、Embedding、Tool、Callbacks、Indexer、Retriever、ChatTemplate 等外部组件集成。 |
| 发布记录 & 迁移指引 | v0.1 到 v0.9 的版本更新、不兼容变更与迁移说明。 |
| FAQ | 常见问题。 |

### 11.2 快速开始章节应纳入的能力

后续项目开发时，快速开始文档可作为功能落地顺序：

1. `ChatModel` 与 `Message`：先跑通最小模型调用。
2. `ChatModelAgent`、`Runner`、`AgentEvent`：再引入 Agent 与事件流。
3. `Memory` 与 `Session`：需要多轮上下文和持久化对话时接入。
4. `Tool` 与文件系统访问：让 Agent 具备外部动作能力。
5. `Middleware`：把横切扩展、工具动态注入、模型包装放到中间件。
6. `Callback` 与 Trace：接入可观测性。
7. `Interrupt/Resume`：处理人工确认、审批、暂停恢复。
8. `Graph Tool`：把复杂确定性流程包装给 Agent 调用。
9. `Skill Middleware`：通过技能包扩展 Agent 能力。
10. `A2UI` 协议：需要流式 UI 组件时使用。
11. `TurnLoop`：处理抢占、中止、多轮生命周期等高级控制。

### 11.3 核心模块文档清单

#### Components 组件

- `Document Loader` / `Document Parser`：加载与解析外部文档。
- `Embedding`：文本向量化。
- `Document Transformer`：切分、过滤、重排、合并文档。
- `Lambda`：将普通 Go 函数转为可编排节点。
- `Indexer`：写入索引或向量库。
- `Retriever`：召回相关上下文。
- `ChatTemplate`：将业务变量转换为模型消息。
- `ChatModel`：模型调用与流式输出。
- `ToolsNode & Tool`：工具声明、工具执行、工具节点。
- `AgenticModel`、`AgenticChatTemplate`、`AgenticToolsNode & Tool`：Agentic Runtime 相关 Beta 能力。

#### Chain & Graph & Workflow 编排

- `Chain/Graph` 编排介绍：理解节点、边、分支与执行拓扑。
- 编排设计理念：按数据流和控制流拆分业务流程。
- `Workflow` 编排框架：适合更业务化的状态流和多输入输出。
- 流式编程要点：重点理解 fake stream、concat、copy、merge。
- Callback 用户手册：统一观测组件和图运行。
- `CallOption` 能力与规范：按全局、组件类型、指定节点分发调用选项。
- `Interrupt & CheckPoint`：中断、检查点与恢复执行。

#### Flow 集成

- `ReAct Agent`：经典模型-工具循环。
- `Host Multi-Agent`：多 Agent 托管与协作。

#### ADK - Agent Development Kit

- Quickstart 与概述：ADK 基础使用和整体架构。
- Agent 抽象：理解 `Agent` 输入、输出、事件和动作。
- Agent 协作：多 Agent、AgentTool、子 Agent 协同。
- Agent 实现：`ChatModelAgent`、`Plan-Execute Agent`、`DeepAgents`。
- `ChatModel Failover`：模型故障转移。
- `Agent Runner` 与扩展：统一运行、事件消费、恢复。
- Human-in-the-loop：人在回路技术架构。
- `ChatModelAgentMiddleware`：模型、工具、上下文、状态的扩展点。
- Middleware：`FileSystem`、`Skill`、`Summarization`、`Reduction`、`PlanTask`、`ToolSearch`、`PatchToolCalls`、`AgentsMD`。
- Agent Callback：Agent 级可观测性。
- Agent Cancel 与 `TurnLoop`：取消、抢占、中止、多轮生命周期。

#### 应用开发工具链

- Eino Dev 插件安装。
- Eino Dev 可视化编排。
- Eino Dev 可视化调试。

### 11.4 官网“字节跳动 Eino 实践”要点

该实践页用“队员、战术、工具、独门秘笈”类比 Eino 开发模型，可沉淀为以下工程规则：

#### 组件选择规则

先确定“需要哪个组件抽象”，再确定“使用哪个具体实现”：

| 抽象 | 作用 | 常见实现方向 |
| --- | --- | --- |
| `ChatModel` | 与大模型交互，输入消息上下文，输出模型消息。 | OpenAI、Claude、Gemini、Ark、Ollama。 |
| `Tool` | 与外部世界交互，根据模型工具调用执行动作。 | 搜索、文件、Git、任务管理、业务系统 API。 |
| `Retriever` | 获取相关上下文，让模型输出基于事实。 | ElasticSearch、VikingDB、Redis VectorStore。 |
| `ChatTemplate` | 将外部输入转换成预设 prompt。 | `DefaultChatTemplate`。 |
| `Document Loader` | 加载指定文本。 | WebURL、S3、File。 |
| `Document Transformer` | 按规则转换文本。 | Markdown/HTML Splitter、Reranker。 |
| `Indexer` | 存储文档并建立索引。 | ElasticSearch、Redis、VikingDB。 |
| `Embedding` | 文本转向量，作为 Indexer/Retriever 共同依赖。 | OpenAI、Ark。 |
| `Lambda` | 用户自定义函数节点。 | JSON parser、DTO converter、branch condition。 |

#### 编排选择规则

- `Chain`：简单链式有向图，适合数据单向流动、无复杂分支的场景，例如 `ChatTemplate -> ChatModel`。
- `Graph`：灵活有向图，适合分支、工具调用、并行、复杂拓扑，例如“最多执行一次 ToolCall 的 Agent”。
- `Workflow`：面向业务流程、状态与多步骤编排。

图建模时统一使用：

- Node：组件实例或 Lambda。
- Edge：一对一数据流转。
- Branch：N 选 1 控制流转。
- START/END：统一入口与出口。

#### CallOption 分发规则

官网实践强调调用选项按范围分发：

- 全局生效：如 `WithCallbacks(handler)` 注入整个图。
- 按组件类型生效：如只对 ChatModel 节点设置 `model.WithTemperature(0.5)`。
- 按指定节点生效：如 `WithCallbacks(handler).DesignateNode("node_1")`。

工程建议：

- 模型参数、工具参数、callback、trace 配置都应通过 Option 注入。
- 避免在组件内部硬编码运行期参数。
- 节点名稳定是指定节点 Option、trace、resume 的基础。

#### 流式编程规则

官网实践强调开发者只需实现真实业务需要的流式范式，编排层负责自动处理：

1. 上游流式、下游非流式：自动 concat。
2. 上游非流式、下游流式：自动转为 `StreamReader[T]`。
3. 编排层处理流复制、合并、桥接等细节。

工程建议：

- ChatModel 通常实现 Invoke 与 Stream。
- Lambda 可按需要实现 Invoke、Stream、Collect、Transform 任一范式。
- 对 TTFT 敏感的路径应尽量保持全链路 Stream/Transform，不要中途 Collect。

### 11.5 官网 Eino 智能助手实践沉淀

官网实践给出一个 RAG ReAct Agent 场景：根据用户请求，从知识库检索必要信息，并按需调用工具完成任务。

#### 实践架构

分两阶段：

1. Knowledge Indexing：将领域知识切分、向量化、写入 VectorStore。
2. Eino Agent：请求进入后召回上下文，构造 prompt，交给 ReAct Agent 循环决策工具调用或最终回答。

#### Knowledge Indexing 流程

推荐 Graph：

```text
Document Source
  -> Document Loader
  -> Document Transformer / Markdown Splitter
  -> Embedding
  -> Indexer / VectorStore
  -> []string IDs
```

实践要点：

- 递归扫描指定目录下 Markdown 文件。
- 按标题或语义块切分文档，平衡向量化尺寸限制与召回效果。
- 使用同一个 Embedding 模型完成索引和查询。
- 可用 Redis Stack / RedisSearch 作为本地 VectorStore。
- 生成的 EinoDev 编排代码需要人工补全组件构造配置。

#### Eino Agent 流程

推荐 Graph：

```text
UserMessage
  -> Lambda: 提取 query / 转模板变量
  -> Retriever: 从 VectorStore 召回上下文
  -> ChatTemplate: 组合问题、上下文、历史
  -> ReAct Agent / ChatModelAgent
  -> Tool calls / Final Answer
```

可用工具示例：

- DuckDuckGo：互联网搜索。
- EinoTool：获取 Eino 工程信息、仓库链接、文档链接。
- GitClone：克隆指定仓库。
- TaskManager：添加、查看、删除任务。
- OpenURL：打开本地文件或 Web 链接。

工程建议：

- RAG 与 Tool Calling 结合时，Retriever 提供事实上下文，Tool 提供行动能力。
- ChatTemplate 中应明确区分用户问题、召回上下文、工具使用约束、输出格式。
- ReAct 循环需要设置最大迭代次数，避免工具循环失控。
- 将 Agent 封装成 HTTP 服务时，Runner 事件流应映射为前端可消费的增量输出。

#### 可观测性实践

官网实践提到可接入：

- APMPlus：查看 Trace 与 Metrics。
- Langfuse：查看请求 Trace。

工程建议：

- 所有 Graph、Agent、Tool 节点接入 callbacks。
- trace 中保留 graph name、node name、agent name、tool name。
- 对模型调用记录 token、latency、stream 首包时间、错误。
- 对 Retriever 记录 query、topK、score、命中文档 ID。
- 对 Tool 记录参数摘要、耗时、错误，不记录敏感明文。

### 11.6 官网版本与迁移文档维护

后续升级 Eino 时必须检查官网发布记录：

- `v0.1.*`：first release。
- `v0.2.*`：second release。
- `v0.3.*`：tiny break change。
- `v0.4.*`：compose optimization。
- `v0.5.*`：ADK implementation。
- `v0.6.*`：jsonschema optimization。
- `v0.7.*`：interrupt resume refactor。
- `v0.8.*`：adk middlewares，包含不兼容更新。
- `v0.9.*`：agentic-runtime，包含更新注意事项。

升级流程：

1. 先读目标版本发布记录和不兼容说明。
2. 检查 `components` 接口、`compose` 编排、ADK middleware、interrupt/resume、agentic-runtime 是否有破坏性变化。
3. 跑完整单测、集成测试、流式测试和人工中断恢复测试。
4. 对业务项目更新本 spec 中的版本注意事项。

## 12. 关键文件索引

- `README.zh_CN.md`：框架介绍、快速上手、代码规范。
- `doc.go`：根 package 定位。
- `schema/doc.go`：核心 schema 与流式约束。
- `schema/message.go`：传统消息模型。
- `schema/agentic_message.go`：Agentic 消息模型。
- `schema/document.go`：文档数据结构。
- `schema/tool.go`：工具元信息与参数 schema。
- `schema/stream.go`：流式读写实现。
- `components/types.go`：组件类型、`Typer`、callback checker。
- `components/model/interface.go`：模型接口。
- `components/tool/interface.go`：工具接口。
- `components/prompt/interface.go`：Prompt 模板接口。
- `components/retriever/interface.go`：检索接口。
- `components/indexer/interface.go`：索引接口。
- `components/embedding/interface.go`：Embedding 接口。
- `components/document/interface.go`：文档加载与转换接口。
- `compose/runnable.go`：`Runnable` 与四种执行范式。
- `compose/generic_graph.go`：`Graph` 泛型入口。
- `compose/chain.go`：`Chain` 链式编排。
- `compose/workflow.go`：`Workflow`。
- `compose/tool_node.go`：工具节点。
- `compose/checkpoint.go`：编排检查点。
- `callbacks/interface.go`：callback API 与约束。
- `adk/interface.go`：Agent 事件、动作、消息类型。
- `adk/chatmodel.go`：`ChatModelAgent`。
- `adk/runner.go`：`Runner`、run/resume/checkpoint。
- `adk/agent_tool.go`：Agent 作为 Tool。
- `adk/workflow.go`：Sequential/Parallel/Loop Agent。
- `adk/prebuilt/deep/deep.go`：DeepAgent。
