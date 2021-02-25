package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"imagemagick-ui/lib"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/thoas/go-funk"

	"github.com/wailsapp/wails"
)

var fSHelper = lib.NewFSHelper()
var covertExtFiles = []string{".gif", ".pdf"}

type File struct {
	Id      string `json:"id"`
	Data    string `json:"data"` // base64字符串
	Name    string `json:"name"`
	Size    int    `json:"size"`
	Ext     string `json:"ext"`    // 文件扩展名
	Status  int    `json:"status"` // 文件状态
	mw      *Magick
	logger  *wails.CustomLogger
	runtime *wails.Runtime
	conf    *lib.Config
}

func NewFile(fileJson string, config *lib.Config) (*File, error) {
	file := &File{
		mw:   NewMagick(),
		conf: config,
	}
	if err := json.Unmarshal([]byte(fileJson), &file); err != nil {
		return file, err
	}
	file.Ext = filepath.Ext(file.Name)
	return file, nil
}

func (f *File) WailsInit(runtime *wails.Runtime) {
	f.runtime = runtime
	f.logger = f.runtime.Log.New("File")
	f.logger.Infof("File init %s", f.Name)
}

// 解析base64字符串
func (f *File) Decode() ([]byte, error) {
	return base64.StdEncoding.DecodeString(f.Data)
}

// 文件处理完毕之后将其写入至本地文件
func (f *File) Write() error {
	// 文件处理开始
	f.Status = Running

	bytes, err := f.Decode()
	if err != nil {
		return err
	}
	if err := f.mw.ReadImageBlob(bytes); err != nil {
		return err
	}
	// 调整图像大小
	width, height := f.mw.Resize(f.conf.App.Width, f.conf.App.Height)
	// 生成文件写入路径
	p, err := f.filepath(width, height)
	if err != nil {
		return err
	}
	if funk.ContainsString(covertExtFiles, strings.ToLower(f.Ext)) {
		// 输出多张图片，比如：gif转换为png
		// eg: xxx.gif => xxx-0.png xxx-1.png ...
		if err := f.mw.WriteImages(p, false); err != nil {
			return err
		}
	} else {
		if err := f.mw.WriteImage(p); err != nil {
			return err
		}
	}
	// 文件处理结束
	f.Status = Done
	return nil
}

// 返回文件保存路径
func (f *File) filepath(width, height uint) (string, error) {
	p := path.Join(f.conf.App.OutDir, f.baseName())
	fp := path.Join(p, f.renameWidthHeight(width, height))
	// 检查是否存在该目录
	if !fSHelper.DirExists(p) {
		if err := fSHelper.MkDirs(p); err != nil {
			return "", err
		}
		return fp, nil
	}
	return fp, nil
}

// 去除文件名的扩展名
// eg: xxx.png => xxx
func (f *File) baseName() string {
	if f.Ext != "" {
		name := strings.Replace(f.Name, f.Ext, "", 1)
		return name
	}
	return f.Name
}

// 重命名
// eg: xxx.png => xxx.jpg
func (f *File) rename() string {
	return fmt.Sprintf("%s.%s", f.baseName(), f.conf.App.Target)
}

// 根据图片width和height重命名文件
// eg: xxx.png => xxx-200x200.png
func (f *File) renameWidthHeight(width, height uint) string {
	if width > 0 && height > 0 {
		name := f.baseName()
		return fmt.Sprintf("%s-%dx%d.%s", name, width, height, f.conf.App.Target)
	}
	return f.rename()
}

// 设置"Magick"
func (f *File) SetMagick(m *Magick) {
	f.mw = m
}

func (f *File) Display() error {
	if err := f.mw.DisplayImage(os.Getenv("DISPLAY")); err != nil {
		return err
	}
	return nil
}

func (f *File) Destroy() {
	f.logger.Infof("file destroy <= %s", f.Name)
	f.mw.Destroy()
}
