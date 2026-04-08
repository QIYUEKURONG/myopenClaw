package storage

import (
	"bufio"
	"encoding/json"
	"myopenclaw/types"
	"os"
	"path/filepath"
)

// 追加一条消息到历史文件（每次对话后调用）
func AppendMessage(sessionId string, msgs []types.LLMMessage) error {
	path, err := GetSessionDir()
	if err != nil {
		return err
	}
	storePath := filepath.Join(path, sessionId+".jsonl")

	file, err := os.OpenFile(storePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, msg := range msgs {
		writeVal, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		_, err = file.WriteString(string(writeVal) + "\n")
	}

	return err
}

// 加载某个 session 的全部历史消息
func LoadMessages(sessionId string) ([]types.LLMMessage, error) {
	path, err := GetSessionDir()
	if err != nil {
		return nil, err
	}
	storePath := filepath.Join(path, sessionId+".jsonl")
	file, err := os.OpenFile(storePath, os.O_RDONLY, 0600)
	if os.IsNotExist(err) {
		return make([]types.LLMMessage, 0), nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	result := make([]types.LLMMessage, 0)
	for scanner.Scan() {
		line := scanner.Text()
		var item types.LLMMessage
		err = json.Unmarshal([]byte(line), &item)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
