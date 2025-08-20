package ui

import (
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type chineseTheme struct {
	fyne.Theme
}

// NewChineseTheme creates a new theme with Chinese font support
func NewChineseTheme() fyne.Theme {
	return &chineseTheme{Theme: theme.DefaultTheme()}
}

func (c *chineseTheme) Font(style fyne.TextStyle) fyne.Resource {
	// 返回默认字体，配合环境变量FYNE_FONT使用
	return theme.DefaultTheme().Font(style)
}

func (c *chineseTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (c *chineseTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (c *chineseTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14 // 适合中文显示的字体大小
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameHeadingText:
		return 18
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// 简化版本：使用Fyne内置的字体回退机制
type simpleChineseTheme struct {
	fyne.Theme
}

func NewSimpleChineseTheme() fyne.Theme {
	return &simpleChineseTheme{Theme: theme.DefaultTheme()}
}

func (s *simpleChineseTheme) Font(style fyne.TextStyle) fyne.Resource {
	// 返回默认字体，配合环境变量FYNE_FONT使用
	return theme.DefaultTheme().Font(style)
}

func (s *simpleChineseTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14 // 稍大的字体便于中文显示
	default:
		return theme.DefaultTheme().Size(name)
	}
}