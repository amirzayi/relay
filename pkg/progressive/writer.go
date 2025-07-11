package progressive

import "io"

type progressiveWrite struct {
	w                io.Writer
	onWrite          func(transferred int64, percent int)
	totalSize        int64
	transferredBytes int64
}

func NewWriter(w io.Writer, onWrite func(transferred int64, percent int), totalSize int64) io.Writer {
	return &progressiveWrite{w, onWrite, totalSize, 0}
}

func (pw *progressiveWrite) Write(p []byte) (n int, err error) {
	n, err = pw.w.Write(p)
	pw.transferredBytes += int64(n)
	pw.onWrite(pw.transferredBytes, int((pw.transferredBytes*100)/pw.totalSize))
	return n, err
}
