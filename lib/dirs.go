package lib

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ReadRecursiveDir(fsys fs.FS) (files []string, err error) {
	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type().IsRegular() {
			files = append(files, path)
		}

		return nil
	})

	return
}

func IsGitRepo(dir string) bool {
	gitPath := filepath.Join(dir, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}

func ResolvePath(dir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var path string
	if strings.HasPrefix(dir, "~/") {
		path = filepath.Join(homeDir, dir[2:])	
	} else if dir == "~" {
		path = homeDir
	} else {
		path = dir
	}

	return path, nil
}
