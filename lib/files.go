package lib

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	godiffpatch "github.com/sourcegraph/go-diff-patch"
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
	aData, err := os.ReadFile(pathA)
	if err != nil {
		return false, err
	}

	bData, err := os.ReadFile(pathB)
	if err != nil {
		return false, err
	}

	patches := godiffpatch.GeneratePatch("", string(aData), string(bData))

	return patches != "", nil
}

func TrimExtension(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}
