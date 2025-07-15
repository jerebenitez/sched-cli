package lib

import "io/fs"

func ReadRecursiveDir(fsys fs.FS) (files []string, err error) {
	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type().IsRegular() {
			files = append(files, path)
		}

		return nil
	})

	return
}

