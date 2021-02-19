package core

import (
	"gopkg.in/gographics/imagick.v3/imagick"
)

type Manager struct {
	files []*File
}

func NewManager() *Manager {
	return &Manager{}
}

// 将文件添加至"Manager"
func (m *Manager) HandleFile(fileJson string) error {
	file, err := NewFile(fileJson)
	if err != nil {
		return err
	}
	if m.files != nil {
		m.files = append(m.files, file)
	}
	m.files = []*File{file}
	return nil
}

// 并发处理文件
func (m *Manager) Convert() (errs []error) {
	for _, file := range m.files {
		if err := file.Write(); err != nil {
			return append(errs, err)
		}
	}
	return errs
}

// 垃圾回收
func (m *Manager) Clear() {
	defer imagick.Terminate()
	m.files = []*File{}
}
