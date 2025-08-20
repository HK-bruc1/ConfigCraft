package components

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	
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
	appTitle := widget.NewLabelWithStyle("DHF Configuration Manager", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	
	versionLabel := widget.NewLabelWithStyle("Version 0.2.1", fyne.TextAlignCenter, fyne.TextStyle{})
	authorLabel := widget.NewLabelWithStyle("Created by Felix", fyne.TextAlignCenter, fyne.TextStyle{})
	
	description := widget.NewLabelWithStyle(
		"DHF AC710N-V300P03 SDK 现代化可视配置工具\n将复杂配置文件转换为友好的图形界面\n\nA modern visual configuration tool for DHF SDK",
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

// showOpenDialog 显示文件打开对话框
func (t *Toolbar) showOpenDialog() {
	if t.window == nil {
		return
	}

	// 创建文件选择对话框，只允许选择YAML文件
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, t.window)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		filePath := reader.URI().Path()
		log.Printf("Selected file to open: %s", filePath)

		// 检查文件扩展名
		ext := filepath.Ext(filePath)
		if ext != ".yaml" && ext != ".yml" {
			dialog.ShowError(fmt.Errorf("请选择YAML格式文件（.yaml或.yml）"), t.window)
			return
		}

		// 调用打开回调
		if t.openCallback != nil {
			t.openCallback(filePath)
		}
	}, t.window)

	// 设置文件过滤器
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".yaml", ".yml"}))
	
	// 正确设置起始目录 - 创建有效的ListableURI
	if workDir, err := os.Getwd(); err == nil {
		log.Printf("Setting file dialog to start from: %s", workDir)
		
		// 创建文件URI - 在Windows上需要正确的格式
		var dirURI fyne.URI
		if strings.HasPrefix(workDir, "/") {
			// Unix路径
			dirURI = storage.NewFileURI(workDir)
		} else {
			// Windows路径 - 手动构造file:// URI
			winPath := strings.ReplaceAll(workDir, "\\", "/")
			if !strings.HasPrefix(winPath, "/") {
				winPath = "/" + winPath
			}
			if uri, err := storage.ParseURI("file://" + winPath); err == nil {
				dirURI = uri
			} else {
				log.Printf("Failed to parse URI: %v", err)
				dirURI = storage.NewFileURI(workDir)
			}
		}
		
		// 尝试设置位置，需要转换为ListableURI
		if dirURI != nil {
			if listableURI, ok := dirURI.(fyne.ListableURI); ok {
				fileDialog.SetLocation(listableURI)
				log.Printf("Set file dialog location to: %s", dirURI.String())
			} else {
				log.Printf("Warning: URI is not ListableURI: %s", dirURI.String())
			}
		}
	}
	
	fileDialog.Resize(fyne.NewSize(800, 600))
	fileDialog.Show()
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

	// 否则显示保存对话框选择新位置
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
	
	// 正确设置保存对话框的起始目录
	if workDir, err := os.Getwd(); err == nil {
		log.Printf("Setting save dialog to start from: %s", workDir)
		
		// 创建文件URI - 在Windows上需要正确的格式
		var dirURI fyne.URI
		if strings.HasPrefix(workDir, "/") {
			// Unix路径
			dirURI = storage.NewFileURI(workDir)
		} else {
			// Windows路径 - 手动构造file:// URI
			winPath := strings.ReplaceAll(workDir, "\\", "/")
			if !strings.HasPrefix(winPath, "/") {
				winPath = "/" + winPath
			}
			if uri, err := storage.ParseURI("file://" + winPath); err == nil {
				dirURI = uri
			} else {
				log.Printf("Failed to parse URI: %v", err)
				dirURI = storage.NewFileURI(workDir)
			}
		}
		
		// 尝试设置位置，需要转换为ListableURI
		if dirURI != nil {
			if listableURI, ok := dirURI.(fyne.ListableURI); ok {
				saveDialog.SetLocation(listableURI)
				log.Printf("Set save dialog location to: %s", dirURI.String())
			} else {
				log.Printf("Warning: URI is not ListableURI: %s", dirURI.String())
			}
		}
	}
	
	saveDialog.Resize(fyne.NewSize(800, 600))
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