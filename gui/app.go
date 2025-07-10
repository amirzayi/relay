package gui

import (
	"context"
	_ "embed"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const version = "v1.0.0"

const appID = "923803ee-54c2-491a-83e8-baf89b97ba06"

var appCtx context.Context

// App struct
type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) id() string {
	return appID
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	appCtx = ctx
}

func (a *App) Version() string {
	return version
}

func (a *App) onSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	secondInstanceArgs := secondInstanceData.Args
	runtime.WindowUnminimise(appCtx)
	runtime.Show(appCtx)
	runtime.MessageDialog(appCtx, runtime.MessageDialogOptions{
		Title:   "Second Instance",
		Message: "A second instance of the application was launched",
		Type:    runtime.ErrorDialog,
	})
	go runtime.EventsEmit(appCtx, "launchArgs", secondInstanceArgs)
}
