package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

const (
	ImageMagick       = "ImageMagick-7.0.10"
	ImageMagickTmp    = "/tmp/" + ImageMagick      // 临时文件
	SourceDir         = "../source/" + ImageMagick // 源文件目录
	LibMagickWand     = "libMagickWand-7.Q16HDRI.8.dylib"
	LibMagickCore     = "libMagickCore-7.Q16HDRI.8.dylib"
	LibMagickWandFile = "/tmp/" + ImageMagick + "/lib/" + LibMagickWand
	LibMagickCoreFile = "/tmp/" + ImageMagick + "/lib/" + LibMagickCore
	CgoCflagsImagick  = "-Xpreprocessor -fopenmp -DMAGICKCORE_HDRI_ENABLE=1 -DMAGICKCORE_QUANTUM_DEPTH=16 -I" + ImageMagickTmp + "/include/ImageMagick-7"
	CgoLdflagsImagick = "-g -O2 -L" + ImageMagickTmp + "/lib/ -lMagickWand-7.Q16HDRI.8 -lMagickCore-7.Q16HDRI.8"
)

type RunMac struct {
	Verbose    bool
	PackageApp bool
	BuildMode  string
}

func NewRunMac(op *Options) *RunMac {
	return &RunMac{Verbose: op.Verbose, BuildMode: op.BuildMode, PackageApp: op.PackageApp}
}

func (r *RunMac) init() (*RunMac, error) {
	fsHelper := cmd.NewFSHelper()
	program := cmd.NewProgramHelper(true)
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)
	buildSpinner.Start("初始化...")
	if err := os.RemoveAll(ImageMagickTmp); err != nil {
		buildSpinner.Error(err.Error())
		return nil, err
	}
	// 创建临时文件
	if err := fsHelper.MkDir(ImageMagickTmp); err != nil {
		buildSpinner.Error(err.Error())
		return nil, err
	}
	dir, err := fsHelper.LocalDir(SourceDir)
	if err != nil {
		buildSpinner.Error(err.Error())
		return nil, err
	}
	filenames, err := dir.GetAllFilenames()
	if err != nil {
		buildSpinner.Error(err.Error())
		return nil, err
	}
	filenames.Each(func(s string) {
		r := strings.Replace(s, os.Getenv("PWD")+"/source", "/tmp", 1)
		if fsHelper.DirExists(s) {
			if err := fsHelper.MkDir(r); err != nil {
				buildSpinner.Error(err.Error())
			}
		} else {
			if err := fsHelper.CopyFile(s, r); err != nil {
				buildSpinner.Error(err.Error())
			}
		}
	})
	if err := program.RunCommandArray([]string{"install_name_tool", "-id", LibMagickWandFile, LibMagickWandFile}); err != nil {
		buildSpinner.Error(err.Error())
		return nil, err
	}
	if err := program.RunCommandArray([]string{
		"install_name_tool",
		"-change",
		"/" + ImageMagick + "/lib/" + LibMagickCore,
		LibMagickCoreFile,
		LibMagickWandFile,
	}); err != nil {
		buildSpinner.Error(err.Error())
		return nil, err
	}
	if err := program.RunCommandArray([]string{"install_name_tool", "-id", LibMagickCoreFile, LibMagickCoreFile}); err != nil {
		buildSpinner.Error(err.Error())
		return nil, err
	}
	buildSpinner.Success(fmt.Sprintf("生成%s临时文件", ImageMagick))
	return r, nil
}

func (r RunMac) Build() {
	imagick := "go build -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick"
	wails := "wails build"
	// debug模式
	if r.BuildMode == cmd.BuildModeDebug {
		wails += " -d"
	}
	// mac下打包成.app的包
	if r.PackageApp {
		wails += " -p"
	}
	if r.Verbose {
		// 输出详细信息
		imagick = "go build -v -x -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick"
		wails += " -verbose"
	}
	logger := cmd.NewLogger()
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)
	buildSpinner.Start("自定义 CGO_CFLAGS CGO_LDFLAGS 编译 gopkg.in/gographics/imagick.v3/imagick")
	program := cmd.NewProgramHelper(true)
	if err := os.Setenv("CGO_CFLAGS", CgoCflagsImagick); err != nil {
		buildSpinner.Error(err.Error())
		return
	}
	if err := os.Setenv("CGO_LDFLAGS", CgoLdflagsImagick); err != nil {
		buildSpinner.Error(err.Error())
		return
	}
	// 先编译imagick包
	if err := program.RunCommand(imagick); err != nil {
		buildSpinner.Error(err.Error())
		return
	}
	buildSpinner.Success("imagick编译完成")
	if err := program.RunCommand(wails); err != nil {
		logger.Error(err.Error())
		return
	}
}
