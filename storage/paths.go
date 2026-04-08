package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetSessionDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Failed to get user home dir: %v", err)
	}
	current := filepath.Join(dir, ".myopenclaw", "sessions")

	return current, nil
}
