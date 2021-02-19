package main

import (
	"imagemagick-ui/lib/core"
	"log"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func main() {

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "imagemagick-ui",
		JS:     js,
		CSS:    css,
		Colour: "rgba(255,255,255,1)",
	})
	manager := core.NewManager()
	app.Bind(manager)
	if err := app.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
