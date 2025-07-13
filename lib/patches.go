package lib

import "slices"

type MissingReturn struct {
	File string
	Missing string
}

func CheckPatches(orig, patches []string) (result []MissingReturn) {
	for _, o := range orig {
		if !slices.Contains(patches, o) {
			result = append(result, MissingReturn{o, "patch"})
		}
	}

	for _, p := range patches {
		if !slices.Contains(orig, p) {
			result = append(result, MissingReturn{p, "original"})
		}
	}

	if result == nil {
		return []MissingReturn{}
	}

	return
}
