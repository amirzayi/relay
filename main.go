package main

import (
	"os"

	"github.com/amirzayi/relay/cmd"
	"github.com/amirzayi/relay/gui"
)

func main() {
	if len(os.Args) > 1 {
		cmd.Execute()
	} else {
		gui.Start()
	}
}
