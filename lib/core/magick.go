package core

import (
	"gopkg.in/gographics/imagick.v3/imagick"
)

type Magick struct {
	*imagick.MagickWand
}

func NewMagick() *Magick {
	return &Magick{
		MagickWand: imagick.NewMagickWand(),
	}
}

// 调整图像大小
func (m *Magick) Resize(w, h uint) (rw, rh uint) {
	width := m.GetImageWidth()
	height := m.GetImageHeight()
	resize := Resize{
		Width:  width,
		Height: height,
	}
	rw, rh = resize.Base(w, h)
	if m.AdaptiveResizeImage(rw, rh) != nil {
		return width, height
	}
	return rw, rh
}
