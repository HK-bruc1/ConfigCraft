package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Toolbar struct {
	container *fyne.Container
	window    fyne.Window
	
	newCallback    func()
	openCallback   func()
	saveCallback   func()
	exportCallback func()
}

func (t *Toolbar) SetWindow(window fyne.Window) {
	t.window = window
}

func NewToolbar() *Toolbar {
	toolbar := &Toolbar{}
	
	// 创建简洁的工具栏按钮
	newBtn := widget.NewButton("New", func() {
		if toolbar.newCallback != nil {
			toolbar.newCallback()
		}
	})
	
	openBtn := widget.NewButton("Open", func() {
		if toolbar.openCallback != nil {
			toolbar.openCallback()
		}
	})
	
	saveBtn := widget.NewButton("Save", func() {
		if toolbar.saveCallback != nil {
			toolbar.saveCallback()
		}
	})
	
	exportBtn := widget.NewButton("Export to DHF", func() {
		if toolbar.exportCallback != nil {
			toolbar.exportCallback()
		}
	})
	exportBtn.Importance = widget.HighImportance // 高亮导出按钮
	
	// 添加关于按钮
	aboutBtn := widget.NewButton("About", func() {
		toolbar.showAboutDialog()
	})
	aboutBtn.Importance = widget.LowImportance
	
	// 创建单行工具栏布局
	toolbar.container = container.NewHBox(
		newBtn,
		openBtn,
		saveBtn,
		widget.NewSeparator(),
		exportBtn,
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

func (t *Toolbar) Container() *fyne.Container {
	return t.container
}

func (t *Toolbar) SetNewCallback(callback func()) {
	t.newCallback = callback
}

func (t *Toolbar) SetOpenCallback(callback func()) {
	t.openCallback = callback
}

func (t *Toolbar) SetSaveCallback(callback func()) {
	t.saveCallback = callback
}

func (t *Toolbar) SetExportCallback(callback func()) {
	t.exportCallback = callback
}