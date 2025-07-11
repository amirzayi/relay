package gui

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var (
	//go:embed all:frontend/dist
	assets embed.FS

	//go:embed frontend/src/logo.png
	icon []byte
)

func Start() {
	app := NewApp()
	host := NewHost()
	client := NewSender()

	err := wails.Run(&options.App{
		Title:      "Relay File Sharing",
		Width:      400,
		Height:     500,
		MinWidth:   300,
		MinHeight:  400,
		Fullscreen: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 7, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
			host,
			client,
		},
		Windows: &windows.Options{
			DisablePinchZoom: true,
		},
		Linux: &linux.Options{
			Icon:        icon,
			ProgramName: "Relay File Sharing",
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               app.id(),
			OnSecondInstanceLaunch: app.onSecondInstanceLaunch,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
