package core

import (
	"encoding/base64"
	"encoding/json"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func init() {
	imagick.Initialize()
}

type File struct {
	Data      string `json:"data"` // base64字符串
	Name      string `json:"name"`
	Size      string `json:"size"`
	Extension string `json:"extension"` // 文件扩展名
	mw        *imagick.MagickWand
}

func NewFile(fileJson string) (*File, error) {
	file := &File{mw: imagick.NewMagickWand()}
	if err := json.Unmarshal([]byte(fileJson), &file); err != nil {
		return file, err
	}
	return file, nil
}

// 解析base64字符串
func (f *File) decode() error {
	if bytes, err := base64.StdEncoding.DecodeString(f.Data); err != nil {
		// 读取文件内容
		return f.mw.ReadImageBlob(bytes)
	}
	return nil
}

// 文件处理完毕之后将其写入至本地文件
func (f *File) Write() error {
	if err := f.decode(); err != nil {
		return err
	}
	if err := f.mw.ResizeImage(uint(200), uint(200), imagick.FILTER_LANCZOS); err != nil {
		return err
	}
	if err := f.mw.WriteImage("resize-file.png"); err != nil {
		return err
	}
	return nil
}
