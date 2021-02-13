package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

const (
	ImageMagick       = "ImageMagick-7.0.10"
	ImageMagickTmp    = "/tmp/" + ImageMagick         // 临时文件
	SourceDir         = "../../source/" + ImageMagick // 源文件目录
	LibMagickWand     = "libMagickWand-7.Q16HDRI.8.dylib"
	LibMagickCore     = "libMagickCore-7.Q16HDRI.8.dylib"
	LibMagickWandFile = "/tmp/" + ImageMagick + "/lib/" + LibMagickWand
	LibMagickCoreFile = "/tmp/" + ImageMagick + "/lib/" + LibMagickCore
	CgoCflagsImagick  = "-Xpreprocessor -fopenmp -DMAGICKCORE_HDRI_ENABLE=1 -DMAGICKCORE_QUANTUM_DEPTH=16 -I" + ImageMagickTmp + "/include/ImageMagick-7"
	CgoLdflagsImagick = "-g -O2 -L" + ImageMagickTmp + "/lib/ -lMagickWand-7.Q16HDRI.8 -lMagickCore-7.Q16HDRI.8"
)

func init() {
	fsHelper := cmd.NewFSHelper()
	program := cmd.NewProgramHelper(true)
	buildSpinner := spinner.NewSpinner()
	buildSpinner.Start("初始化...")
	if err := os.RemoveAll(ImageMagickTmp); err != nil {
		buildSpinner.Error(err.Error())
		return
	}
	// 创建临时文件
	if err := fsHelper.MkDir(ImageMagickTmp); err != nil {
		buildSpinner.Error(err.Error())
		return
	}
	dir, err := fsHelper.LocalDir(SourceDir)
	if err != nil {
		buildSpinner.Error(err.Error())
		return
	}
	filenames, err := dir.GetAllFilenames()
	if err != nil {
		buildSpinner.Error(err.Error())
		return
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
		return
	}
	if err := program.RunCommandArray([]string{
		"install_name_tool",
		"-change",
		"/" + ImageMagick + "/lib/" + LibMagickCore,
		LibMagickCoreFile,
		LibMagickWandFile,
	}); err != nil {
		buildSpinner.Error(err.Error())
		return
	}
	if err := program.RunCommandArray([]string{"install_name_tool", "-id", LibMagickCoreFile, LibMagickCoreFile}); err != nil {
		buildSpinner.Error(err.Error())
		return
	}
	buildSpinner.Success(fmt.Sprintf("生成%s临时文件", ImageMagick))
}

func f() {
	logger := cmd.NewLogger()
	program := cmd.NewProgramHelper(true)
	if err := os.Setenv("CGO_CFLAGS", CgoCflagsImagick); err != nil {
		logger.Error(err.Error())
		return
	}
	if err := os.Setenv("CGO_LDFLAGS", CgoLdflagsImagick); err != nil {
		logger.Error(err.Error())
		return
	}
	// 先编译imagick包
	if err := program.RunCommand("go build -v -x -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick"); err != nil {
		logger.Error(err.Error())
		return
	}
	if err := program.RunCommand("wails build -d"); err != nil {
		logger.Error(err.Error())
		return
	}
}
