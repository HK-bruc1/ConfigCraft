package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ncruces/zenity"
)

// ZenityFileDialog 提供使用zenity的跨平台原生文件对话框
type ZenityFileDialog struct{}

// NewZenityFileDialog 创建zenity文件对话框实例
func NewZenityFileDialog() *ZenityFileDialog {
	return &ZenityFileDialog{}
}

// ShowOpenDialog 显示文件打开对话框
func (zfd *ZenityFileDialog) ShowOpenDialog(title string) (string, error) {
	
	// 设置zenity选项
	options := []zenity.Option{
		zenity.Title(title),
		zenity.FileFilter{
			Name:     "YAML配置文件",
			Patterns: []string{"*.yaml", "*.yml"},
		},
	}
	
	// zenity会自动从当前工作目录开始，无需切换目录
	// 如果需要指定特定目录，可以添加Filename选项
	
	// 显示文件选择对话框
	filePath, err := zenity.SelectFile(options...)
	if err != nil {
		if err == zenity.ErrCanceled {
			return "", fmt.Errorf("用户取消了文件选择")
		}
		return "", fmt.Errorf("文件对话框错误: %v", err)
	}
	
	// 规范化文件路径 - 确保使用正确的路径分隔符
	filePath = filepath.Clean(filePath)
	
	return filePath, nil
}

// ShowSaveDialog 显示文件保存对话框
func (zfd *ZenityFileDialog) ShowSaveDialog(title, defaultName string) (string, error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}
	
	// 构造默认文件路径
	defaultPath := filepath.Join(currentDir, defaultName)
	
	// 设置zenity选项
	options := []zenity.Option{
		zenity.Title(title),
		zenity.ConfirmOverwrite(),
		zenity.FileFilter{
			Name:     "YAML配置文件",
			Patterns: []string{"*.yaml", "*.yml"},
		},
		zenity.Filename(defaultPath), // 设置默认文件路径
	}
	
	// 显示保存对话框
	filePath, err := zenity.SelectFileSave(options...)
	if err != nil {
		if err == zenity.ErrCanceled {
			return "", fmt.Errorf("用户取消了文件选择")
		}
		return "", fmt.Errorf("文件对话框错误: %v", err)
	}
	
	// 规范化文件路径 - 确保使用正确的路径分隔符
	filePath = filepath.Clean(filePath)
	
	return filePath, nil
}

// ValidateYAMLFile 验证文件是否为YAML格式
func (zfd *ZenityFileDialog) ValidateYAMLFile(filePath string) error {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".yaml" && ext != ".yml" {
		return fmt.Errorf("请选择YAML格式文件（.yaml或.yml）")
	}
	return nil
}

// GetCurrentDirectory 获取当前工作目录
func (zfd *ZenityFileDialog) GetCurrentDirectory() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return ""
}