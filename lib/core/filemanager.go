package core

type Manager struct {
	files []*File
}

func NewManager() *Manager {
	return &Manager{}
}

// 将文件添加至"Manager"
func (m *Manager) HandleFile(fileJson string) error {
	if file, err := NewFile(fileJson); err != nil {
		m.files = append(m.files, file)
		return err
	}
	return nil
}

// 并发处理文件
func (m *Manager) Convert() (errs []error) {
	if err := m.files[0].Write(); err != nil {
		return append(errs, err)
	}
	return errs
}
