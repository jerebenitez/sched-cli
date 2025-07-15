package lib

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

func FileExists(path string) (bool, error){
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func FilesAreDifferent(pathA, pathB string) (bool, error) {
	if aExists, err := FileExists(pathA); err != nil || !aExists {
		return false, fmt.Errorf("file A missing or error: %w", err)
	}

	if bExists, err := FileExists(pathB); err != nil || !bExists {
		return false, fmt.Errorf("file B missing or error: %w", err)
	}

	aData, err := os.ReadFile(pathA)
	if err != nil {
		return false, err
	}

	bData, err := os.ReadFile(pathB)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(aData, bData) {
		return true, nil
	}

	return false, nil
}

func TrimExtension(path string) string {
	if idx := strings.LastIndex(path, "."); idx > 0 {
		return path[0:idx]
	}

	return path
}
