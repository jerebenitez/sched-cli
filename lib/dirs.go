package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"
)

func ReadRecursiveDir(parentFolder, folder string) (files []string) {
	fileSystem := os.DirFS(parentFolder + "/" + folder)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if d.Type().IsRegular() {
			files = append(files, path)
		}

		return nil
	})

	return files
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

func CanApplyPatch(sourcePath, patchPath string) (bool, error) {
	if aExists, err := FileExists(sourcePath); err != nil || !aExists {
		return false, fmt.Errorf("source file missing or error: %w", err)
	}

	if bExists, err := FileExists(patchPath); err != nil || !bExists {
		return false, fmt.Errorf("patch missing or error: %w", err)
	}

	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		log.Fatalf("Failed to open /dev/null: %v", err)
	}
	defer devnull.Close() // Ensure /dev/null is closed when done

	diff := exec.Command(
		"patch", "--dry-run", "-p1", "--strip=0", 
		sourcePath, 
		patchPath,
	)
	// We only care about exit code
	diff.Stdout = devnull

	if err := diff.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}

	return true, nil

}

func TrimExtension(path string) string {
	if idx := strings.LastIndex(path, "."); idx > 0 {
		return path[0:idx]
	}

	return path
}
