package localstore

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/vrischmann/userdir"
)

const ConfigDir = "imagemagick-ui"

// 应用程序配置文件保存的目录
type LocalStore struct {
	ConfDir string
}

// 创建实例
func NewLocalStore() *LocalStore {
	return &LocalStore{
		ConfDir: path.Join(userdir.GetConfigHome(), ConfigDir),
	}
}

// 读取配置文件
func (l LocalStore) Load(filename string) ([]byte, error) {
	p := path.Join(l.ConfDir, filename)
	d, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return d, err
}

// 将配置写入本地文件
func (l *LocalStore) Store(data []byte, filename string) error {
	p := path.Join(l.ConfDir, filename)
	if err := EnsureDirExists(l.ConfDir); err != nil {
		return err
	}
	if err := ioutil.WriteFile(p, data, 0777); err != nil {
		return err
	}
	return nil
}

// 检查文件是否存在，不存在则创建
func EnsureDirExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.Mkdir(path, 0777); err != nil {
			return err
		}
	}
	return nil
}
