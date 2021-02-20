package main

import (
	"imagemagick-ui/lib/core"
	"log"
	"os"
	"path"

	"gopkg.in/gographics/imagick.v3/imagick"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	outDir := path.Join(homeDir, "Desktop", "imagemagick-ui")
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.Mkdir(outDir, 0777); err != nil {
			return
		}
		return
	}
}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()
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
