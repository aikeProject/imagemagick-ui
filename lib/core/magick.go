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
func (m *Magick) Resize(w, h uint) error {
	width := m.GetImageWidth()
	height := m.GetImageHeight()
	// 保持图像纵横比
	if width > height {
		height = uint((float32(w) / float32(width)) * float32(height))
		width = w
	} else if width < height {
		width = uint((float32(h) / float32(height)) * float32(width))
		height = h
	}
	return m.AdaptiveResizeImage(width, height)
}
