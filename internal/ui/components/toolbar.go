package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Toolbar struct {
	container *fyne.Container
	
	newCallback    func()
	openCallback   func()
	saveCallback   func()
	exportCallback func()
}

func NewToolbar() *Toolbar {
	toolbar := &Toolbar{}
	
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
	
	exportBtn := widget.NewButton("Export", func() {
		if toolbar.exportCallback != nil {
			toolbar.exportCallback()
		}
	})
	
	// 创建紧凑的工具栏
	toolbar.container = container.NewHBox(
		newBtn,
		openBtn,
		saveBtn,
		widget.NewSeparator(),
		exportBtn,
	)
	
	return toolbar
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