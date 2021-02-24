package core

import (
	"encoding/base64"
	"encoding/json"
	"imagemagick-ui/lib"
	"os"
	"path"

	"github.com/wailsapp/wails"
)

type File struct {
	Id        string `json:"id"`
	Data      string `json:"data"` // base64字符串
	Name      string `json:"name"`
	Size      int    `json:"size"`
	Extension string `json:"extension"` // 文件扩展名
	Status    int    `json:"status"`    // 文件状态
	mw        *Magick
	logger    *wails.CustomLogger
	runtime   *wails.Runtime
	conf      *lib.Config
}

func NewFile(fileJson string, config *lib.Config) (*File, error) {
	file := &File{
		mw:   NewMagick(),
		conf: config,
	}
	if err := json.Unmarshal([]byte(fileJson), &file); err != nil {
		return file, err
	}
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
	if err := f.mw.Resize(f.conf.App.Width, f.conf.App.Height); err != nil {
		return err
	}
	if err := f.mw.WriteImage(path.Join(f.conf.App.OutDir, f.Name)); err != nil {
		return err
	}

	// 文件处理结束
	f.Status = Done
	return nil
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
