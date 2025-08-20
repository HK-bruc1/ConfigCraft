package main

import (
	"log"
	"os"
	"path/filepath"

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
	
	// 加载中文schema（中文字体问题已解决）
	schemaPath := filepath.Join("assets", "schemas", "dhf-schema.yaml")
	if err := app.LoadSchema(schemaPath); err != nil {
		log.Printf("Warning: Failed to load schema: %v", err)
	} else {
		log.Printf("Successfully loaded Chinese schema: %s", schemaPath)
	}
	
	app.Run()
}