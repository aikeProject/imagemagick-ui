package lib

import "os"

// 文件操作集合
type FSHelper struct {
}

func NewFSHelper() *FSHelper {
	result := &FSHelper{}
	return result
}

// 检查文件夹是否存在
func (h *FSHelper) DirExists(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	// 检测给定路径是否是目录
	return info.Mode().IsDir()
}

// 创建该路径上的所有目录
func (h FSHelper) MkDirs(path string, mode ...os.FileMode) error {
	var perms os.FileMode
	perms = 0700
	if len(mode) == 1 {
		perms = mode[0]
	}
	return os.MkdirAll(path, perms)
}
