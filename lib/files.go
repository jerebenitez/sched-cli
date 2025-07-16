package lib

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
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

	diff := diffmatchpatch.New()
	diffs := diff.DiffMain(string(aData), string(bData), false)

	areDifferent := len(diffs) > 1 || diffs[0].Type != diffmatchpatch.DiffEqual

	return areDifferent, nil
}

func TrimExtension(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}
