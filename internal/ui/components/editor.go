package components

import (
	"configcraft/internal/models"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type ConfigEditor struct {
	container  fyne.CanvasObject
	content    *fyne.Container
	schema     *models.Schema
	userConfig *models.UserConfig
	window     fyne.Window // 添加窗口引用以支持弹窗
}

func NewConfigEditor() *ConfigEditor {
	content := container.NewVBox()
	
	// 添加简单的欢迎提示
	welcomeCard := widget.NewCard("Welcome", "Select a configuration group from the left panel to begin editing.", container.NewVBox())
	content.Add(welcomeCard)
	
	// 创建滚动容器
	scrollContainer := container.NewScroll(content)
	scrollContainer.SetMinSize(fyne.NewSize(400, 300))
	
	// 简化布局，移除多余的标题
	container := container.NewPadded(scrollContainer)
	
	return &ConfigEditor{
		container: container,
		content:   content,
	}
}

func (ce *ConfigEditor) Container() fyne.CanvasObject {
	return ce.container
}

func (ce *ConfigEditor) SetSchema(schema *models.Schema) {
	ce.schema = schema
}

func (ce *ConfigEditor) SetConfig(config *models.UserConfig) {
	ce.userConfig = config
	// 触发当前显示内容的刷新
	ce.content.Refresh()
}

func (ce *ConfigEditor) SetWindow(window fyne.Window) {
	ce.window = window
}

func (ce *ConfigEditor) ShowSection(sectionID string) {
	ce.content.Objects = nil
	
	if ce.schema == nil {
		ce.content.Add(widget.NewLabel("No schema loaded"))
		ce.content.Refresh()
		return
	}
	
	parts := strings.Split(sectionID, ".")
	
	if len(parts) == 1 {
		ce.showSectionFields(sectionID)
	} else if len(parts) == 2 {
		ce.showGroupFields(parts[0], parts[1])
	}
	
	ce.content.Refresh()
}

func (ce *ConfigEditor) showSectionFields(sectionID string) {
	section, exists := ce.schema.Sections[sectionID]
	if !exists {
		errorCard := widget.NewCard("Error", "Section not found: "+sectionID, container.NewVBox())
		ce.content.Add(errorCard)
		return
	}
	
	// 创建现代化的分组标题卡片
	headerCard := widget.NewCard(section.Name, "Configure the settings below", container.NewVBox())
	ce.content.Add(headerCard)
	
	// 重新设计字段布局：每个字段独立成卡片
	fieldsContainer := container.NewVBox()
	
	// 按字段键排序以确保一致的显示顺序
	fieldKeys := make([]string, 0, len(section.Fields))
	for fieldKey := range section.Fields {
		fieldKeys = append(fieldKeys, fieldKey)
	}
	sort.Strings(fieldKeys)
	
	for _, fieldKey := range fieldKeys {
		field := section.Fields[fieldKey]
		fieldWidget := ce.createFieldWidget(sectionID+"."+fieldKey, field)
		
		// 每个字段都有自己的卡片，确保明确的视觉分离
		fieldCard := widget.NewCard("", "", fieldWidget)
		fieldsContainer.Add(fieldCard)
		
		// 添加间距
		fieldsContainer.Add(widget.NewSeparator())
	}
	
	ce.content.Add(fieldsContainer)
}

func (ce *ConfigEditor) showGroupFields(sectionID, groupID string) {
	section, exists := ce.schema.Sections[sectionID]
	if !exists {
		errorCard := widget.NewCard("Error", "Section not found: "+sectionID, container.NewVBox())
		ce.content.Add(errorCard)
		return
	}
	
	group, exists := section.Groups[groupID]
	if !exists {
		errorCard := widget.NewCard("Error", "Group not found: "+groupID, container.NewVBox())
		ce.content.Add(errorCard)
		return
	}
	
	// 创建现代化的组标题卡片
	headerCard := widget.NewCard(group.Name, "Configure the group settings below", container.NewVBox())
	ce.content.Add(headerCard)
	
	// 重新设计组字段布局：每个字段独立成卡片
	fieldsContainer := container.NewVBox()
	
	// 按字段键排序以确保一致的显示顺序
	fieldKeys := make([]string, 0, len(group.Fields))
	for fieldKey := range group.Fields {
		fieldKeys = append(fieldKeys, fieldKey)
	}
	sort.Strings(fieldKeys)
	
	for _, fieldKey := range fieldKeys {
		field := group.Fields[fieldKey]
		fieldWidget := ce.createFieldWidget(sectionID+"."+groupID+"."+fieldKey, field)
		
		// 每个字段都有自己的卡片，确保明确的视觉分离
		fieldCard := widget.NewCard("", "", fieldWidget)
		fieldsContainer.Add(fieldCard)
		
		// 添加间距
		fieldsContainer.Add(widget.NewSeparator())
	}
	
	ce.content.Add(fieldsContainer)
}

func (ce *ConfigEditor) createFieldWidget(fieldPath string, field models.ConfigField) fyne.CanvasObject {
	// 创建规整的字段布局容器
	fieldContainer := container.NewVBox()
	
	// === 第一行：标题行（标签 + 帮助按钮） ===
	headerRow := container.NewBorder(nil, nil, nil, nil)
	
	// 主标签
	titleLabel := widget.NewLabelWithStyle(field.Label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleLabel.TextStyle.Monospace = false
	
	headerContent := container.NewHBox(titleLabel)
	
	// 如果有tooltip，添加统一样式的帮助按钮
	if field.Tooltip != "" {
		helpBtn := widget.NewButtonWithIcon("", nil, func() {
			ce.showHelpDialog(field.Label, field.Tooltip)
		})
		helpBtn.SetText("💡") // 使用灯泡图标
		helpBtn.Importance = widget.LowImportance
		headerContent.Add(helpBtn)
	}
	
	headerRow.Objects = []fyne.CanvasObject{headerContent}
	fieldContainer.Add(headerRow)
	
	// === 第二行：描述信息（如果有） ===
	if field.Description != "" {
		descText := widget.NewLabelWithStyle("📝 "+field.Description, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
		descText.Wrapping = fyne.TextWrapWord
		fieldContainer.Add(descText)
		
		// 添加小间距
		fieldContainer.Add(widget.NewSeparator())
	}
	
	// === 第三行：控件区域 ===
	var controlWidget fyne.CanvasObject
	switch field.Type {
	case "select":
		controlWidget = ce.createSelectWidget(fieldPath, field)
	case "combo":
		controlWidget = ce.createComboWidget(fieldPath, field)
	case "text":
		controlWidget = ce.createTextWidget(fieldPath, field)
	case "number":
		controlWidget = ce.createNumberWidget(fieldPath, field)
	case "boolean":
		controlWidget = ce.createBooleanWidget(fieldPath, field)
	default:
		controlWidget = ce.createTextWidget(fieldPath, field)
	}
	
	fieldContainer.Add(controlWidget)
	
	// 使用边框容器添加统一的内边距
	return container.NewPadded(fieldContainer)
}

// 统一的帮助对话框显示方法
func (ce *ConfigEditor) showHelpDialog(title, content string) {
	// 创建格式化的帮助内容
	helpContent := container.NewVBox(
		widget.NewLabelWithStyle("字段说明", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewRichTextFromMarkdown(content),
	)
	
	// 创建统一样式的帮助对话框
	helpDialog := dialog.NewCustom("帮助信息", "关闭", helpContent, ce.window)
	helpDialog.Resize(fyne.NewSize(450, 250))
	helpDialog.Show()
}

func (ce *ConfigEditor) createSelectWidget(fieldPath string, field models.ConfigField) fyne.CanvasObject {
	var options []string
	var values []interface{}
	
	for _, option := range field.Options {
		options = append(options, option.Label)
		values = append(values, option.Value)
	}
	
	selectWidget := widget.NewSelect(options, func(selected string) {
		for i, option := range field.Options {
			if option.Label == selected {
				ce.setValue(fieldPath, values[i])
				break
			}
		}
	})
	
	if currentValue := ce.getValue(fieldPath); currentValue != nil {
		for i, value := range values {
			if value == currentValue {
				selectWidget.SetSelected(options[i])
				break
			}
		}
	} else if field.Default != nil {
		for i, value := range values {
			if value == field.Default {
				selectWidget.SetSelected(options[i])
				break
			}
		}
	}
	
	return selectWidget
}

func (ce *ConfigEditor) createTextWidget(fieldPath string, field models.ConfigField) fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.OnChanged = func(text string) {
		ce.setValue(fieldPath, text)
	}
	
	if currentValue := ce.getValue(fieldPath); currentValue != nil {
		if str, ok := currentValue.(string); ok {
			entry.SetText(str)
		}
	} else if field.Default != nil {
		if str, ok := field.Default.(string); ok {
			entry.SetText(str)
		}
	}
	
	return entry
}

func (ce *ConfigEditor) createNumberWidget(fieldPath string, field models.ConfigField) fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.OnChanged = func(text string) {
		if val, err := strconv.Atoi(text); err == nil {
			ce.setValue(fieldPath, val)
		}
	}
	
	if currentValue := ce.getValue(fieldPath); currentValue != nil {
		if num, ok := currentValue.(int); ok {
			entry.SetText(strconv.Itoa(num))
		}
	} else if field.Default != nil {
		if num, ok := field.Default.(int); ok {
			entry.SetText(strconv.Itoa(num))
		}
	}
	
	return entry
}

func (ce *ConfigEditor) createBooleanWidget(fieldPath string, field models.ConfigField) fyne.CanvasObject {
	check := widget.NewCheck("", func(checked bool) {
		ce.setValue(fieldPath, checked)
	})
	
	if currentValue := ce.getValue(fieldPath); currentValue != nil {
		if b, ok := currentValue.(bool); ok {
			check.SetChecked(b)
		}
	} else if field.Default != nil {
		if b, ok := field.Default.(bool); ok {
			check.SetChecked(b)
		}
	}
	
	return check
}

// 新增：创建可编辑下拉框
func (ce *ConfigEditor) createComboWidget(fieldPath string, field models.ConfigField) fyne.CanvasObject {
	// 准备选项
	var options []string
	var values []interface{}
	
	for _, option := range field.Options {
		options = append(options, option.Label)
		values = append(values, option.Value)
	}
	
	// 创建一个容器，包含下拉框和文本输入框
	entry := widget.NewEntry()
	if field.Placeholder != "" {
		entry.PlaceHolder = field.Placeholder
	}
	
	// 创建选择框用于快速选择预设值
	var selectWidget *widget.Select
	selectWidget = widget.NewSelect(append([]string{"选择预设值..."}, options...), func(selected string) {
		if selected == "选择预设值..." {
			return
		}
		
		// 找到对应的值并设置到输入框
		for i, option := range field.Options {
			if option.Label == selected {
				if str, ok := values[i].(string); ok {
					entry.SetText(str)
				} else {
					entry.SetText(fmt.Sprintf("%v", values[i]))
				}
				ce.setValue(fieldPath, values[i])
				break
			}
		}
		// 重置选择框显示
		selectWidget.SetSelected("选择预设值...")
	})
	
	// 文本输入框变化时更新配置值
	entry.OnChanged = func(text string) {
		// 尝试匹配预设值
		for _, option := range field.Options {
			if fmt.Sprintf("%v", option.Value) == text {
				ce.setValue(fieldPath, option.Value)
				return
			}
		}
		// 如果不是预设值，直接使用文本值
		ce.setValue(fieldPath, text)
	}
	
	// 设置初始值
	if currentValue := ce.getValue(fieldPath); currentValue != nil {
		entry.SetText(fmt.Sprintf("%v", currentValue))
	} else if field.Default != nil {
		entry.SetText(fmt.Sprintf("%v", field.Default))
		ce.setValue(fieldPath, field.Default)
	}
	
	// 创建简洁的布局：上下结构，没有多余标签
	comboContainer := container.NewVBox(
		selectWidget,
		widget.NewSeparator(), // 添加分隔线区分两个控件
		entry,
	)
	
	return comboContainer
}

func (ce *ConfigEditor) getValue(fieldPath string) interface{} {
	if ce.userConfig != nil {
		return ce.userConfig.Values[fieldPath]
	}
	return nil
}

func (ce *ConfigEditor) setValue(fieldPath string, value interface{}) {
	if ce.userConfig == nil {
		ce.userConfig = &models.UserConfig{Values: make(map[string]interface{})}
	}
	ce.userConfig.Values[fieldPath] = value
}