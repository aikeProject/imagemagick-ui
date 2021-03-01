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
	AppMacOS          = "build/" + AppName + ".app/Contents/MacOS/" + AppName
	Frameworks        = "build/" + AppName + ".app/Contents/Frameworks"
)

type RunMac struct {
	Verbose    bool
	PackageApp bool
	BuildMode  string
}

func NewRunMac(op *Options) (*RunMac, error) {
	r := &RunMac{Verbose: op.Verbose, BuildMode: op.BuildMode, PackageApp: op.PackageApp}
	return r, r.init()
}

func (r *RunMac) init() error {
	if err := r.CompileTmpMagick(); err != nil {
		return err
	}
	if err := r.BuildImagick(); err != nil {
		return err
	}
	return nil
}

// 在"/tmp"目录下生成临时文件"ImageMagick-7.0.10"
// 用于运行"gopkg.in/gographics/imagick.v3/imagick"
func (r *RunMac) CompileTmpMagick() error {
	fsHelper := cmd.NewFSHelper()
	program := cmd.NewProgramHelper(true)
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)
	buildSpinner.Start("初始化...")
	if err := os.RemoveAll(ImageMagickTmp); err != nil {
		buildSpinner.Error(err.Error())
		return err
	}
	// 创建临时文件
	if err := fsHelper.MkDir(ImageMagickTmp); err != nil {
		buildSpinner.Error(err.Error())
		return err
	}
	dir, err := fsHelper.LocalDir(SourceDir)
	if err != nil {
		buildSpinner.Error(err.Error())
		return err
	}
	filenames, err := dir.GetAllFilenames()
	if err != nil {
		buildSpinner.Error(err.Error())
		return err
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
		return err
	}
	if err := program.RunCommandArray([]string{
		"install_name_tool",
		"-change",
		"/" + ImageMagick + "/lib/" + LibMagickCore,
		LibMagickCoreFile,
		LibMagickWandFile,
	}); err != nil {
		buildSpinner.Error(err.Error())
		return err
	}
	if err := program.RunCommandArray([]string{"install_name_tool", "-id", LibMagickCoreFile, LibMagickCoreFile}); err != nil {
		buildSpinner.Error(err.Error())
		return err
	}
	buildSpinner.Success(fmt.Sprintf("生成%s临时文件", ImageMagick))
	return nil
}

// 编译 imagick 包
// 自定义 CGO_CFLAGS CGO_LDFLAGS 编译 gopkg.in/gographics/imagick.v3/imagick
func (r *RunMac) BuildImagick() error {
	imagick := "go build -tags no_pkgconfig " + ImagickPackage
	if r.Verbose {
		imagick = "go build -v -x -tags no_pkgconfig " + ImagickPackage
	}
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)
	buildSpinner.Start("自定义 CGO_CFLAGS CGO_LDFLAGS 编译 gopkg.in/gographics/imagick.v3/imagick")
	program := cmd.NewProgramHelper(true)
	if err := os.Setenv("CGO_CFLAGS", CgoCflagsImagick); err != nil {
		buildSpinner.Error(err.Error())
		return err
	}
	if err := os.Setenv("CGO_LDFLAGS", CgoLdflagsImagick); err != nil {
		buildSpinner.Error(err.Error())
		return err
	}
	// 先编译imagick包
	if err := program.RunCommand(imagick); err != nil {
		buildSpinner.Error(err.Error())
		return err
	}
	buildSpinner.Success("imagick编译完成")
	return nil
}

// 打包
func (r *RunMac) Build() error {
	fsHelper := cmd.NewFSHelper()
	program := cmd.NewProgramHelper(true)
	wails := "wails build"
	// mac下打包成.app的包
	if r.PackageApp {
		wails += " -p"
	} else if r.BuildMode == cmd.BuildModeDebug {
		// debug模式
		wails += " -d"
	}
	if r.Verbose {
		// 输出详细信息
		wails += " -verbose"
	}

	if err := program.RunCommand(wails); err != nil {
		return err
	}

	// 安装包
	if r.PackageApp {
		if err := fsHelper.MkDirs(Frameworks); err != nil {
			return err
		}
		copyFile := []string{LibMagickWand, LibMagickCore}
		for _, f := range copyFile {
			if err := fsHelper.CopyFile("/tmp/"+ImageMagick+"/lib/"+f, Frameworks+"/"+f); err != nil {
				return err
			}
		}
		// 修改"ImageMagick"三方库动态链接地址
		// 从安装包内部加载外部模块
		// 添加"@rpath"地址
		commands := [][]string{
			{
				"install_name_tool",
				"-add_rpath",
				"@loader_path/../Frameworks",
				AppMacOS,
			},
			{
				"install_name_tool",
				"-change",
				LibMagickWandFile,
				"@rpath/" + LibMagickWand,
				AppMacOS,
			},
			{
				"install_name_tool",
				"-change",
				LibMagickCoreFile,
				"@rpath/" + LibMagickCore,
				AppMacOS,
			},
			{
				"install_name_tool",
				"-id",
				"@rpath/" + LibMagickWand,
				Frameworks + "/" + LibMagickWand,
			},
			{
				"install_name_tool",
				"-change",
				LibMagickCoreFile,
				"@rpath/" + LibMagickCore,
				Frameworks + "/" + LibMagickWand,
			},
			{
				"install_name_tool",
				"-id",
				"@rpath/" + LibMagickCore,
				Frameworks + "/" + LibMagickCore,
			},
		}

		for _, command := range commands {
			if err := program.RunCommandArray(command); err != nil {
				return err
			}
		}
	}

	return nil
}

// 启动本地开发服务
func (r *RunMac) Serve() error {
	wails := "wails serve"
	program := cmd.NewProgramHelper(true)
	if err := program.RunCommand(wails); err != nil {
		return err
	}
	return nil
}
