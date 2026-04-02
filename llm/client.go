package llm

import (
	"context"
	"myopenclaw/tools"
	"myopenclaw/types"
)

type LLMClient interface {
	Chat(ctx context.Context, messages []types.LLMMessage, tools map[string]tools.Tool) (*LLMResponse, error)
}

type LLMResponse struct {
	Content   string           // LLM 的文本回复
	ToolCalls []types.ToolCall // LLM 请求调用的工具列表（通用格式）
}
