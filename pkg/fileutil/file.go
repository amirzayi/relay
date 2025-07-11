package fileutil

import (
	"encoding/gob"
	"fmt"
	"io"
)

type File struct {
	Name    string
	Size    int64
	Path    string
	Parents []string
}

func (f File) ShortedName(length, start, end int) string {
	s := f.Name
	if len(s) > length {
		s = fmt.Sprintf("%s...%s", s[:start], s[len(s)-end:])
	}
	return s
}

type Files []File

func (fs Files) SendDetails(w io.Writer) error {
	return gob.NewEncoder(w).Encode(fs)
}

func GetDetails(r io.Reader) (Files, error) {
	var files Files
	err := gob.NewDecoder(r).Decode(&files)
	return files, err
}
