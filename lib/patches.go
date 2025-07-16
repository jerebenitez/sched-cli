package lib

import (
	"bytes"
	"fmt"
	"os"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

func ApplyPatch(sourcePath, patchPath string, isDryRun bool) (bool, error) {
	patchFile, err := os.Open(patchPath)
	if err != nil {
		return false, fmt.Errorf("open patch file: %w", err)
	}
	defer patchFile.Close()

	patches, _, err := gitdiff.Parse(patchFile)
	if err != nil {
		return false, fmt.Errorf("parse patch file: %w", err)
	}

	if len(patches) == 0 {
		return false, fmt.Errorf("no files to patch")
	}

	for _, patch := range patches {
		if patch.NewName == "/dev/null" || patch.NewName == "" {
			continue
		}

		content, err := os.ReadFile(sourcePath)
		if err != nil {
			return false, fmt.Errorf("read target file: %w", err)
		}

		var out bytes.Buffer
		err = gitdiff.Apply(&out, bytes.NewReader(content), patch)
		if err != nil {
			return false, nil
		}
		
		if !isDryRun {
			if err := os.WriteFile(sourcePath, out.Bytes(), 0644); err != nil {
				return false, fmt.Errorf("write patched file %s: %w", sourcePath, err)
			}
		}

	}

	return true, nil
}
