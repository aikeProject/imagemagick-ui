package core

import (
	"encoding/json"
	"imagemagick-ui/lib"
	"runtime/debug"
	"sync"

	"gopkg.in/gographics/imagick.v3/imagick"

	"github.com/wailsapp/wails"
)

// 文件状态常量
const (
	NotStarted  = iota // 初始状态
	Start              // 文件数据发送中
	SendSuccess        // 文件数据已发送到golang程序中
	Running            // 文件处理中
	Done               // 处理完毕
)

type Manager struct {
	files   []*File
	mw      *Magick
	conf    *lib.Config
	runtime *wails.Runtime
	logger  *wails.CustomLogger
}

func NewManager(config *lib.Config) *Manager {
	return &Manager{
		conf: config,
	}
}

func (m *Manager) WailsInit(runtime *wails.Runtime) error {
	m.runtime = runtime
	m.logger = m.runtime.Log.New("FileManager")
	m.logger.Info("File Manager initialized...")
	return nil
}

// 将文件添加至"Manager"
func (m *Manager) HandleFile(fileJson string) error {
	file, err := NewFile(fileJson, m.conf)
	if err != nil {
		return err
	}
	file.WailsInit(m.runtime)
	m.files = append(m.files, file)
	m.logger.Infof("文件添加至 Manager <= %s", file.Name)
	return nil
}

// 并发处理文件
func (m *Manager) Convert(idStr string) (err error) {
	var wg sync.WaitGroup
	var files []*File
	var ids []string

	err = json.Unmarshal([]byte(idStr), &ids)
	if err != nil {
		return err
	}
	m.logger.Infof("ids %v", ids)
	if len(ids) > 0 {
		for _, id := range ids {
			if file := m.getByIdFile(id); file != nil {
				files = append(files, file)
			}
		}
	} else {
		files = m.files
	}

	// 将文件状态重置回未处理状态
	for _, file := range files {
		file.Status = SendSuccess
	}

	m.SetMagick()
	wg.Add(len(files))
	for _, f := range files {
		go func(file *File, w *sync.WaitGroup) {
			err = file.Write()
			if err != nil {
				m.logger.Errorf("文件 %s 处理失败, 错误: %v", file.Id, err)
			} else {
				m.logger.Infof("success: %s", file.Name)
				file.runtime.Events.Emit("file:complete", Complete{
					Id:     file.Id,
					Status: file.Status,
				})
			}
			w.Done()
		}(f, &wg)
	}
	wg.Wait()
	defer m.Destroy()
	return err
}

// 通过Id查找文件
func (m *Manager) getByIdFile(id string) *File {
	for _, file := range m.files {
		if file.Id == id {
			return file
		}
	}
	return nil
}

// 实例化Magick
func (m *Manager) SetMagick() {
	m.mw = NewMagick()

	for _, file := range m.files {
		file.SetMagick(m.mw)
	}
}

// 垃圾回收
func (m *Manager) Clear() {
	defer func() {
		imagick.Terminate()
		debug.FreeOSMemory()
	}()
	m.files = []*File{}
}

// 垃圾回收
func (m *Manager) Destroy() {
	m.logger.Infof("Destroy")
	m.mw.Destroy()
	m.mw = nil
}
