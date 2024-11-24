package utils

import (
	"fmt"
	"strings"
)

// drawProgressBar draw percentage progress bar
func drawProgressBar(percent, width int, barName string) {
	progress := strings.Repeat("", 100-percent/(100/width))
	progress += strings.Repeat("█", percent/(100/width))

	fmt.Printf("\r%d%% |%-*s| %s", percent, width, progress, barName)
}
