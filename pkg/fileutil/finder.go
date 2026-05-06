package fileutil

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func GetFilesByPaths(paths ...string) (Files, error) {
	var files Files
	for _, pathDir := range paths {
		err := filepath.Walk(pathDir, func(path string, info fs.FileInfo, err error) error {
			files = append(files, File{
				Name:    info.Name(),
				Size:    info.Size(),
				Path:    path,
				Parents: strings.Split(strings.TrimPrefix(strings.TrimSuffix(path, filepath.Base(path)), pathDir), string(filepath.Separator)),
			})
			return err
		})
		if err != nil {
			return files, err
		}
	}
	return files, nil
}
