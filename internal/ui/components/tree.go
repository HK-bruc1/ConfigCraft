package components

import (
	"configcraft/internal/models"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TreeNode 自定义树节点结构
type TreeNode struct {
	id          string
	name        string
	isSection   bool
	isExpanded  bool
	children    []*TreeNode
	parent      *TreeNode
	widget      fyne.CanvasObject
}

// ConfigTree 自定义树形控件，解决Fyne Tree的闪烁问题
type ConfigTree struct {
	container         fyne.CanvasObject
	scroll           *container.Scroll
	vbox             *fyne.Container
	schema           *models.Schema
	selectionCallback func(string)
	nodes            map[string]*TreeNode
	selectedNode     *TreeNode
}

func NewConfigTree() *ConfigTree {
	ct := &ConfigTree{
		nodes: make(map[string]*TreeNode),
	}
	
	// 创建主容器
	vbox := container.NewVBox()
	scrollContainer := container.NewScroll(vbox)
	scrollContainer.SetMinSize(fyne.NewSize(250, 300))
	
	ct.container = scrollContainer
	ct.scroll = scrollContainer
	ct.vbox = vbox
	
	return ct
}

func (ct *ConfigTree) Container() fyne.CanvasObject {
	return ct.container
}

func (ct *ConfigTree) LoadSchema(schema *models.Schema) {
	ct.schema = schema
	ct.rebuildTree()
}

func (ct *ConfigTree) SetSelectionCallback(callback func(string)) {
	ct.selectionCallback = callback
}

// rebuildTree 重建整个树结构
func (ct *ConfigTree) rebuildTree() {
	if ct.schema == nil {
		return
	}
	
	// 清空现有节点
	ct.vbox.RemoveAll()
	ct.nodes = make(map[string]*TreeNode)
	
	// 创建根节点
	rootNode := &TreeNode{
		id:         "root",
		name:       "配置分组",
		isSection:  false,
		isExpanded: true,
		children:   make([]*TreeNode, 0),
	}
	ct.nodes["root"] = rootNode
	
	// 为每个section创建节点 - 按固定顺序遍历以确保界面一致性
	sectionKeys := make([]string, 0, len(ct.schema.Sections))
	for sectionKey := range ct.schema.Sections {
		sectionKeys = append(sectionKeys, sectionKey)
	}
	// 使用与app.go相同的排序逻辑
	sectionOrder := map[string]int{
		"basic":      1,
		"key_actions": 2,
		"led_config": 3,
		"factory":    4,
		"advanced":   5,
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
	
	for _, sectionKey := range sectionKeys {
		section := ct.schema.Sections[sectionKey]
		sectionNode := &TreeNode{
			id:         sectionKey,
			name:       section.Name,
			isSection:  true,
			isExpanded: false,
			children:   make([]*TreeNode, 0),
			parent:     rootNode,
		}
		ct.nodes[sectionKey] = sectionNode
		rootNode.children = append(rootNode.children, sectionNode)
		
		// 为每个group创建子节点 - 按固定顺序遍历
		groupKeys := make([]string, 0, len(section.Groups))
		for groupKey := range section.Groups {
			groupKeys = append(groupKeys, groupKey)
		}
		sort.Strings(groupKeys) // 按字典序排序
		
		for _, groupKey := range groupKeys {
			group := section.Groups[groupKey]
			groupID := sectionKey + "." + groupKey
			groupNode := &TreeNode{
				id:         groupID,
				name:       group.Name,
				isSection:  false,
				isExpanded: false,
				children:   make([]*TreeNode, 0),
				parent:     sectionNode,
			}
			ct.nodes[groupID] = groupNode
			sectionNode.children = append(sectionNode.children, groupNode)
		}
	}
	
	// 渲染树结构
	ct.renderTree()
}

// renderTree 渲染整个树到界面
func (ct *ConfigTree) renderTree() {
	ct.vbox.RemoveAll()
	
	if rootNode, exists := ct.nodes["root"]; exists {
		ct.renderNodeChildren(rootNode, 0)
	}
	
	ct.vbox.Refresh()
}

// renderNodeChildren 渲染节点的子节点
func (ct *ConfigTree) renderNodeChildren(parentNode *TreeNode, depth int) {
	for _, child := range parentNode.children {
		ct.renderNode(child, depth)
		
		// 如果节点展开，递归渲染子节点
		if child.isExpanded && len(child.children) > 0 {
			ct.renderNodeChildren(child, depth+1)
		}
	}
}

// renderNode 渲染单个节点
func (ct *ConfigTree) renderNode(node *TreeNode, depth int) {
	// 创建水平容器
	nodeContainer := container.NewHBox()
	
	// 添加紧凑的缩进空间
	if depth > 0 {
		indentSpace := widget.NewLabel(strings.Repeat("  ", depth)) // 改为2个空格
		nodeContainer.Add(indentSpace)
	}
	
	// 创建简约箭头图标
	if len(node.children) > 0 {
		// 使用简洁的箭头
		var expandIcon string
		if node.isExpanded {
			expandIcon = "▼" // 下箭头，表示已展开
		} else {
			expandIcon = "▶" // 右箭头，表示收缩
		}
		
		// 创建可点击的图标按钮
		expandButton := widget.NewButton(expandIcon, func() {
			ct.toggleNode(node)
		})
		expandButton.Importance = widget.LowImportance
		
		// 设置更小的尺寸
		expandButton.Resize(fyne.NewSize(20, 20))
		nodeContainer.Add(expandButton)
	} else {
		// 叶子节点留空，保持对齐
		spacer := widget.NewLabel("")
		spacer.Resize(fyne.NewSize(20, 20))
		nodeContainer.Add(spacer)
	}
	
	// 简洁的节点文本，不添加额外图标
	nodeText := node.name
	
	// 创建节点文本按钮（用于选择）
	selectButton := widget.NewButton(nodeText, func() {
		ct.selectNode(node)
	})
	
	// 根据节点类型和选中状态设置样式
	if ct.selectedNode == node {
		// 选中节点高亮显示
		selectButton.Importance = widget.HighImportance
	} else if node.isSection {
		// Section节点使用中等重要性
		selectButton.Importance = widget.MediumImportance
	} else {
		// Group节点使用低重要性
		selectButton.Importance = widget.LowImportance
	}
	
	nodeContainer.Add(selectButton)
	
	// 保存widget引用以备更新
	node.widget = nodeContainer
	
	ct.vbox.Add(nodeContainer)
}

// toggleNode 展开/收缩节点
func (ct *ConfigTree) toggleNode(node *TreeNode) {
	node.isExpanded = !node.isExpanded
	
	// 重新渲染整个树，但这次没有Fyne Tree的刷新问题
	ct.renderTree()
}

// selectNode 选择节点
func (ct *ConfigTree) selectNode(node *TreeNode) {
	ct.selectedNode = node
	
	// 重新渲染以更新选择状态
	ct.renderTree()
	
	// 触发回调
	if ct.selectionCallback != nil {
		ct.selectionCallback(node.id)
	}
}

// ForceRefresh 强制刷新 - 重建树结构
func (ct *ConfigTree) ForceRefresh() {
	ct.rebuildTree()
}