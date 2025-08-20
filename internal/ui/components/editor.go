package components

import (
	"dhf-config-manager/internal/models"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ConfigEditor struct {
	container  fyne.CanvasObject
	content    *fyne.Container
	schema     *models.Schema
	userConfig *models.UserConfig
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
	
	// 将字段分组到卡片中，每个卡片最多3-4个字段
	fieldsContainer := container.NewVBox()
	fieldCount := 0
	currentCard := container.NewVBox()
	
	for fieldKey, field := range section.Fields {
		fieldWidget := ce.createFieldWidget(sectionID+"."+fieldKey, field)
		currentCard.Add(fieldWidget)
		fieldCount++
		
		// 每3个字段创建一个新卡片
		if fieldCount%3 == 0 {
			card := widget.NewCard("Settings", "", currentCard)
			fieldsContainer.Add(card)
			currentCard = container.NewVBox()
		}
	}
	
	// 添加剩余字段
	if len(currentCard.Objects) > 0 {
		card := widget.NewCard("Settings", "", currentCard)
		fieldsContainer.Add(card)
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
	
	// 将组内字段分组到卡片中
	fieldsContainer := container.NewVBox()
	fieldCount := 0
	currentCard := container.NewVBox()
	
	for fieldKey, field := range group.Fields {
		fieldWidget := ce.createFieldWidget(sectionID+"."+groupID+"."+fieldKey, field)
		currentCard.Add(fieldWidget)
		fieldCount++
		
		// 每4个字段创建一个新卡片
		if fieldCount%4 == 0 {
			card := widget.NewCard("Group Settings", "", currentCard)
			fieldsContainer.Add(card)
			currentCard = container.NewVBox()
		}
	}
	
	// 添加剩余字段
	if len(currentCard.Objects) > 0 {
		card := widget.NewCard("Group Settings", "", currentCard)
		fieldsContainer.Add(card)
	}
	
	ce.content.Add(fieldsContainer)
}

func (ce *ConfigEditor) createFieldWidget(fieldPath string, field models.ConfigField) fyne.CanvasObject {
	var fieldWidget fyne.CanvasObject
	
	switch field.Type {
	case "select":
		fieldWidget = ce.createSelectWidget(fieldPath, field)
	case "text":
		fieldWidget = ce.createTextWidget(fieldPath, field)
	case "number":
		fieldWidget = ce.createNumberWidget(fieldPath, field)
	case "boolean":
		fieldWidget = ce.createBooleanWidget(fieldPath, field)
	default:
		fieldWidget = ce.createTextWidget(fieldPath, field)
	}
	
	// 创建现代化的字段标签，带有描述性样式
	label := widget.NewLabelWithStyle(field.Label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
	// 使用垂直布局，标签在上，控件在下
	fieldContainer := container.NewVBox(
		label,
		fieldWidget,
	)
	
	// 添加内边距使布局更美观
	return container.NewPadded(fieldContainer)
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