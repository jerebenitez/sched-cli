package lib

import (
	"slices"
	"testing"
)

type inputs struct {
	a []string
	b []string
}

var comparetests = []struct {
	name string
    in inputs
    out []string
}{
	{"Empty slices", inputs{}, []string{}},
	{"Equal slices", inputs{
		[]string{"sys/sys/proc.h"},
		[]string{"sys/sys/proc.h"},
	}, []string{}},
	{"Missing patch", inputs{
		[]string{"sys/sys/proc.h"},
		[]string{},
	}, []string{"sys/sys/proc.h"}},
	{"Missing original", inputs{
		[]string{},
		[]string{"sys/sys/proc.h"},
	}, []string{}},
}

func TestCompareDirs(t *testing.T) {
	for _, tt := range comparetests {
        t.Run(tt.name, func(t *testing.T) {
            s := CompareDirs(tt.in.a, tt.in.b)
            if !slices.Equal(s, tt.out) {
                t.Errorf("got %q, want %q", s, tt.out)
            }
        })
    }
}
