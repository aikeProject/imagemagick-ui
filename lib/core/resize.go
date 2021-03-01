package core

type Resize struct {
	Width  uint
	Height uint
}

// 保持图像纵横比比例
func (r *Resize) Base(w, h uint) (rw, rh uint) {
	rw = r.Width
	rh = r.Height

	switch {
	case w > 0 && h == 0:
		return w, r.rHeight(w)
	case w == 0 && h > 0:
		return r.rWidth(h), h
	case w > 0 && h > 0:
		if r.Width >= r.Height {
			rh = r.rHeight(w)
			rw = w
		} else if r.Width <= r.Height {
			rw = r.rWidth(h)
			rh = h
		}
		return rw, rh
	}

	return rw, rh
}

func (r *Resize) rWidth(h uint) uint {
	return uint((float32(h) / float32(r.Height)) * float32(r.Width))
}

func (r Resize) rHeight(w uint) uint {
	return uint((float32(w) / float32(r.Width)) * float32(r.Height))
}
