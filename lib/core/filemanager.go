package core

import (
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
	m.logger.Infof("文件添加至 Manager <= %s", file.Name)
	return nil
}

// 并发处理文件
func (m *Manager) Convert() (errs []error) {
	var wg sync.WaitGroup
	wg.Add(m.countUnconverted())
	for _, f := range m.files {
		if f.Status == Done {
			continue
		}
		go func(file *File, w *sync.WaitGroup) {
			if err := file.Write(); err != nil {
				m.logger.Errorf("文件 %s 处理失败, 错误: %v", file.Id, err)
				errs = append(errs, err)
			} else {
				m.logger.Infof("处理完成的文件: %s", file.Name)
				file.runtime.Events.Emit("file:complete", Complete{
					Id:     file.Id,
					Status: file.Status,
				})
			}
			w.Done()
		}(f, &wg)
	}
	wg.Wait()
	return errs
}

// 垃圾回收
func (m *Manager) Clear() {
	defer func() {
		imagick.Terminate()
		debug.FreeOSMemory()
	}()
	m.files = []*File{}
}

// 尚未处理的文件数量
func (m *Manager) countUnconverted() int {
	c := 0
	for _, file := range m.files {
		if file.Status != Done {
			c++
		}
	}
	return c
}
