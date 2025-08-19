package ui

import (
	"dhf-config-manager/internal/config"
	"dhf-config-manager/internal/models"
	"dhf-config-manager/internal/ui/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	fyneApp    fyne.App
	window     fyne.Window
	parser     *config.Parser
	schema     *models.Schema
	userConfig *models.UserConfig
	
	tree       *components.ConfigTree
	editor     *components.ConfigEditor
	toolbar    *components.Toolbar
}

func NewApp() *App {
	fyneApp := app.New()
	fyneApp.SetIcon(nil)
	
	// 设置中文支持 - 暂时禁用自定义主题
	// fyneApp.Settings().SetTheme(&chineseTheme{})
	
	window := fyneApp.NewWindow("DHF Configuration Manager")
	
	// 设置窗口属性 - 更合理的默认尺寸
	window.Resize(fyne.NewSize(900, 650))
	window.SetFixedSize(false) // 允许调整大小
	window.CenterOnScreen()
	
	return &App{
		fyneApp: fyneApp,
		window:  window,
		parser:  config.NewParser(),
	}
}

func (a *App) Initialize() error {
	a.tree = components.NewConfigTree()
	a.editor = components.NewConfigEditor()
	a.toolbar = components.NewToolbar()
	
	a.setupLayout()
	a.setupCallbacks()
	
	return nil
}

func (a *App) setupLayout() {
	// 创建主分割面板，设置合适的比例
	mainSplit := container.NewHSplit(
		a.tree.Container(),
		a.editor.Container(),
	)
	mainSplit.SetOffset(0.28) // 左侧占28%，右侧占72%
	
	// 设置左侧树形组件的合理宽度
	a.tree.Container().Resize(fyne.NewSize(220, 0))
	
	// 创建工具栏
	toolbar := container.NewVBox(
		a.toolbar.Container(),
		widget.NewSeparator(),
	)
	
	// 创建状态栏
	statusBar := container.NewBorder(
		nil, nil,
		widget.NewLabel("Ready"),
		widget.NewLabel("DHF Config Manager v1.0"),
		nil,
	)
	
	// 整体布局：顶部工具栏，中间主内容，底部状态栏
	content := container.NewBorder(
		toolbar,      // 顶部
		statusBar,    // 底部
		nil,          // 左侧
		nil,          // 右侧
		mainSplit,    // 中心内容
	)
	
	a.window.SetContent(content)
}

func (a *App) setupCallbacks() {
	a.tree.SetSelectionCallback(func(nodeID string) {
		a.editor.ShowSection(nodeID)
	})
	
	a.toolbar.SetNewCallback(func() {
		a.userConfig = &models.UserConfig{Values: make(map[string]interface{})}
		a.editor.SetConfig(a.userConfig)
		a.refreshTree()
	})
	
	a.toolbar.SetOpenCallback(func() {
	})
	
	a.toolbar.SetSaveCallback(func() {
	})
	
	a.toolbar.SetExportCallback(func() {
	})
}

func (a *App) refreshTree() {
	if a.schema != nil {
		a.tree.LoadSchema(a.schema)
	}
}

func (a *App) LoadSchema(filePath string) error {
	if err := a.parser.LoadSchema(filePath); err != nil {
		return err
	}
	
	a.schema = a.parser.GetSchema()
	a.tree.LoadSchema(a.schema)
	a.editor.SetSchema(a.schema)
	
	return nil
}

func (a *App) Run() {
	a.window.ShowAndRun()
}