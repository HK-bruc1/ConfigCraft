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
	
	// 设置窗口属性 - 保持原定尺寸
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
	// 设置toolbar的window引用
	a.toolbar.SetWindow(a.window)
	
	// 左侧区域：配置分组导航
	leftPanel := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Configuration Groups", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		a.tree.Container(),
	)
	
	// 右侧区域：详细配置选项
	rightPanel := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Configuration Options", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		a.editor.Container(),
	)
	
	// 创建主分割面板
	mainSplit := container.NewHSplit(leftPanel, rightPanel)
	mainSplit.SetOffset(0.3) // 左侧30%，右侧70%
	
	// 整体布局：工具栏在顶部，分界线，主要内容区域
	content := container.NewBorder(
		// 顶部：工具栏 + 分界线
		container.NewVBox(
			a.toolbar.Container(),
			widget.NewSeparator(),
		),
		// 底部：状态栏
		container.NewBorder(
			widget.NewSeparator(), nil,
			widget.NewLabel("Ready"),
			widget.NewLabel("DHF Configuration Manager"),
			nil,
		),
		nil, nil, // 左右留空
		// 中心：主要内容区域
		mainSplit,
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