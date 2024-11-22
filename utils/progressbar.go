package utils

import (
	"fmt"
	"strings"
)

func DrawProgressBar(percent, width int, fileName string) {

	// draw transferred data progress bar
	progressBar := strings.Repeat("", 100-percent/(100/width))
	progressBar += strings.Repeat("█", percent/(100/width))

	fmt.Printf("\r%d%% [%-*s] %s", percent, width, progressBar, fileName)
}
