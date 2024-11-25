package utils

import (
	"fmt"
	"io"
)

// DrawRWProgressbar reads data from an io.Reader and writes it to an io.Writer with buffering while displaying a progress bar.
func DrawRWProgressbar(r io.Reader, w io.Writer, size int64, bufferSize, barSize int, barTitle string) error {
	return drawRWProgressbar(r, w, size, bufferSize, func(syncedBytes int64) {
		donePercent := int(syncedBytes * 100 / size)
		drawProgressBar(donePercent, barSize, barTitle)
	})
}

// WriteFromReader reads data from an io.Reader and writes it to an io.Writer with buffering.
func WriteFromReader(r io.Reader, w io.Writer, size int64, bufferSize int) error {
	return drawRWProgressbar(r, w, size, bufferSize, nil)
}

func drawRWProgressbar(r io.Reader, w io.Writer, size int64, bufferSize int, progressFn func(syncedBytes int64)) error {
	lr := io.LimitReader(r, size)
	rw := io.TeeReader(lr, w)

	buffer := make([]byte, bufferSize)
	syncedBytes := int64(0)
	for syncedBytes < size {
		n, err := rw.Read(buffer)
		if err != nil {
			return fmt.Errorf("\r\nfailed to sync buffer between network and disk: %v", err)
		}
		syncedBytes += int64(n)
		if progressFn != nil {
			progressFn(syncedBytes)
		}
	}
	return nil
}
