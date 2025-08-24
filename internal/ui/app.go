package ui

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	
	"dhf-config-manager/internal/config"
	"dhf-config-manager/internal/models"
	"dhf-config-manager/internal/ui/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	
	currentFilePath string // 记录当前打开的文件路径
}

func NewApp() *App {
	fyneApp := app.New()
	fyneApp.SetIcon(nil)
	
	// 使用简化的中文主题支持
	fyneApp.Settings().SetTheme(NewSimpleChineseTheme())
	
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
	// 设置toolbar和editor的window引用
	a.toolbar.SetWindow(a.window)
	a.editor.SetWindow(a.window)
	
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
	
	a.toolbar.SetOpenCallback(func(filePath string) {
		a.openConfigFile(filePath)
	})
	
	a.toolbar.SetSaveCallback(func(filePath string) {
		a.saveConfigFile(filePath)
	})
	
	a.toolbar.SetHasOpenFileCallback(func() bool {
		return a.currentFilePath != ""
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

// openConfigFile 打开YAML配置文件 - 支持schema和配置文件
func (a *App) openConfigFile(filePath string) {
	log.Printf("Opening config file: %s", filePath)
	
	// 尝试作为schema文件加载
	if err := a.parser.LoadSchema(filePath); err == nil {
		// 成功加载为schema文件
		a.schema = a.parser.GetSchema()
		a.userConfig = &models.UserConfig{Values: make(map[string]interface{})}
		a.currentFilePath = ""  // schema文件不是配置文件
		a.editor.SetSchema(a.schema)
		a.editor.SetConfig(a.userConfig)
		
		// 刷新界面
		a.refreshTree()
		
		// 自动选择第一个section
		if len(a.schema.Sections) > 0 {
			sectionKeys := make([]string, 0, len(a.schema.Sections))
			for sectionKey := range a.schema.Sections {
				sectionKeys = append(sectionKeys, sectionKey)
			}
			// 使用相同的排序逻辑
			sectionOrder := map[string]int{
				"basic": 1, "call_actions": 2, "music_actions": 3, "led_config": 4, "special_functions": 5, "advanced": 6,
			}
			sort.Slice(sectionKeys, func(i, j int) bool {
				orderI, existsI := sectionOrder[sectionKeys[i]]
				orderJ, existsJ := sectionOrder[sectionKeys[j]]
				if existsI && existsJ {
					return orderI < orderJ
				} else if existsI {
					return true
				} else if existsJ {
					return false
				}
				return sectionKeys[i] < sectionKeys[j]
			})
			
			a.editor.ShowSection(sectionKeys[0])
		}
		
		// 显示成功消息
		message := fmt.Sprintf("Schema文件已成功加载！\n\n文件路径: %s\n配置分组数: %d\n支持增强功能: 描述信息、提示、可编辑下拉框", 
			filePath, len(a.schema.Sections))
		dialog.ShowInformation("Schema加载成功", message, a.window)
		return
	}
	
	// 作为用户配置文件加载
	userConfig, err := a.parser.LoadUserConfig(filePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("无法加载配置文件: %v", err), a.window)
		return
	}
	
	// 从配置文件内容动态生成schema
	dynamicSchema := a.generateSchemaFromConfig(userConfig)
	
	// 更新应用状态
	a.schema = dynamicSchema
	a.userConfig = userConfig
	a.currentFilePath = filePath // 记录当前文件路径
	a.editor.SetSchema(a.schema)
	a.editor.SetConfig(a.userConfig)
	
	// 刷新界面
	a.refreshTree()
	
	// 自动选择第一个section
	if len(a.schema.Sections) > 0 {
		for sectionKey := range a.schema.Sections {
			a.editor.ShowSection(sectionKey)
			break
		}
	}
	
	// 显示成功消息
	message := fmt.Sprintf("配置文件已成功加载！\n\n文件路径: %s\n配置项数: %d\n自动识别分组数: %d", 
		filePath, len(a.userConfig.Values), len(a.schema.Sections))
	dialog.ShowInformation("打开成功", message, a.window)
}

// saveConfigFile 保存配置文件并生成conf文件 - 智能保存版本
func (a *App) saveConfigFile(requestedPath string) {
	// 如果有当前文件路径，直接保存到原文件；否则使用用户选择的路径
	var targetPath string
	if a.currentFilePath != "" {
		targetPath = a.currentFilePath
		log.Printf("Saving to original file: %s", targetPath)
	} else {
		targetPath = requestedPath
		log.Printf("Saving to new file: %s", targetPath)
	}
	
	if a.userConfig == nil {
		dialog.ShowError(fmt.Errorf("没有可保存的配置数据"), a.window)
		return
	}
	
	// 使用parser保存配置并生成conf文件
	if err := a.parser.SaveConfigWithConf(a.userConfig, targetPath); err != nil {
		dialog.ShowError(err, a.window)
		return
	}
	
	// 生成conf文件路径
	confPath := strings.TrimSuffix(targetPath, filepath.Ext(targetPath)) + ".conf"
	
	// 显示成功消息
	var message string
	if a.currentFilePath != "" {
		message = fmt.Sprintf("配置已成功保存并覆盖原文件！\n\nYAML配置: %s\nDHF配置: %s", targetPath, confPath)
	} else {
		message = fmt.Sprintf("配置保存成功！\n\nYAML配置: %s\nDHF配置: %s", targetPath, confPath)
	}
	
	dialog.ShowInformation("保存成功", message, a.window)
	log.Printf("Successfully saved YAML and generated conf file")
}

// generateSchemaFromConfig 从配置文件动态生成schema
func (a *App) generateSchemaFromConfig(userConfig *models.UserConfig) *models.Schema {
	schema := &models.Schema{
		SchemaVersion: "1.0",
		DisplayName:   "动态配置",
		Sections:      make(map[string]models.ConfigSection),
	}
	
	// 按键路径分组配置项
	sectionGroups := make(map[string]map[string]interface{})
	
	for keyPath, value := range userConfig.Values {
		parts := strings.Split(keyPath, ".")
		
		var sectionKey, fieldKey string
		
		if len(parts) == 2 {
			// 二级结构: section.field
			sectionKey = parts[0]
			fieldKey = parts[1]
		} else if len(parts) >= 3 {
			// 三级或更多: section.group.field 或 section.group.subgroup.field
			sectionKey = parts[0]
			fieldKey = strings.Join(parts[1:], ".")
		} else {
			// 一级结构: 放入"基础配置"分组
			sectionKey = "basic"
			fieldKey = parts[0]
		}
		
		if sectionGroups[sectionKey] == nil {
			sectionGroups[sectionKey] = make(map[string]interface{})
		}
		sectionGroups[sectionKey][fieldKey] = value
	}
	
	// 为每个section创建配置 - 按固定顺序遍历以确保界面显示一致
	sectionKeys := make([]string, 0, len(sectionGroups))
	for sectionKey := range sectionGroups {
		sectionKeys = append(sectionKeys, sectionKey)
	}
	// 定义section的显示顺序优先级
	sectionOrder := map[string]int{
		"basic":      1,
		"key_actions": 2,
		"led_config": 3,
		"factory":    4,
		"advanced":   5,
	}
	// 按优先级和字典序排序
	sort.Slice(sectionKeys, func(i, j int) bool {
		orderI, existsI := sectionOrder[sectionKeys[i]]
		orderJ, existsJ := sectionOrder[sectionKeys[j]]
		
		if existsI && existsJ {
			return orderI < orderJ
		} else if existsI {
			return true // 有优先级的排在前面
		} else if existsJ {
			return false
		}
		return sectionKeys[i] < sectionKeys[j] // 都没有优先级时按字典序
	})
	
	for _, sectionKey := range sectionKeys {
		fields := sectionGroups[sectionKey]
		section := models.ConfigSection{
			Name:   a.getSectionDisplayName(sectionKey),
			Icon:   "settings",
			Fields: make(map[string]models.ConfigField),
			Groups: make(map[string]models.ConfigGroup),
		}
		
		// 处理字段 - 按固定顺序遍历以确保字段显示一致
		fieldKeys := make([]string, 0, len(fields))
		for fieldKey := range fields {
			fieldKeys = append(fieldKeys, fieldKey)
		}
		sort.Strings(fieldKeys) // 按字典序排序字段
		
		for _, fieldKey := range fieldKeys {
			value := fields[fieldKey]
			fieldParts := strings.Split(fieldKey, ".")
			
			if len(fieldParts) == 1 {
				// 直接字段
				section.Fields[fieldKey] = a.createFieldFromValue(fieldKey, value)
			} else {
				// 分组字段
				groupKey := fieldParts[0]
				subFieldKey := strings.Join(fieldParts[1:], ".")
				
				if _, exists := section.Groups[groupKey]; !exists {
					section.Groups[groupKey] = models.ConfigGroup{
						Name:   a.getGroupDisplayName(groupKey),
						Fields: make(map[string]models.ConfigField),
					}
				}
				
				group := section.Groups[groupKey]
				group.Fields[subFieldKey] = a.createFieldFromValue(subFieldKey, value)
				section.Groups[groupKey] = group
			}
		}
		
		schema.Sections[sectionKey] = section
	}
	
	log.Printf("Generated schema with %d sections", len(schema.Sections))
	return schema
}

// createFieldFromValue 根据值的类型创建配置字段
func (a *App) createFieldFromValue(key string, value interface{}) models.ConfigField {
	field := models.ConfigField{
		Label:    a.getFieldDisplayName(key),
		Default:  value,
		Required: false,
	}
	
	switch v := value.(type) {
	case bool:
		field.Type = "boolean"
	case int, int64, float64:
		field.Type = "number"
	case string:
		// 如果是已知的枚举值，创建选择框
		if a.isEnumValue(v) {
			field.Type = "select"
			field.Options = a.getEnumOptions(v)
		} else {
			field.Type = "text"
		}
	default:
		field.Type = "text"
	}
	
	return field
}

// getSectionDisplayName 获取section的显示名称
func (a *App) getSectionDisplayName(key string) string {
	nameMap := map[string]string{
		"basic":      "基础配置",
		"key_actions": "按键配置",
		"led_config": "LED配置",
		"factory":    "工厂设置",
		"advanced":   "高级设置",
	}
	
	if name, exists := nameMap[key]; exists {
		return name
	}
	return strings.Title(key) + "配置"
}

// getGroupDisplayName 获取group的显示名称
func (a *App) getGroupDisplayName(key string) string {
	nameMap := map[string]string{
		"call_scenario":      "通话场景",
		"music_scenario":     "音乐场景",
		"connection_status":  "连接状态",
		"system_events":      "系统事件",
		"call_events":        "通话事件",
		"tws_connected":      "TWS已连接",
		"tws_disconnected":   "TWS未连接",
	}
	
	if name, exists := nameMap[key]; exists {
		return name
	}
	return strings.Title(key)
}

// getFieldDisplayName 获取field的显示名称
func (a *App) getFieldDisplayName(key string) string {
	nameMap := map[string]string{
		"ic_model":           "IC型号",
		"vm_operation":       "VM操作",
		"pa_control":         "功放控制",
		"low_power_warn_time": "低电提醒时间",
		"active_click":       "通话中单击",
		"incoming_click":     "来电单击",
		"bluetooth_connected": "蓝牙已连接",
		"reset_mode":         "重置模式",
		"auto_power_on":      "自动开机",
	}
	
	if name, exists := nameMap[key]; exists {
		return name
	}
	return strings.Title(strings.ReplaceAll(key, "_", " "))
}

// isEnumValue 判断是否为枚举值
func (a *App) isEnumValue(value string) bool {
	enumPrefixes := []string{"APP_MSG_", "LED_", "FACTORY_", "DISP_"}
	for _, prefix := range enumPrefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

// getEnumOptions 获取枚举选项
func (a *App) getEnumOptions(value string) []models.ConfigOption {
	// 基于当前值推断可能的选项
	if strings.HasPrefix(value, "APP_MSG_") {
		return []models.ConfigOption{
			{Value: "APP_MSG_NULL", Label: "无操作"},
			{Value: "APP_MSG_CALL_ANSWER", Label: "接听"},
			{Value: "APP_MSG_CALL_HANGUP", Label: "挂断"},
			{Value: "APP_MSG_VOL_UP", Label: "音量+"},
			{Value: "APP_MSG_VOL_DOWN", Label: "音量-"},
			{Value: "APP_MSG_MUSIC_PP", Label: "播放/暂停"},
			{Value: "APP_MSG_MUSIC_NEXT", Label: "下一首"},
			{Value: "APP_MSG_MUSIC_PREV", Label: "上一首"},
			{Value: "APP_MSG_OPEN_SIRI", Label: "打开Siri"},
		}
	} else if strings.HasPrefix(value, "LED_") {
		return []models.ConfigOption{
			{Value: "LED_BLUE_ON", Label: "蓝灯常亮"},
			{Value: "LED_RED_ON", Label: "红灯常亮"},
			{Value: "LED_GREEN_ON", Label: "绿灯常亮"},
			{Value: "LED_BLUE_FAST", Label: "蓝灯快闪"},
			{Value: "LED_RED_SLOW", Label: "红灯慢闪"},
			{Value: "LED_OFF", Label: "关闭"},
		}
	}
	
	// 默认只提供当前值作为选项
	return []models.ConfigOption{
		{Value: value, Label: value},
	}
}