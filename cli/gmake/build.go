package main

import (
	"fmt"
	"log"
	"strings"
)

func getSupportedPlatforms() []string {
	return []string{
		"darwin/amd64",
		"linux/amd64",
		"linux/arm-7",
		"windows/amd64",
	}
}

func init() {
	log.Print("123")
	var packageApp = false
	var platform = ""
	commandDescription := `This command will check to ensure all pre-requistes are installed prior to building. If not, it will attempt to install them. Building comprises of a number of steps: install frontend dependencies, build frontend, pack frontend, compile main application.`
	initCmd := app.Command("build", "Builds your Wails project").
		LongDescription(commandDescription).BoolFlag("p", "Package application on successful build", &packageApp)
	var b strings.Builder
	for _, plat := range getSupportedPlatforms() {
		fmt.Fprintf(&b, " - %s\n", plat)
	}
	initCmd.StringFlag("x",
		fmt.Sprintf("Cross-compile application to specified platform via xgo\n%s", b.String()),
		&platform)

	initCmd.Action(func() error {
		log.Print("command build")
		log.Printf("packageApp %v ", packageApp)
		return nil
	})
}
