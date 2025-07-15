package lib

import (
	"fmt"
	"os"
	"os/exec"
)

func ApplyPatch(sourcePath, patchPath string, isDryRun bool) (bool, error) {
	if aExists, err := FileExists(sourcePath); err != nil || !aExists {
		return false, fmt.Errorf("source file missing or error: %w", err)
	}

	if bExists, err := FileExists(patchPath); err != nil || !bExists {
		return false, fmt.Errorf("patch missing or error: %w", err)
	}

	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return false, fmt.Errorf("failed to open /dev/null: %v", err)
	}
	defer devnull.Close()

	dryRun := ""
	if isDryRun {
		dryRun = "--dry-run"
	}

	// TODO: test that passing an empty string as an argument doesn't break
	diff := exec.Command(
		"patch", dryRun, "-p1", "--strip=0", 
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
