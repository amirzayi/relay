package fileutil

import (
	"io"
	"os"
	"path/filepath"

	"github.com/amirzayi/relay/config"
	"github.com/amirzayi/relay/pkg/progressive"
)

// TODO: make save path configurable
func (f File) ProgressiveWrite(r io.Reader, onWrite func(transferred int64, percent int)) (int64, error) {
	filePath := filepath.Join(f.Parents...)
	filePath = filepath.Join(config.DefaultDirectory(), filePath, f.Name)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return 0, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	pw := progressive.NewWriter(file, onWrite, f.Size)
	return io.CopyN(pw, r, f.Size)
}

func (f File) ProgressiveRead(w io.Writer, onWrite func(transferred int64, percent int)) (int64, error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	pw := progressive.NewWriter(w, onWrite, f.Size)
	return io.Copy(pw, file)
}
