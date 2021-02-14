package main

import (
	"fmt"
	"runtime"

	"github.com/wailsapp/wails/cmd"
)

type Options struct {
	Verbose    bool
	PackageApp bool
	BuildMode  string
}

func init() {
	var platform = runtime.GOOS
	var verbose = false
	var debugMode = false
	var packageApp = false
	initCmd := app.Command("build", "打包项目").
		LongDescription("在打包之前提前处理好ImageMagick包相关配置").
		BoolFlag("d", "启用debug模式", &debugMode).
		BoolFlag("p", "打包成应用程序", &packageApp).
		BoolFlag("verbose", "打印详细日志", &verbose)

	// Build application
	buildMode := cmd.BuildModeProd
	if debugMode {
		buildMode = cmd.BuildModeDebug
	}

	initCmd.Action(func() error {
		op := &Options{
			Verbose:    verbose,
			PackageApp: packageApp,
			BuildMode:  buildMode,
		}
		switch platform {
		case "darwin":
			runMac, err := NewRunMac(op)
			if err != nil {
				return err
			}
			runMac.Build()
			return nil
		case "windows":
			return nil
		case "linux":
			return nil
		default:
			return fmt.Errorf("platform '%s' not supported for bundling yet", platform)
		}
	})
}
