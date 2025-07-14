package lib

import (
	"errors"
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

func FileExists(path string) bool{
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		log.Fatal(err)
		return false
	}	
}

func FilesAreDifferent(pathA, pathB string) bool {
	if !FileExists(pathA) || !FileExists(pathB) {
		return false
	}

	// TODO: Migrate this entirely to go
	diff := exec.Command("diff", "-q", pathA, pathB)

	if err := diff.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return true
			}
			log.Fatal(err)
			return false
		} else {
			return false
		}
	}

	return false
}

func CanApplyPatch(sourcePath, patchPath string) bool {
	if !FileExists(sourcePath) || !FileExists(patchPath) {
		return false
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
		// Consider any error as a condition for not applicability
		// TODO: maybe improve this
		return false
	}

	return true

}

func TrimExtension(path string) string {
	if idx := strings.LastIndex(path, "."); idx > 0 {
		return path[0:idx]
	}

	return path
}
