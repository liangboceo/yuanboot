# Agent Tools 开发示范 Spec

本文用于指导 `sendex-server/pkg/agent/tools` 下新增 Agent Tool，确保工具编码、注册、入参 schema、执行 handler、上下文回调、结果结构和验证方式保持一致。

## 1. 适用范围

适用于以下场景：

- 新增一个可被 Agent 调用的工具，例如生成文档、生成 PPT、生成 PRD、联网搜索。
- 将工具展示到系统 Agent 配置页，允许选择工具编码。
- 将工具接入 `SessionTools`，供会话和工单处理链路调用。

## 2. 目录约定

```text
sendex-server/pkg/agent/tools/
├── registry.go      # 工具编码、工具展示定义、按编码选择工具
├── session.go    # 工具上下文 Context 与 SessionTools 聚合入口
├── helpers.go       # 通用参数解析辅助函数
├── presentation.go  # PPT 工具示例
├── prd.go           # PRD 工具示例
└── websearch.go     # 联网搜索工具示例
```

新增工具优先独立成文件：

```text
pkg/agent/tools/<feature>.go
```

除非工具非常简单且与已有工具强相关，否则不要继续堆到 `session.go`。

## 3. Tool 基础结构

项目统一使用 `sendex-server/pkg/agent.Tool` 定义工具：

```go
type Tool struct {
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	InputSchema     map[string]any `json:"inputSchema"`
	ReadOnlyHint    bool           `json:"readOnlyHint"`
	DestructiveHint bool           `json:"destructiveHint"`
	IdempotentHint  bool           `json:"idempotentHint"`
	OpenWorldHint   bool           `json:"openWorldHint"`
	Handler         ToolHandler    `json:"-"`
}
```

约定：

- `Name` 必须使用 `registry.go` 中定义的 `ToolCodeXxx` 常量，不要写裸字符串。
- `Description` 必须描述工具能力、适用场景和产出，便于 LLM 正确选择。
- `InputSchema` 使用 JSON Schema 风格，必须包含 `type`、`properties`、`required`。
- `Handler` 只做确定性执行、参数校验和结果返回，不在工具内直接访问 Controller。
- 所有底层错误使用 `%w` 包装，便于排查。

## 4. 工具编码规范

在 `pkg/agent/tools/registry.go` 中新增工具编码：

```go
const (
	ToolCodeGeneratePRD = "sendex_generate_prd"
)
```

命名规则：

| 类型 | 规则 | 示例 |
|---|---|---|
| Go 常量 | `ToolCode` + 动作/能力名 | `ToolCodeGeneratePRD` |
| 工具 code | `sendex_` + snake_case 能力名 | `sendex_generate_prd` |
| 工具函数 | 动宾结构 + `Tool` 后缀 | `GeneratePRDTool` |
| 工具文件 | snake_case 或功能短名 | `prd.go` |

要求：

- 工具 code 必须稳定，已上线后不要随意改名，否则已有 Agent 配置中的 `toolCodes` 会失效。
- 新工具如需在系统配置页可选，必须加入 `Definitions()`。

## 5. 工具展示注册

在 `Definitions()` 中增加展示项：

```go
func Definitions() []dto.SysAgentToolResp {
	return []dto.SysAgentToolResp{
		{Code: ToolCodeGeneratePRD, Name: "生成PRD", Description: "生成产品需求文档并保存到资料库"},
	}
}
```

要求：

- `Code` 必须使用 `ToolCodeXxx`。
- `Name` 面向用户，简短清晰。
- `Description` 面向配置页，说明工具价值。
- 暂不开放给配置页的工具可以保留编码和实现，但不要加入 `Definitions()`，或明确用注释说明原因。

## 6. 工具接入 SessionTools

在 `pkg/agent/tools/session.go` 的 `SessionTools(ctx Context)` 中加入工具实例：

```go
func SessionTools(ctx Context) []agent.Tool {
	return []agent.Tool{
		GeneratePRDTool(ctx.SaveDoc, ctx.SessionID),
	}
}
```

要求：

- 需要保存资料库文档的工具，通过 `ctx.SaveDoc` 注入，不直接操作数据库。
- 需要读取资料库文档的工具，通过 `ctx.ReadDoc` 注入。
- 需要读取会话历史的工具，通过 `ctx.GetMessages` 注入。
- 工具内部不应依赖全局 session 状态，使用闭包传入的 `sessionID`。

## 7. Context 回调规范

当前 `tools.Context` 定义：

```go
type Context struct {
	SessionID   uint
	SessionName string
	GetMessages GetSessionMessagesFunc
	SaveDoc     SaveDocumentFunc
	ReadDoc     ReadDocumentFunc
}
```

工具与业务系统交互优先通过这些回调完成：

| 回调 | 场景 |
|---|---|
| `GetMessages` | 读取当前会话历史 |
| `SaveDoc` | 生成文档、PPT、PRD 后保存到资料库 |
| `ReadDoc` | 根据文档 ID 读取资料库文档内容 |

如果新增工具需要新的系统能力，应优先扩展 `Context`，再由 `session_controller.go` 和 `workorder_processor.go` 注入实现。

## 8. InputSchema 编写规范

示例：

```go
InputSchema: map[string]any{
	"type": "object",
	"properties": map[string]any{
		"filename": map[string]any{"type": "string", "description": "PRD文件名，如 product-prd.md"},
		"title":    map[string]any{"type": "string", "description": "PRD标题或产品/功能名称"},
		"features": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{"type": "string", "description": "功能名称"},
				},
				"required": []string{"name"},
			},
			"description": "核心功能列表",
		},
	},
	"required": []string{"title", "features"},
}
```

要求：

- 字段名使用 camelCase，例如 `targetUsers`、`nonFunctionalReqs`。
- 所有字段必须写 `description`，让模型知道如何填参。
- 数组必须声明 `items`。
- 对象数组必须声明内部 `properties` 和必要的 `required`。
- 必填字段只放真正执行不可缺少的参数；可推导字段使用默认值。

## 9. Handler 开发规范

标准结构：

```go
Handler: func(execCtx context.Context, args map[string]any) (any, error) {
	filename := StringArg(args, "filename")
	content := StringArg(args, "content")
	if content == "" {
		return nil, fmt.Errorf("文档内容不能为空")
	}

	docID, meta, err := saveDoc(filename, content, sessionID)
	if err != nil {
		return nil, fmt.Errorf("保存文档失败: %w", err)
	}

	result := map[string]any{
		"docId":    docID,
		"filename": filename,
		"message":  fmt.Sprintf("文档「%s」已生成并保存到资料库，文档ID: %s", filename, docID),
	}
	for k, v := range meta {
		result[k] = v
	}
	return result, nil
}
```

要求：

- 参数读取优先使用 `StringArg`、`NumberArg` 等辅助函数。
- 必填参数在 handler 内再次校验，不能只依赖 schema。
- 文件名、扩展名、空内容、数组为空等边界必须处理。
- 返回 `map[string]any`，字段尽量稳定。
- 生成文档类工具必须返回 `docId`、`filename`、`message`。
- 如返回完整内容，字段统一使用 `content`。
- 不要在工具中写 HTTP 响应、SSE、数据库事务；这些能力通过回调或上层控制器完成。

## 10. 文件生成类工具规范

文件生成类工具包括：

- `sendex_generate_document`
- `sendex_generate_presentation`
- `sendex_generate_prd`

`sendex_generate_prd` 当前约定生成 HTML 格式 PRD，并将 `index.html`、`README.txt` 等完整产物打包为 ZIP 压缩包保存；下载时应解码为真实 `.zip` 文件。

通用要求：

- 保存统一使用 `SaveDocumentFunc`：

```go
type SaveDocumentFunc func(filename, content string, sessionID uint) (docID string, docMeta map[string]any, err error)
```

- 文本/Markdown 工具直接保存文本内容。
- 二进制文件需要编码为可保存字符串，例如 PPTX 使用 `data:application/...;base64,` 前缀。
- 返回结果中合并 `docMeta`，但不要覆盖核心语义字段，必要时跳过冲突字段。
- 默认文件名必须带扩展名。

## 11. 示例：生成 PRD 工具

### 11.1 工具编码

```go
const (
	ToolCodeGeneratePRD = "sendex_generate_prd"
)
```

### 11.2 展示定义

```go
{Code: ToolCodeGeneratePRD, Name: "生成PRD", Description: "生成产品需求文档并保存到资料库"}
```

### 11.3 聚合入口

```go
GeneratePRDTool(ctx.SaveDoc, ctx.SessionID),
```

### 11.4 工具实现骨架

```go
func GeneratePRDTool(saveDoc SaveDocumentFunc, sessionID uint) agent.Tool {
	return agent.Tool{
		Name:        ToolCodeGeneratePRD,
		Description: "生成产品需求文档(PRD)并保存到资料库。",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title":      map[string]any{"type": "string", "description": "PRD标题"},
				"background": map[string]any{"type": "string", "description": "项目背景"},
				"features":   map[string]any{"type": "array", "items": map[string]any{"type": "object"}},
			},
			"required": []string{"title", "background", "features"},
		},
		ReadOnlyHint:    false,
		DestructiveHint: false,
		IdempotentHint:  false,
		OpenWorldHint:   false,
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			// 1. 参数校验
			// 2. 构造 HTML
			// 3. 打包 ZIP
			// 4. 调用 saveDoc 保存 data:application/zip;base64,...
			// 5. 返回 docId / filename / content / message
			return nil, nil
		},
	}
}
```

## 12. 示例：联网搜索工具

联网搜索类工具与文件生成类不同，不应保存资料库文档。

要求：

- `ReadOnlyHint` 设置为 `true`。
- `OpenWorldHint` 设置为 `true`。
- HTTP 请求必须带 `context.WithTimeout`。
- 外部接口失败时返回清晰错误。
- 返回结果应包含 `query`、`source`、`summary`、`results`。

示例返回：

```go
map[string]any{
	"query":   query,
	"source":  "DuckDuckGo Instant Answer",
	"summary": summary,
	"results": results,
	"message": fmt.Sprintf("找到 %d 条与「%s」相关的搜索结果", len(results), query),
}
```

## 13. 工具选择链路

工具选择链路如下：

```text
SysAgent.ToolCodes
  -> utils.UnmarshalStringSlice
  -> tools.Select(ctx, toolCodes)
  -> SessionTools(ctx)
  -> agent.StreamChatWithTools(...)
  -> Agent.convertTools(...)
  -> Eino ADK ToolsNode 执行 Handler
```

相关入口：

- `internal/controller/session_controller.go`
- `internal/service/workorder_processor.go`
- `internal/controller/sys_controller.go`

新增工具后至少确认：

- `tools.Select` 能按工具 code 选中该工具。
- `Definitions()` 能在系统工具选项中返回该工具。
- `SessionTools` 返回的工具 `Name` 与 `ToolCodeXxx` 完全一致。

## 14. 验证规范

新增或修改工具后执行：

```bash
cd sendex-server
gofmt -w pkg/agent/tools/<feature>.go pkg/agent/tools/registry.go pkg/agent/tools/session.go
go test ./pkg/agent/tools ./pkg/agent
```

如影响 Controller、Service 或 DTO，再补充相关包测试：

```bash
go test ./internal/controller ./internal/service ./internal/dto
```

## 15. 常见问题

### 15.1 工具配置页看不到新工具

检查：

- 是否在 `registry.go` 新增 `ToolCodeXxx`。
- 是否在 `Definitions()` 中加入展示项。
- 前端是否重新拉取 Agent options。

### 15.2 Agent 选择了工具但执行失败

检查：

- `SessionTools` 是否加入工具实例。
- 工具 `Name` 是否等于配置中的 `toolCodes`。
- `InputSchema.required` 是否与 handler 校验一致。
- handler 是否对空值、类型转换失败做了保护。

### 15.3 生成文档后无法下载或预览

检查：

- 是否通过 `SaveDocumentFunc` 保存。
- `filename` 是否包含正确扩展名。
- 文本内容是否直接保存，二进制内容是否正确 base64/data URL 编码。
- 返回结果是否包含 `docId`。

## 16. 开发 Checklist

新增工具时逐项确认：

- [ ] 新增 `pkg/agent/tools/<feature>.go`。
- [ ] 工具函数返回 `agent.Tool`。
- [ ] `Name` 使用 `ToolCodeXxx` 常量。
- [ ] `Description` 清晰描述能力和场景。
- [ ] `InputSchema` 完整声明字段、类型、必填项和描述。
- [ ] handler 内完成参数校验和默认值处理。
- [ ] 需要系统能力时通过 `tools.Context` 回调注入。
- [ ] 在 `registry.go` 新增工具编码。
- [ ] 需要展示时加入 `Definitions()`。
- [ ] 在 `SessionTools` 中加入工具实例。
- [ ] 执行 `gofmt`。
- [ ] 执行 `go test ./pkg/agent/tools ./pkg/agent`。
