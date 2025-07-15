package lib

import (
	"slices"
	"testing"
	"testing/fstest"
)

type recursiveOutput struct {
	dirs []string
	err error
}

func TestReadRecursiveDir(t *testing.T) {
	recursivetests := map[string]struct {
		fsys fstest.MapFS
		output recursiveOutput
	} {
		"Empty directory": {
			fsys: fstest.MapFS{},
			output: recursiveOutput{
				dirs: []string{},
				err: nil,
			},
		},
		"Real example": {
			fsys: fstest.MapFS{
				"sys/conf/files": {},
				"sys/kern/imgact_elf.c": {},
				"sys/kern/kern_exec.c": {},
				"sys/kern/kern_thread.c": {},
				"sys/kern/sched_4bsd.c": {},
				"sys/Makefile": {},
				"sys/sys/proc.h": {},
			},
			output: recursiveOutput{
				dirs: []string{
					"sys/conf/files",
					"sys/kern/imgact_elf.c",
					"sys/kern/kern_exec.c",
					"sys/kern/kern_thread.c",
					"sys/kern/sched_4bsd.c",
					"sys/Makefile",
					"sys/sys/proc.h",
				},
				err: nil,
			},
		},
	}

	for name, test := range recursivetests {
		t.Run(name, func(t *testing.T) {

			dirs, err := ReadRecursiveDir(test.fsys)
			if slices.Equal(dirs, test.output.dirs) && err != test.output.err {
				t.Fatalf("ReadRecursiveDir(%v) \n\treturned %q, %q;\n\texpected %q, %q",
				test.fsys,
				dirs, err,
				test.output.dirs, test.output.err,
			)
			}
		})
	}
}
