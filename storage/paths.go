package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetSessionDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Errorf("Failed to get user home dir: %v", err)
		return "", err
	}

	//session
}
