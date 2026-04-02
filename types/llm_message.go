package types

// LLMMessage 是通用的 LLM 消息格式（抽象层）
type LLMMessage struct {
	Role       string     // "system" | "user" | "assistant" | "tool"
	Content    string     // 文本内容
	ToolCalls  []ToolCall // LLM 请求调用的工具列表
	ToolCallID string     // tool role 消息对应的 call_id
}

// ToolCall 是通用的工具调用格式（抽象层）
type ToolCall struct {
	ID        string                 // 工具调用的唯一 ID
	Name      string                 // 工具名称
	Arguments map[string]interface{} // 工具参数（已解析的 JSON）
}
