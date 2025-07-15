package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	tmp := t.TempDir()

	existingFile := filepath.Join(tmp, "exists.txt")
	err := os.WriteFile(existingFile, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	exists, err := FileExists(existingFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !exists {
		t.Errorf("expected true for existing file")
	}
}

func TestFileDoesNotExist(t *testing.T) {
	tmp := t.TempDir()

	nonExistent := filepath.Join(tmp, "missing.txt")
	exists, err := FileExists(nonExistent)
	if err != nil {
		t.Errorf("unexpected error for non-existent file: %v", err)
	}
	if exists {
		t.Errorf("expected false for non-existent file")
	}
}

func TestInvalidPath(t *testing.T) {
	invalidPath := string([]byte{0x00}) // null byte is illegal in paths
	exists, err := FileExists(invalidPath)
	if err == nil {
		t.Errorf("expected error for invalid path, got nil")
	}
	if exists {
		t.Errorf("expected false for invalid path")
	}
}

func TestFilesAreDifferentIdenticalFiles(t *testing.T) {
	tmp := t.TempDir()

	fileA := filepath.Join(tmp, "a.txt")
	fileB := filepath.Join(tmp, "b.txt")

	content := []byte("same content")
	if err := os.WriteFile(fileA, content, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fileB, content, 0644); err != nil {
		t.Fatal(err)
	}

	different, err := FilesAreDifferent(fileA, fileB)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if different {
		t.Errorf("expected files to be the same")
	}
}

func TestFilesAreDifferentDifferentFiles(t *testing.T) {
	tmp := t.TempDir()

	fileA := filepath.Join(tmp, "a.txt")
	fileB := filepath.Join(tmp, "b.txt")

	if err := os.WriteFile(fileA, []byte("some content"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(fileB, []byte("different content"), 0644); err != nil {
		t.Fatal(err)
	}
	different, err := FilesAreDifferent(fileA, fileB)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !different {
		t.Errorf("expected files to be different")
	}
}

func TestFilesAreDifferentMissingFile(t *testing.T) {
	tmp := t.TempDir()

	fileA := filepath.Join(tmp, "a.txt")

	missingFile := filepath.Join(tmp, "does-not-exist.txt")
	different, err := FilesAreDifferent(fileA, missingFile)
	if err == nil {
		t.Errorf("expected error for missing file B")
	}
	if different {
		t.Errorf("expected result to be false when file B is missing")
	}
}

func TestFilesAreDifferentMissingFiles(t *testing.T) {
	tmp := t.TempDir()
	missingFile := filepath.Join(tmp, "does-not-exist.txt")

	_, err := FilesAreDifferent(missingFile, missingFile)
	if err == nil {
		t.Errorf("expected error for missing files")
	}
}

func TestFilesAreDifferentUnreadableFile(t *testing.T) {
	tmp := t.TempDir()
	fileA := filepath.Join(tmp, "a.txt")
	protectedFile := filepath.Join(tmp, "protected.txt")
	os.WriteFile(protectedFile, []byte("hi"), 0000) // no permissions
	defer os.Chmod(protectedFile, 0644)             // cleanup permissions

	_, err := FilesAreDifferent(fileA, protectedFile)
	if err == nil {
		t.Errorf("expected error due to unreadable file")
	}
}
