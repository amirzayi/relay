package utils

import "fmt"

var sizes = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}

const base = 1024

func ConvertByteSizeToHumanReadable(size float64) string {
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
