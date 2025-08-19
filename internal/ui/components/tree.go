package components

import (
	"dhf-config-manager/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ConfigTree struct {
	container         *fyne.Container
	tree             *widget.Tree
	schema           *models.Schema
	selectionCallback func(string)
}

func NewConfigTree() *ConfigTree {
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return []string{}
		},
		IsBranch: func(uid string) bool {
			return uid == ""
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			label := node.(*widget.Label)
			label.SetText("Loading...")
		},
	}
	
	// 创建标题标签
	titleLabel := widget.NewLabel("Configuration")
	titleLabel.Alignment = fyne.TextAlignCenter
	
	// 创建带标题的容器
	container := container.NewBorder(
		container.NewVBox(
			titleLabel,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		container.NewScroll(tree), // 添加滚动支持
	)
	
	ct := &ConfigTree{
		container: container,
		tree:      tree,
	}
	
	tree.OnSelected = func(uid string) {
		if ct.selectionCallback != nil {
			ct.selectionCallback(uid)
		}
	}
	
	return ct
}

func (ct *ConfigTree) Container() *fyne.Container {
	return ct.container
}

func (ct *ConfigTree) LoadSchema(schema *models.Schema) {
	ct.schema = schema
	
	ct.tree.ChildUIDs = func(uid string) []string {
		if uid == "" {
			var sections []string
			for key := range schema.Sections {
				sections = append(sections, key)
			}
			return sections
		}
		
		if section, exists := schema.Sections[uid]; exists {
			var groups []string
			for key := range section.Groups {
				groups = append(groups, uid+"."+key)
			}
			return groups
		}
		
		return []string{}
	}
	
	ct.tree.IsBranch = func(uid string) bool {
		if uid == "" {
			return true
		}
		
		if section, exists := schema.Sections[uid]; exists {
			return len(section.Groups) > 0
		}
		
		return false
	}
	
	ct.tree.UpdateNode = func(uid string, branch bool, node fyne.CanvasObject) {
		label := node.(*widget.Label)
		
		if uid == "" {
			label.SetText("Configuration")
			return
		}
		
		if section, exists := schema.Sections[uid]; exists {
			label.SetText(section.Name)
			return
		}
		
		for sectionKey, section := range schema.Sections {
			for groupKey, group := range section.Groups {
				if uid == sectionKey+"."+groupKey {
					label.SetText(group.Name)
					return
				}
			}
		}
		
		label.SetText(uid)
	}
	
	ct.tree.Refresh()
}

func (ct *ConfigTree) SetSelectionCallback(callback func(string)) {
	ct.selectionCallback = callback
}