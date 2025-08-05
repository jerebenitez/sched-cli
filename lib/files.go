package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	godiffpatch "github.com/sourcegraph/go-diff-patch"
)

func Touch(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

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

func InstallFiles(pathToSrc string, files []string) error {
	for _, file := range files {
		path, f := filepath.Split(file)
		target := filepath.Join(pathToSrc, path)
		if err := os.MkdirAll(target, os.ModePerm); err != nil {
			return fmt.Errorf("error creating folder %s", target)
		}

		// TODO: Test that this behaves nicely when compiling the kernel
		filePath := filepath.Join(target, f)
		err := os.Link(filePath, target)
		if err != nil {
			return fmt.Errorf("error copying file %s", file)
		}
	}

	return nil
}

func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
