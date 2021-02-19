package core

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/wailsapp/wails"

	"gopkg.in/gographics/imagick.v3/imagick"
)

// 文件状态常量
const (
	NotStarted = iota
	Start
	Running
	Done
)

type File struct {
	Data      string `json:"data"` // base64字符串
	Name      string `json:"name"`
	Size      int    `json:"size"`
	Extension string `json:"extension"` // 文件扩展名
	Status    int    `json:"status"`    // 文件状态
	mw        *imagick.MagickWand
	runtime   *wails.Runtime
	logger    *wails.CustomLogger
}

func NewFile(fileJson string) (*File, error) {
	file := &File{}
	if err := json.Unmarshal([]byte(fileJson), &file); err != nil {
		return file, err
	}
	return file, nil
}

func (f *File) WailsInit(runtime *wails.Runtime) {
	f.runtime = runtime
	f.logger = f.runtime.Log.New("File")
	f.logger.Info("File Item initialized...")
}

// 解析base64字符串
func (f *File) Decode() error {
	f.SetMagic()
	bytes, err := base64.StdEncoding.DecodeString(f.Data)
	if err != nil {
		return err
	}
	// 读取文件内容
	return f.mw.ReadImageBlob(bytes)
}

// 文件处理完毕之后将其写入至本地文件
func (f *File) Write() error {
	defer f.Destroy()
	if err := f.Decode(); err != nil {
		return err
	}
	width := f.mw.GetImageWidth()
	height := f.mw.GetImageHeight()
	if err := f.mw.ResizeImage(uint(width/2), uint(height/2), imagick.FILTER_LANCZOS); err != nil {
		return err
	}
	//if err := f.mw.WriteImage("resize-file.png"); err != nil {
	//	return err
	//}
	if err := f.Display(); err != nil {
		return err
	}
	return nil
}

func (f *File) SetMagic() {
	f.mw = imagick.NewMagickWand()
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
