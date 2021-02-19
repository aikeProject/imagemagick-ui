package core

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"gopkg.in/gographics/imagick.v3/imagick"
)

type File struct {
	Data      string `json:"data"` // base64字符串
	Name      string `json:"name"`
	Size      int    `json:"size"`
	Extension string `json:"extension"` // 文件扩展名
	mw        *imagick.MagickWand
}

func NewFile(fileJson string) (*File, error) {
	file := &File{}
	if err := json.Unmarshal([]byte(fileJson), &file); err != nil {
		return file, err
	}
	return file, nil
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
	log.Println("destroy <= " + f.Name)
	f.mw.Destroy()
}
