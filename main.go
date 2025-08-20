package main

import (
	"log"
	"os"

	"dhf-config-manager/internal/ui"
)

func main() {
	// 设置UTF-8编码环境
	os.Setenv("LANG", "zh_CN.UTF-8")
	os.Setenv("LC_ALL", "zh_CN.UTF-8")
	
	// 强制Fyne使用黑体字体以支持中文显示（使用TTF而不是TTC格式）
	os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\simhei.ttf")
	
	log.Println("Font path set to:", os.Getenv("FYNE_FONT"))
	
	app := ui.NewApp()
	
	if err := app.Initialize(); err != nil {
		log.Fatal("Failed to initialize app:", err)
	}
	
	log.Println("DHF Configuration Manager started. Please load a configuration file using the toolbar.")
	
	app.Run()
}