package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"myopenclaw/tools"
	"myopenclaw/types"
	"net/http"
	"os"
	"time"
)

type DeepSeekClient struct {
	APIKey  string
	BaseURL string
	Model   string
}

func NewDeepSeekClient() *DeepSeekClient {
	return &DeepSeekClient{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	}
}

type DeepSeekToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type DeepSeekMessage struct {
	//角色
	Role string `json:"role"`
	//内容
	Content string `json:"content"`
	//工具 为什么message里面会有tool呢 是因为role可能是tool
	ToolCalls []DeepSeekToolCall `json:"tool_calls"`
	//工具的ID 这个是为了让大模型确认一下已经调用过工具的顺序等。 比如{"123","wo"},{"456","ai"},{"789","ni"}
	ToolCallID string `json:"tool_call_id"`
}

type DeepSeekToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}
type DeepSeekTool struct {
	Type     string               `json:"type"`
	Function DeepSeekToolFunction `json:"function"`
}
type DeepSeekRequest struct {
	//要发送的消息
	Messages []DeepSeekMessage `json:"messages"`
	Model    string            `json:"model"`
	Stream   bool              `json:"stream"` // ← 添加这个
	Tools    []DeepSeekTool    `json:"tools"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message DeepSeekMessage `json:"message"`
	} `json:"choices"`
}

func convertToCommonLLMTool(ToolCalls []DeepSeekToolCall) []types.ToolCall {
	result := make([]types.ToolCall, 0, len(ToolCalls))

	for _, toolCall := range ToolCalls {
		var item types.ToolCall
		item.ID = toolCall.ID
		item.Name = toolCall.Function.Name
		var args map[string]interface{}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			args = make(map[string]interface{})
		}
		item.Arguments = args
		result = append(result, item)
	}

	return result
}

func convertToDeepSeekMessages(messages []types.LLMMessage) []DeepSeekMessage {
	if len(messages) == 0 {
		return nil
	}
	result := make([]DeepSeekMessage, 0, len(messages))

	for _, message := range messages {
		var item DeepSeekMessage
		item.Role = message.Role
		item.Content = message.Content

		if len(message.ToolCalls) > 0 {
			for _, toolCall := range message.ToolCalls {
				jsonStr, _ := json.Marshal(toolCall.Arguments)

				item.ToolCalls = append(item.ToolCalls, DeepSeekToolCall{
					ID:   toolCall.ID,
					Type: "function",
					Function: struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					}{
						Name:      toolCall.Name,
						Arguments: string(jsonStr),
					},
				})

			}
		}
		if message.ToolCallID != "" {
			item.ToolCallID = message.ToolCallID
		}

		result = append(result, item)
	}

	return result
}

func (d *DeepSeekClient) Chat(ctx context.Context, messages []types.LLMMessage, tools map[string]tools.Tool) (*LLMResponse, error) {
	//构建请求message
	var deepSeekRequest DeepSeekRequest
	deepSeekRequest.Messages = convertToDeepSeekMessages(messages)
	deepSeekRequest.Model = d.Model

	for _, tool := range tools {
		deepSeekRequest.Tools = append(deepSeekRequest.Tools, DeepSeekTool{
			Type: "function",
			Function: DeepSeekToolFunction{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters:  tool.Parameters(),
			},
		})
	}

	requestBody, err := json.Marshal(deepSeekRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deep seek request: %w", err)
	}

	req, err := http.NewRequest("POST", d.BaseURL+"/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.APIKey)

	// 调试信息
	fmt.Printf("[DEBUG] Request URL: %s\n", d.BaseURL+"/chat/completions")
	fmt.Printf("[DEBUG] Request Body: %s\n", string(requestBody))

	httpClient := &http.Client{Timeout: 10 * time.Second}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to send request, status code: %v", resp.StatusCode)
	}

	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var deepSeekResponse DeepSeekResponse
	err = json.Unmarshal(respbody, &deepSeekResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if len(deepSeekResponse.Choices) == 0 {
		return nil, fmt.Errorf("LLM no choices found")
	}

	return &LLMResponse{
		Content:   deepSeekResponse.Choices[0].Message.Content,
		ToolCalls: convertToCommonLLMTool(deepSeekResponse.Choices[0].Message.ToolCalls),
	}, nil
}
