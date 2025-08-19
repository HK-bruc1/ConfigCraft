package ui

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type chineseTheme struct{}

func (c *chineseTheme) Font(style fyne.TextStyle) fyne.Resource {
	// 使用系统默认字体，通常支持中文
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
		return 16 // 更大的字体便于中文显示
	case theme.SizeNameCaptionText:
		return 12
	default:
		return theme.DefaultTheme().Size(name)
	}
}