package components

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	
	"configcraft/internal/version"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type Toolbar struct {
	container *fyne.Container
	window    fyne.Window
	
	openCallback func(filePath string)
	saveCallback func(filePath string)
	hasOpenFile  func() bool // 检查是否有已打开的文件
}

func (t *Toolbar) SetWindow(window fyne.Window) {
	t.window = window
}

func NewToolbar() *Toolbar {
	toolbar := &Toolbar{}
	
	// 创建OPEN按钮
	openBtn := widget.NewButton("打开配置", func() {
		toolbar.showOpenDialog()
	})
	openBtn.Importance = widget.MediumImportance
	
	// 创建SAVE按钮
	saveBtn := widget.NewButton("保存配置", func() {
		toolbar.showSaveDialog()
	})
	saveBtn.Importance = widget.HighImportance // 高亮保存按钮
	
	// 创建About按钮
	aboutBtn := widget.NewButton("关于", func() {
		toolbar.showAboutDialog()
	})
	aboutBtn.Importance = widget.LowImportance
	
	// 创建简洁的3按钮工具栏
	toolbar.container = container.NewHBox(
		openBtn,
		saveBtn,
		widget.NewSeparator(),
		aboutBtn,
	)
	
	return toolbar
}

func (t *Toolbar) showAboutDialog() {
	if t.window == nil {
		return
	}
	
	// 创建现代化的关于对话框内容
	appTitle := widget.NewLabelWithStyle(version.AppName, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	
	versionLabel := widget.NewLabelWithStyle("Version "+version.Version, fyne.TextAlignCenter, fyne.TextStyle{})
	authorLabel := widget.NewLabelWithStyle("Created by Felix", fyne.TextAlignCenter, fyne.TextStyle{})
	
	description := widget.NewLabelWithStyle(
		"Universal Configuration Management Tool\n将复杂配置文件转换为友好的图形界面\n\nA modern visual configuration tool with YAML support",
		fyne.TextAlignCenter, 
		fyne.TextStyle{},
	)
	
	// 添加一些间距和分隔线
	content := container.NewVBox(
		appTitle,
		widget.NewSeparator(),
		container.NewPadded(container.NewVBox(
			versionLabel,
			authorLabel,
		)),
		widget.NewSeparator(),
		container.NewPadded(description),
	)
	
	aboutDialog := dialog.NewCustom("About", "Close", content, t.window)
	aboutDialog.Resize(fyne.NewSize(450, 300))
	aboutDialog.Show()
}

// showOpenDialog 显示文件打开对话框（使用zenity原生对话框）
func (t *Toolbar) showOpenDialog() {
	if t.window == nil {
		return
	}

	// 使用zenity原生文件对话框
	zenityDialog := NewZenityFileDialog()
	
	// 显示打开对话框
	filePath, err := zenityDialog.ShowOpenDialog("选择配置文件")
	if err != nil {
		// 如果是用户取消，不显示错误
		if !strings.Contains(err.Error(), "用户取消") {
			dialog.ShowError(err, t.window)
		}
		return
	}

	// 验证文件格式
	if err := zenityDialog.ValidateYAMLFile(filePath); err != nil {
		dialog.ShowError(err, t.window)
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		dialog.ShowError(fmt.Errorf("文件不存在或无法访问: %v", err), t.window)
		return
	}

	// 调用打开回调
	if t.openCallback != nil {
		t.openCallback(filePath)
	}
}

// showSaveDialog 智能保存 - 优先直接保存，否则显示对话框
func (t *Toolbar) showSaveDialog() {
	if t.window == nil {
		return
	}

	// 如果已有打开的文件，直接保存
	if t.hasOpenFile != nil && t.hasOpenFile() {
		if t.saveCallback != nil {
			t.saveCallback("") // 空路径表示保存到原文件
		}
		return
	}

	// 否则显示保存对话框选择新位置（混合方案：Fyne对话框 + 目录切换技巧）
	currentDir, err := os.Getwd()
	if err == nil {
		log.Printf("Current working directory: %s", currentDir)
		// 临时切换到当前目录，确保对话框从正确位置开始
		originalDir, _ := os.Getwd()
		if err := os.Chdir(currentDir); err == nil {
			defer func() {
				os.Chdir(originalDir)
			}()
		}
	}

	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, t.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		filePath := writer.URI().Path()
		log.Printf("Selected file to save: %s", filePath)

		// 确保文件扩展名为.yaml
		if filepath.Ext(filePath) == "" {
			filePath += ".yaml"
		}

		// 调用保存回调
		if t.saveCallback != nil {
			t.saveCallback(filePath)
		}
	}, t.window)

	// 设置默认文件名和过滤器
	saveDialog.SetFileName("dhf_config.yaml")
	saveDialog.SetFilter(storage.NewExtensionFileFilter([]string{".yaml", ".yml"}))
	
	saveDialog.Resize(fyne.NewSize(900, 650))
	saveDialog.Show()
}

// SetOpenCallback 设置打开文件回调
func (t *Toolbar) SetOpenCallback(callback func(filePath string)) {
	t.openCallback = callback
}

// SetSaveCallback 设置保存文件回调
func (t *Toolbar) SetSaveCallback(callback func(filePath string)) {
	t.saveCallback = callback
}

// SetHasOpenFileCallback 设置检查是否有打开文件的回调
func (t *Toolbar) SetHasOpenFileCallback(callback func() bool) {
	t.hasOpenFile = callback
}

func (t *Toolbar) Container() *fyne.Container {
	return t.container
}