package core

import (
	"encoding/base64"
	"encoding/json"
)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) HandleFile(fileJson string) error {
	file := &File{}
	if err := json.Unmarshal([]byte(fileJson), &file); err != nil {
		return err
	}
	_, err := base64.StdEncoding.DecodeString(file.Data)
	if err != nil {
		return err
	}
	return nil
}
