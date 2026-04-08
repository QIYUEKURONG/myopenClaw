package types

// LLMMessage 是通用的 LLM 消息格式（抽象层）
type LLMMessage struct {
	Role       string     `json:"role"`       // "system" | "user" | "assistant" | "tool"
	Content    string     `json:"content"`    // 文本内容
	ToolCalls  []ToolCall `json:"toolCalls"`  // LLM 请求调用的工具列表
	ToolCallID string     `json:"toolCallId"` // tool role 消息对应的 call_id
}

// ToolCall 是通用的工具调用格式（抽象层）
type ToolCall struct {
	ID        string                 `json:"id"`        // 工具调用的唯一 ID
	Name      string                 `json:"name"`      // 工具名称
	Arguments map[string]interface{} `json:"arguments"` // 工具参数（已解析的 JSON）
}
