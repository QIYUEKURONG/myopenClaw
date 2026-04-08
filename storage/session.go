package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// 从 sessions.json 加载索引
func LoadSessionIndex() (map[string]string, error) {
	path, err := GetSessionDir()
	if err != nil {
		return nil, err
	}
	sessionIndexFile := filepath.Join(path, "sessions.json")

	filePath, err := os.ReadFile(sessionIndexFile)
	if os.IsNotExist(err) {
		return make(map[string]string), nil
	}
	if err != nil {
		return nil, fmt.Errorf("Error loading session index: %v", err)
	}

	result := make(map[string]string)

	err = json.Unmarshal(filePath, &result)
	if err != nil {
		return nil, fmt.Errorf("LoadSessionIndex unmarshal error %v", err)
	}

	return result, nil
}

// 把索引保存到 sessions.json
func SaveSessionIndex(data map[string]string) error {
	path, err := GetSessionDir()
	if err != nil {
		return err
	}
	sessionIndexFile := filepath.Join(path, "sessions.json")
	//judage
	err = os.MkdirAll(filepath.Dir(sessionIndexFile), 0700)
	if err != nil {
		return fmt.Errorf("SaveSessionIndex create error %v", err)

	}

	writeData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("SaveSessionIndex marshal error %v", err)

	}

	//write
	err = os.WriteFile(sessionIndexFile, writeData, 0600)
	if err != nil {
		return fmt.Errorf("SaveSessionIndex write error %v", err)
	}

	return nil
}
