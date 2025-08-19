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
	
	app := ui.NewApp()
	
	if err := app.Initialize(); err != nil {
		log.Fatal("Failed to initialize app:", err)
	}
	
	schemaPath := filepath.Join("assets", "schemas", "dhf-real-schema.yaml")
	if err := app.LoadSchema(schemaPath); err != nil {
		log.Printf("Warning: Failed to load schema: %v", err)
	}
	
	app.Run()
}