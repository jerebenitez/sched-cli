package lib

import (
	"bytes"
	"os"
	"os/exec"
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
