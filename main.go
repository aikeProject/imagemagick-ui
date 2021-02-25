package main

import (
	"imagemagick-ui/lib"
	"imagemagick-ui/lib/core"
	"log"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()
	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")
	app := wails.CreateApp(&wails.AppConfig{
		Width:  824,
		Height: 568,
		Title:  "imagemagick-ui",
		JS:     js,
		CSS:    css,
		Colour: "rgba(255,255,255,1)",
	})
	newConfig := lib.NewConfig()
	manager := core.NewManager(newConfig)
	app.Bind(manager)
	app.Bind(newConfig)
	if err := app.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
