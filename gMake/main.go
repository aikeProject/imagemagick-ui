package main

import (
	"os"
	"os/exec"

	"github.com/wailsapp/wails/cmd"
)

const (
	AppName        = "imagemagick-ui"
	ImagickPackage = "gopkg.in/gographics/imagick.v3/imagick"
)

type Options struct {
	Verbose    bool
	PackageApp bool
	BuildMode  string
}

// Create Logger
var logger = cmd.NewLogger()

// Create main app
var app = cmd.NewCli("GMake", "Makefile的Go版本")

// Main!
func main() {
	err := app.Run()
	if err != nil {
		logger.Error(err.Error())
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
