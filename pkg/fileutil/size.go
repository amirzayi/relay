package fileutil

import "fmt"

func (fs Files) TotalSize() (totalSize int64) {
	for _, f := range fs {
		totalSize += f.Size
	}
	return
}

func (f File) HumanReadableSize() string {
	return convertSizeToHumanReadable(float64(f.Size))
}

func (fs Files) HumanReadableTotalSize() string {
	totalSize := fs.TotalSize()
	return convertSizeToHumanReadable(float64(totalSize))
}

var sizes = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

const base = 1024

func convertSizeToHumanReadable(size float64) string {
	unitsLimit := len(sizes)
	i := 0

	for size >= base && i < unitsLimit {
		size = size / base
		i++
	}
	f := "%.0f %s"
	if i > 1 {
		f = "%.2f %s"
	}
	return fmt.Sprintf(f, size, sizes[i])
}
