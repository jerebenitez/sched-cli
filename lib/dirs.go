package lib

import (
	"io/fs"
	"log"
	"os"
)

func ReadRecursiveDir(parentFolder, folder string) (files []string) {
	fileSystem := os.DirFS(parentFolder + "/" + folder)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if d.Type().IsRegular() {
			files = append(files, path)
		}

		return nil
	})

	return files
}

