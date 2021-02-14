package main

import (
	"fmt"
	"runtime"
)

func init() {
	var platform = runtime.GOOS
	initCmd := app.Command("build", "Builds your Wails project")

	initCmd.Action(func() error {
		switch platform {
		case "darwin":
			PackageMac()
			return nil
		case "windows":
			PackageWin()
			return nil
		case "linux":
			return nil
		default:
			return fmt.Errorf("platform '%s' not supported for bundling yet", platform)
		}
	})
}
