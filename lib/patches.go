package lib

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ApplyPatch(sourcePath, patchPath string, isDryRun bool) (bool, error) {
	if _, err := os.Stat(sourcePath); err != nil {
		return false, err
	}

	dryRun := ""
	if isDryRun {
		dryRun = "--dry-run"
	}

	cmd := exec.Command("patch", dryRun, "-p1", sourcePath)
	patchData, err := os.ReadFile(patchPath)
	if err != nil {
		return false, err
	}

	cmd.Stdin = bytes.NewReader(patchData)
	cmd.Stderr = nil
	cmd.Stdout = nil

	if err := cmd.Run(); err == nil {
		return true, nil
	} else {
		return false, nil
	}
}

func ApplyPatches(src, dir string, patches []string) error {
	for _, patch := range patches {
		sourcePath := filepath.Join(src, TrimExtension(patch))
		patchPath := filepath.Join(dir, patch)
		result, err := ApplyPatch(sourcePath, patchPath, true)
		if err != nil {
			return err
		} else if !result {
			return fmt.Errorf(
				"unable to patch %s: file not compatible with patch",
				TrimExtension(patch),
			)
		}
	}

	return nil
}
