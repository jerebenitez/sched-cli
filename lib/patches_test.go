package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyPatchDryRun(t *testing.T) {
	tmp := t.TempDir()

	original := filepath.Join(tmp, "file.txt")
	patch := filepath.Join(tmp, "change.patch")

	// Create original file
	err := os.WriteFile(original, []byte("line1\nline2\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid unified diff patch
	patchContent := `--- file.txt
+++ file.txt
@@ -1,2 +1,2 @@
-line1
+line1 patched
 line2
`
	err = os.WriteFile(patch, []byte(patchContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test dry-run
	applied, err := ApplyPatch(original, patch, true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !applied {
		t.Errorf("expected patch to be applicable in dry-run")
	}
}

func TestApplyPatch_FileMissing(t *testing.T) {
	tmp := t.TempDir()

	patch := filepath.Join(tmp, "patch.diff")
	os.WriteFile(patch, []byte("irrelevant"), 0644)

	_, err := ApplyPatch("nonexistent.txt", patch, true)
	if err == nil {
		t.Errorf("expected error for missing source file")
	}
}

func TestApplyPatch_InvalidPatch(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.txt")
	patch := filepath.Join(tmp, "bad.patch")

	os.WriteFile(src, []byte("abc\n"), 0644)
	os.WriteFile(patch, []byte("this is not a patch"), 0644)

	applied, err := ApplyPatch(src, patch, true)
	if err == nil && applied {
		t.Errorf("expected invalid patch to fail")
	}
}
