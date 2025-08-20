package components

import (
	"dhf-config-manager/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ConfigTree struct {
	container         fyne.CanvasObject
	scroll           *container.Scroll
	tree             *widget.Tree
	schema           *models.Schema
	selectionCallback func(string)
	isInitialized    bool // 标记是否已初始化
}

func NewConfigTree() *ConfigTree {
	ct := &ConfigTree{}
	
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			// 静态返回，避免动态计算导致的闪烁
			if ct.schema == nil {
				return []string{}
			}
			
			if uid == "" {
				// 返回固定的section列表
				var sections []string
				for key := range ct.schema.Sections {
					sections = append(sections, key)
				}
				return sections
			}
			
			// 返回固定的group列表
			if section, exists := ct.schema.Sections[uid]; exists {
				var groups []string
				for key := range section.Groups {
					groups = append(groups, uid+"."+key)
				}
				return groups
			}
			
			return []string{}
		},
		IsBranch: func(uid string) bool {
			if uid == "" {
				return true
			}
			
			if ct.schema == nil {
				return false
			}
			
			// 简单判断，避免复杂计算
			if section, exists := ct.schema.Sections[uid]; exists {
				return len(section.Groups) > 0
			}
			
			return false
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			// 统一创建标签，减少差异化处理
			label := widget.NewLabel("")
			if branch {
				// 分支节点稍微加粗
				label.TextStyle.Bold = true
			}
			return label
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			label := node.(*widget.Label)
			
			// 防止频繁更新相同内容
			var newText string
			
			if uid == "" {
				newText = "配置分组"
			} else if ct.schema == nil {
				newText = "Loading..."
			} else if section, exists := ct.schema.Sections[uid]; exists {
				newText = section.Name
			} else {
				// 查找group名称
				found := false
				for sectionKey, section := range ct.schema.Sections {
					for groupKey, group := range section.Groups {
						if uid == sectionKey+"."+groupKey {
							newText = group.Name
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				if !found {
					newText = uid
				}
			}
			
			// 只在文本确实需要更新时才设置，避免无意义的重绘
			if label.Text != newText {
				label.SetText(newText)
			}
		},
	}
	
	// 创建滚动容器
	scrollContainer := container.NewScroll(tree)
	scrollContainer.SetMinSize(fyne.NewSize(250, 300))
	
	ct.container = scrollContainer
	ct.scroll = scrollContainer
	ct.tree = tree
	
	// 优化选择处理，防止位置偏移
	tree.OnSelected = func(uid string) {
		if ct.selectionCallback != nil && uid != "" {
			// 延迟回调，避免选择期间的重绘干扰
			go func() {
				ct.selectionCallback(uid)
			}()
		}
	}
	
	return ct
}

func (ct *ConfigTree) Container() fyne.CanvasObject {
	return ct.container
}

func (ct *ConfigTree) LoadSchema(schema *models.Schema) {
	ct.schema = schema
	
	// 避免任何Refresh操作，让回调函数自然响应数据变化
	// 仅在第一次设置schema时标记为已初始化
	if !ct.isInitialized {
		ct.isInitialized = true
		// 只有在确实需要时才刷新
		if ct.tree != nil {
			ct.tree.Refresh()
		}
	}
}

func (ct *ConfigTree) SetSelectionCallback(callback func(string)) {
	ct.selectionCallback = callback
}

// ForceRefresh 强制刷新 - 仅在必要时使用
func (ct *ConfigTree) ForceRefresh() {
	if ct.tree != nil {
		ct.tree.Refresh()
	}
}