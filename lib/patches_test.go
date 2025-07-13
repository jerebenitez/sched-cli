package lib

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Calling CheckPatches with the same files, should return an empty array
func TestCheckPatches(t *testing.T) {
	orig := []string{
		"sys/Makefile/Makefile",
		"sys/conf/files/files",
		"sys/kern/imgact_elf.c/imgact_elf.c",
		"sys/kern/kern_exec.c/kern_exec.c",
		"sys/kern/kern_thread.c/kern_thread.c",
		"sys/kern/sched_4bsd.c/sched_4bsd.c",
		"sys/sys/proc.h/proc.h",
	}
	patches := []string{
		"sys/Makefile/Makefile",
		"sys/conf/files/files",
		"sys/kern/imgact_elf.c/imgact_elf.c",
		"sys/kern/kern_exec.c/kern_exec.c",
		"sys/kern/kern_thread.c/kern_thread.c",
		"sys/kern/sched_4bsd.c/sched_4bsd.c",
		"sys/sys/proc.h/proc.h",
	}
	want := []MissingReturn{}

	got := CheckPatches(orig, patches)

	if !cmp.Equal(want, got) {
        t.Error(cmp.Diff(want, got))
    }
}

// Some missing files
func TestCheckPatchesMissingFiles(t *testing.T) {
	orig := []string{
		"sys/Makefile/Makefile",
		"sys/conf/files/files",
		"sys/sys/proc.h/proc.h",
	}
	patches := []string{
		"sys/kern/imgact_elf.c/imgact_elf.c",
		"sys/kern/kern_exec.c/kern_exec.c",
		"sys/kern/kern_thread.c/kern_thread.c",
		"sys/kern/sched_4bsd.c/sched_4bsd.c",
	}
	want := []MissingReturn{
		{"sys/Makefile/Makefile", "patch"},
		{"sys/conf/files/files", "patch"},
		{"sys/sys/proc.h/proc.h", "patch"},
		{"sys/kern/imgact_elf.c/imgact_elf.c", "original"},
		{"sys/kern/kern_exec.c/kern_exec.c", "original"},
		{"sys/kern/kern_thread.c/kern_thread.c", "original"},
		{"sys/kern/sched_4bsd.c/sched_4bsd.c", "original"},
	}

	got := CheckPatches(orig, patches)

	if !cmp.Equal(want, got) {
        t.Error(cmp.Diff(want, got))
    }
}

func TestCheckPatchesWithEmptySlice(t *testing.T) {
	orig := []string{
		"sys/Makefile/Makefile",
		"sys/conf/files/files",
		"sys/sys/proc.h/proc.h",
	}
	patches := []string{}
	want := []MissingReturn{
		{"sys/Makefile/Makefile", "patch"},
		{"sys/conf/files/files", "patch"},
		{"sys/sys/proc.h/proc.h", "patch"},
	}

	got := CheckPatches(orig, patches)

	if !cmp.Equal(want, got) {
        t.Error(cmp.Diff(want, got))
    }
}
