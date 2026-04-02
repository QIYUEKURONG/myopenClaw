package agent

import (
	"context"
	"errors"
	"fmt"
	"myopenclaw/constant"
	"myopenclaw/llm"
	"myopenclaw/tools"
	"myopenclaw/types"
	"strings"
	"time"
)

type Runtime struct {
	//llm 的client
	//sessionID 对应 llm client
	//SessionToLlm map[string]*LLMClient
	Tools map[string]tools.Tool

	LLMClient llm.LLMClient
}

func NewRuntime() *Runtime {
	echoTool := tools.EchoTool{}

	toolsMap := make(map[string]tools.Tool)
	toolsMap[echoTool.Name()] = &echoTool
	deepSeek := llm.NewDeepSeekClient()

	return &Runtime{Tools: toolsMap, LLMClient: deepSeek}
}

func (r *Runtime) buildSystemPrompt() string {
	// 你来实现：
	// 1. Agent 的身份："你是一个 AI 助手，名字叫 MyOpenClaw"
	// 2. 工具列表：遍历所有工具，添加它们的 Name 和 Description
	// 3. 工作原则：安全、可维护、高性能
	var prompt strings.Builder
	//身份
	prompt.WriteString("你是一个 AI 助手，名字叫 MyOpenClaw.\n\n")

	//工具
	prompt.WriteString("## 你拥有的工具列表有这些：\n\n")

	for _, tool := range r.Tools {
		prompt.WriteString(fmt.Sprintf("### %s\n", tool.Name()))
		prompt.WriteString(tool.Description())
		prompt.WriteString("\n\n")
	}

	prompt.WriteString("## 工作原则\n\n")
	prompt.WriteString("a. 安全第一：不执行危险操作\n")
	prompt.WriteString("b. 可维护性：代码清晰易懂\n")
	prompt.WriteString("c. 高性能：合理利用资源\n")

	return prompt.String()

}

func (r *Runtime) ProcessMessage(ctx context.Context, msg *types.Message) (*types.Response, error) {
	// 你来实现：
	// 1. 构建 System Prompt
	// 2. 调用 LLM（暂时先返回假数据，因为还没有 LLM Client）
	// 3. 返回 Response
	prompt := r.buildSystemPrompt()
	if prompt == "" {
		return nil, errors.New("no prompt")
	}

	llmMessages := []types.LLMMessage{
		{Role: constant.SysRole, Content: prompt},
		{Role: constant.UserRole, Content: msg.Content},
	}

	tryCount := 5
	for i := 0; i < tryCount; i++ {
		//通过模型找到对应的tool 的结构体 现在只有一个deepseek 就写一个吧
		chatResult, err := r.LLMClient.Chat(ctx, llmMessages, r.Tools)
		if err != nil {
			return nil, fmt.Errorf("LLMClient deal chat find error  %v", err)
		}

		if len(chatResult.ToolCalls) == 0 {
			//得到了最后的结果
			return &types.Response{
				SessionID:   msg.SessionID,
				Content:     chatResult.Content,
				CreatedTime: time.Now(),
			}, nil
		}
		llmMessages = append(llmMessages, types.LLMMessage{
			Role:      "assistant",
			Content:   chatResult.Content,
			ToolCalls: chatResult.ToolCalls,
		})

		//执行工具
		for _, tc := range chatResult.ToolCalls {
			toolResult, err := r.executeTool(&tc)
			if err != nil {
				// 工具执行失败，把错误信息作为工具结果返回给 LLM
				toolResult = fmt.Sprintf("Error: %v", err)
			}

			llmMessages = append(llmMessages, types.LLMMessage{
				Role:       "tool",
				Content:    toolResult,
				ToolCallID: tc.ID,
			})
		}

	}

	return nil, fmt.Errorf("exceeded max rounds")
}

func (r *Runtime) executeTool(useTool *types.ToolCall) (string, error) {
	if useTool == nil {
		return "", fmt.Errorf("useTool is nil")
	}

	tool, exists := r.Tools[useTool.Name]
	if !exists {
		return "", fmt.Errorf("tool '%s' not found", useTool.Name)
	}

	result, err := tool.Execute(useTool.Arguments)
	if err != nil {
		return "", fmt.Errorf("execute tool error: %v", err)
	}

	return result, nil
}
