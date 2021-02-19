package core

import "gopkg.in/gographics/imagick.v3/imagick"

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
	m.files = append(m.files, file)
	return nil
}

// 并发处理文件
func (m *Manager) Convert() (errs []error) {
	imagick.Initialize()
	m.files[0].SetMagic()
	if err := m.files[0].Write(); err != nil {
		return append(errs, err)
	}
	defer imagick.Terminate()
	return errs
}
