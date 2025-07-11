package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetFilesByPaths(paths ...string) (Files, error) {
	files := []File{}

	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read information of %s, %v", path, err)
		}

		if !fileInfo.IsDir() {
			files = append(files, File{
				Name:    fileInfo.Name(),
				Size:    fileInfo.Size(),
				Path:    path,
				Parents: nil,
			})
			continue
		}

		dirFiles, err := readDirectoryFilesRecursively(path, filepath.Base(path))
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve files on directory %s, %v", path, err)
		}
		files = append(files, dirFiles...)
	}

	return files, nil
}

// readDirectoryFilesRecursively retrieve directory and subdirectories files
func readDirectoryFilesRecursively(path string, parents ...string) (Files, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve directory data on %s, %v", path, err)
	}

	var files []File

	// we need to separate files and directories iteration over entries because of confusing
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		entryInfo, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to load info of %s, %v", entry.Name(), err)
		}

		files = append(files, File{
			Name:    entry.Name(),
			Size:    entryInfo.Size(),
			Path:    filepath.Join(path, entry.Name()),
			Parents: parents,
		})
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		parents = append(parents, entry.Name())
		innerFiles, err := readDirectoryFilesRecursively(filepath.Join(path, entry.Name()), parents...)
		if err != nil {
			return nil, err
		}
		files = append(files, innerFiles...)
	}

	return files, nil
}
