package utils

import (
	"fmt"
	"io"
)

// DrawRWProgressbar reads data from an io.Reader and writes it to an io.Writer with buffering while displaying a progress bar.
func DrawRWProgressbar(r io.Reader, w io.Writer, size int64, bufferSize, barSize int, barTitle string) error {
	lr := io.LimitReader(r, size)
	rw := io.TeeReader(lr, w)

	buffer := make([]byte, bufferSize)
	syncedBytes := int64(0)
	for syncedBytes < size {
		n, err := rw.Read(buffer)
		if err != nil {
			return fmt.Errorf("failed to sync buffer between network and disk: %v", err)
		}

		syncedBytes += int64(n)
		donePercent := int(syncedBytes * 100 / size)
		drawProgressBar(donePercent, barSize, barTitle)
	}
	return nil
}

// WriteFromReader reads data from an io.Reader and writes it to an io.Writer with buffering.
func WriteFromReader(r io.Reader, w io.Writer, size int64, bufferSize int) error {
	lr := io.LimitReader(r, size)
	rw := io.TeeReader(lr, w)

	buffer := make([]byte, bufferSize)
	totalBytesRead := int64(0)
	for totalBytesRead < size {
		_, err := rw.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to sync buffer between network and disk: %v", err)
		}
	}
	return nil
}
