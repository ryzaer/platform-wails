package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func runApp() error {

	app := NewApp()

	return wails.Run(&options.App{
		Title:            "Elsana Launcher",
		Width:            400,
		Height:           450,
		WindowStartState: options.Normal,
		Frameless:        true,
		DisableResize:    true,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		BackgroundColour: &options.RGBA{
			R: 27,
			G: 38,
			B: 54,
			A: 1,
		},

		OnStartup: app.startup,

		Bind: []interface{}{
			app,
		},

		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "com.elsana.integrator",

			OnSecondInstanceLaunch: app.onSecondInstanceLaunch,
		},

		Windows: &windows.Options{
			ZoomFactor:           1.0,
			IsZoomControlEnabled: false,
			DisablePinchZoom:     true,
		},
	})
}

// Normal EXE entry point.
func main() {

	if err := runApp(); err != nil {
		fmt.Println("Error:", err)
	}
}

//export StartIntegrator
func StartIntegrator() {

	if err := runApp(); err != nil {
		fmt.Println("Error:", err)
	}
}
