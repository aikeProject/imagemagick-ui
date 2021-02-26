package core

import (
	"encoding/json"
	"imagemagick-ui/lib"
	"path"
	"runtime/debug"
	"sync"

	"github.com/thoas/go-funk"

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
	var files []*File
	var ids []string

	err = json.Unmarshal([]byte(idStr), &ids)
	if err != nil {
		return err
	}

	if len(ids) > 0 {
		for _, id := range ids {
			if file := m.getByIdFile(id); file != nil {
				files = append(files, file)
			}
		}
	} else {
		files = m.files
	}
	if len(files) == 0 {
		return nil
	}
	// 将文件状态重置回未处理状态
	for _, file := range files {
		file.Status = SendSuccess
	}

	// 初始化Magick图片处理实例
	m.SetMagick()

	// xxx.png xxx1.png => xxx.gif
	if funk.ContainsString(covertExtFiles, "."+m.conf.App.Target) {
		return m.Write(files)
	}

	// xxx.png => xxx.jpg
	err = m.worker(files)
	return err
}

// 处理文件集合，合并图像
// 例如：将多张图片装换为.gif格式
// xxx.png xxx1.png => xxx.gif
func (m *Manager) Write(files []*File) error {
	for _, file := range files {
		file.Status = Running
		bytes, err := file.Decode()
		if err != nil {
			return err
		}
		if err := m.mw.ReadImageBlob(bytes); err != nil {
			return err
		}
		m.mw.Resize(m.conf.App.Width, m.conf.App.Height)
		// 设置gif帧率
		if err := m.mw.SetImageDelay(5); err != nil {
			return err
		}
		//if err := m.mw.SetImageDispose(imagick.DISPOSE_NONE); err != nil {
		//	return err
		//}
	}
	filename := files[0].rename()
	// adjoin true 多个文件合并为一个文件
	//m.mw.CoalesceImages()
	err := m.mw.WriteImages(path.Join(m.conf.App.OutDir, filename), true)
	if err != nil {
		m.logger.Errorf("文件 %s 转换失败, 错误: %v", filename, err)
	}
	funk.ForEach(files, func(v *File) {
		v.Status = Done
		m.logger.Infof("success: %s", v.Name)
		v.runtime.Events.Emit("file:complete", Complete{
			Id:     v.Id,
			Status: v.Status,
		})
	})
	defer m.Destroy()
	return err
}

// 并发处理多张图片
func (m *Manager) worker(files []*File) (err error) {
	var wg sync.WaitGroup
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
