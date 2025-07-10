package main

import (
	"embed"
	"os"

	"github.com/amirzayi/relay/cmd"
	"github.com/amirzayi/relay/gui"
)

var (
	//go:embed all:gui/frontend/dist
	assets embed.FS

	//go:embed gui/frontend/src/logo.png
	icon []byte
)

func main() {
	if len(os.Args) > 1 {
		cmd.Execute()
	} else {
		gui.Start(assets, icon)
	}
}
