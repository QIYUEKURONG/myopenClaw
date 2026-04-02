package tools

import "fmt"

type EchoTool struct{}

func (e *EchoTool) Name() string {
	return "echo"
}

func (e *EchoTool) Description() string {
	return `Echo Tool :回显输入的内容
		参数：
			- message(string)：要回显的消息
		示例：
			{"message":"Hello World"}
		返回：
		    回显的消息内容
`
}

func (e *EchoTool) Execute(args map[string]interface{}) (string, error) {
	if _, ok := args["message"]; !ok {
		return "", fmt.Errorf("must provide a message")
	}
	output, ok := args["message"].(string)
	if !ok {
		return "", fmt.Errorf("must provide a type of string")
	}

	return output, nil
}

func (e *EchoTool) Parameters() map[string]interface{} {
	ParameterMap := map[string]interface{}{
		"type":     "object",
		"required": []string{"message"},
		"properties": map[string]interface{}{
			"message": map[string]interface{}{"type": "string", "description": "要回显的消息"},
		},
	}

	return ParameterMap
}
