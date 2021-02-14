package main

import (
	"fmt"
	"runtime"
)

func init() {
	var platform = runtime.GOOS
	var verbose = false
	initCmd := app.Command("serve", "项目本地运行").
		LongDescription("在本地运行项目，便于开发调试").
		BoolFlag("verbose", "打印详细日志", &verbose)

	initCmd.Action(func() error {
		op := &Options{
			Verbose: verbose,
		}
		switch platform {
		case "darwin":
			runMac, err := NewRunMac(op)
			if err != nil {
				return err
			}
			if err := runMac.Serve(); err != nil {
				return err
			}
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
