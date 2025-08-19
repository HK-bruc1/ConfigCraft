package main

import (
	"fmt"
	"os"
	"path/filepath"

	"dhf-config-manager/internal/config"
)

func main() {
	fmt.Println("DHF Configuration Manager - CLI Version")
	fmt.Println("=======================================")
	
	parser := config.NewParser()
	
	// Load schema
	schemaPath := filepath.Join("..", "assets", "schemas", "dhf-real-schema.yaml")
	if err := parser.LoadSchema(schemaPath); err != nil {
		fmt.Printf("Error loading schema: %v\n", err)
		return
	}
	
	schema := parser.GetSchema()
	fmt.Printf("Schema loaded: %s (v%s)\n", schema.DisplayName, schema.SchemaVersion)
	fmt.Printf("Configuration sections: %d\n", len(schema.Sections))
	
	for sectionKey, section := range schema.Sections {
		fmt.Printf("  - %s: %s\n", sectionKey, section.Name)
	}
	
	fmt.Println("\nUse GUI version (dhf-config-manager.exe) for full configuration.")
}