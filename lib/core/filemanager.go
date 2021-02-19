package core

import (
	"github.com/wailsapp/wails"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type Manager struct {
	files   []*File
	runtime *wails.Runtime
	logger  *wails.CustomLogger
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) WailsInit(runtime *wails.Runtime) error {
	m.runtime = runtime
	m.logger = m.runtime.Log.New("FileManager")
	m.logger.Info("File Manager initialized...")
	return nil
}

// 将文件添加至"Manager"
func (m *Manager) HandleFile(fileJson string) error {
	file, err := NewFile(fileJson)
	if err != nil {
		return err
	}
	file.WailsInit(m.runtime)
	m.files = append(m.files, file)
	m.logger.Infof("添加至Manager <= %s", file.Name)
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
