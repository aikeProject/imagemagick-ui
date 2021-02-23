package core

import (
	"imagemagick-ui/lib/config"
	"path"
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
	config  *config.Config
	mw      *imagick.MagickWand
	runtime *wails.Runtime
	logger  *wails.CustomLogger
}

func NewManager(config *config.Config) *Manager {
	return &Manager{
		config: config,
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
	file, err := NewFile(fileJson)
	if err != nil {
		return err
	}
	file.WailsInit(m.runtime)
	m.files = append(m.files, file)
	m.logger.Infof("文件添加至 Manager <= %s", file.Name)
	if m.mw == nil {
		m.mw = imagick.NewMagickWand()
	}
	return nil
}

// 并发处理文件
func (m *Manager) Convert() (errs []error) {
	var wg sync.WaitGroup
	wg.Add(m.countUnconverted())
	for _, f := range m.files {
		go func(file *File, w *sync.WaitGroup) {
			bytes, err := file.Decode()
			if err != nil {
				errs = append(errs, err)
			}
			err = m.mw.ReadImageBlob(bytes)
			if err != nil {
				errs = append(errs, err)
			}
			width := m.mw.GetImageWidth()
			height := m.mw.GetImageHeight()
			m.logger.Infof("width %v, height %v", width, height)
			// 保持图像纵横比
			if width > height {
				height = uint((float32(200) / float32(width)) * float32(height))
				width = 200
			} else if width < height {
				width = uint((200 / float32(height)) * float32(width))
				height = 200
			}
			m.logger.Infof("width %v, height %v", width, height)
			if err := m.mw.AdaptiveResizeImage(width, height); err != nil {
				errs = append(errs, err)
			}
			if err := m.mw.WriteImage(path.Join(m.config.App.OutDir, file.Name)); err != nil {
				errs = append(errs, err)
			}
			w.Done()
		}(f, &wg)
	}
	wg.Wait()
	//err := m.mw.WriteImages(path.Join(m.config.App.OutDir, "thumbnail%03d.png"), true)
	//errs = append(errs, err)
	//var wg sync.WaitGroup
	//wg.Add(m.countUnconverted())
	//for _, f := range m.files {
	//	if f.Status == Done {
	//		continue
	//	}
	//	go func(file *File, w *sync.WaitGroup) {
	//		if err := file.Write(); err != nil {
	//			m.logger.Errorf("文件 %s 处理失败, 错误: %v", file.Id, err)
	//			errs = append(errs, err)
	//		} else {
	//			m.logger.Infof("处理完成的文件: %s", file.Name)
	//			file.runtime.Events.Emit("file:complete", Complete{
	//				Id:     file.Id,
	//				Status: file.Status,
	//			})
	//		}
	//		w.Done()
	//	}(f, &wg)
	//}
	//wg.Wait()
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

// 垃圾回收
func (m *Manager) Destroy() {
	m.logger.Infof("Destroy")
	m.mw.Destroy()
	m.mw = nil
}
