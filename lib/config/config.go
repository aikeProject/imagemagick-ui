package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"imagemagick-ui/lib/localstore"
	"os"
	"path"
	"path/filepath"

	"github.com/wailsapp/wails"
)

const filename = "conf.json"

// 本地配置
type App struct {
	OutDir string `json:"outDir"` // 文件保存目录
	Target string `json:"target"` // 文件目标类型 png/jpg/webp...
}

// 应用程序配置
type Config struct {
	App        *App
	LocalStore *localstore.LocalStore
	Runtime    *wails.Runtime
	Logger     *wails.CustomLogger
}

// 返回"Config"实例
func NewConfig() *Config {
	c := &Config{}
	c.LocalStore = localstore.NewLocalStore()
	file, err := c.LocalStore.Load(filename)
	if err != nil {
		c.App, _ = defaults()
	}
	if err := json.Unmarshal(file, &c.App); err != nil {
		fmt.Printf("error %v", err)
	}
	return c
}

// 获取配置
func (c *Config) GetAppConfig() map[string]interface{} {
	return map[string]interface{}{
		"outDir": c.App.OutDir,
		"target": c.App.Target,
	}
}

func (c *Config) WailsInit(runtime *wails.Runtime) error {
	c.Runtime = runtime
	c.Logger = c.Runtime.Log.New("Config")
	c.Logger.Info("Config 初始化...")
	return nil
}

// 打开对话框，选择输出目录
func (c *Config) SetOutDir() string {
	dir := c.Runtime.Dialog.SelectDirectory()
	if dir != "" {
		c.App.OutDir = dir
		c.Logger.Infof("输出目录: %s", dir)
		if err := c.store(); err != nil {
			c.Logger.Errorf("配置保存失败：%v", err)
		}
	}
	return c.App.OutDir
}

// 打开输出目录
func (c *Config) OpenOutputDir() error {
	if err := c.Runtime.Browser.OpenURL(c.App.OutDir); err != nil {
		return err
	}
	return nil
}

// 保存配置
func (c *Config) SetConfig(cfg string) error {
	a := &App{}
	if err := json.Unmarshal([]byte(cfg), &a); err != nil {
		c.Logger.Errorf("failed to unmarshal config: %v", err)
		return err
	}
	c.App = a
	if err := c.store(); err != nil {
		c.Logger.Errorf("failed to store config: %v", err)
		return err
	}
	return nil
}

// 保存配置到配置文件
func (c *Config) store() error {
	js, err := json.Marshal(c.GetAppConfig())
	if err != nil {
		c.Logger.Errorf("应用配置解析失败: %v", err)
		return err
	}
	if err := c.LocalStore.Store(js, filename); err != nil {
		c.Logger.Errorf("应用配置保存失败: %v", err)
		return err
	}
	return nil
}

// 重置为默认配置
func (c *Config) RestoreDefaults() error {
	app, err := defaults()
	if err != nil {
		return err
	}
	c.App = app
	if err := c.store(); err != nil {
		return err
	}
	return nil
}

// 默认配置
func defaults() (*App, error) {
	a := &App{
		Target: "png",
	}
	ud, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get user directory: %v", err)
		return nil, errors.New("无法获取用户目录")
	}

	od := path.Join(ud, "Desktop", localstore.ConfigDir)
	cp := filepath.Clean(od)

	if _, err := os.Stat(od); os.IsNotExist(err) {
		if err := os.Mkdir(od, 0777); err != nil {
			od = "./"
			fmt.Printf("failed to create default output directory: %v", err)
			return nil, errors.New("无法创建默认输出目录")
		}
	}
	a.OutDir = cp
	return a, nil
}
