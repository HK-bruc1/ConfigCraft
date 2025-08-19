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
	container  *fyne.Container
	content    *fyne.Container
	schema     *models.Schema
	userConfig *models.UserConfig
}

func NewConfigEditor() *ConfigEditor {
	content := container.NewVBox()
	
	// 创建标题标签
	titleLabel := widget.NewLabel("Configuration Editor")
	titleLabel.Alignment = fyne.TextAlignCenter
	
	// 创建滚动容器
	scrollContainer := container.NewScroll(content)
	scrollContainer.SetMinSize(fyne.NewSize(300, 200)) // 更紧凑的最小尺寸
	
	container := container.NewBorder(
		container.NewVBox(
			titleLabel,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		scrollContainer,
	)
	
	return &ConfigEditor{
		container: container,
		content:   content,
	}
}

func (ce *ConfigEditor) Container() *fyne.Container {
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
		ce.content.Add(widget.NewLabel("Section not found: " + sectionID))
		return
	}
	
	ce.content.Add(widget.NewCard(section.Name, "", container.NewVBox()))
	
	for fieldKey, field := range section.Fields {
		fieldWidget := ce.createFieldWidget(sectionID+"."+fieldKey, field)
		ce.content.Add(fieldWidget)
	}
}

func (ce *ConfigEditor) showGroupFields(sectionID, groupID string) {
	section, exists := ce.schema.Sections[sectionID]
	if !exists {
		ce.content.Add(widget.NewLabel("Section not found: " + sectionID))
		return
	}
	
	group, exists := section.Groups[groupID]
	if !exists {
		ce.content.Add(widget.NewLabel("Group not found: " + groupID))
		return
	}
	
	ce.content.Add(widget.NewCard(group.Name, "", container.NewVBox()))
	
	for fieldKey, field := range group.Fields {
		fieldWidget := ce.createFieldWidget(sectionID+"."+groupID+"."+fieldKey, field)
		ce.content.Add(fieldWidget)
	}
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
	
	label := widget.NewLabel(field.Label)
	return container.NewBorder(
		nil, nil,
		label,
		nil,
		fieldWidget,
	)
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