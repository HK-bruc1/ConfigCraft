package config

import (
	"fmt"
	"os"

	"dhf-config-manager/internal/models"
	"gopkg.in/yaml.v3"
)

type Parser struct {
	schema *models.Schema
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) LoadSchema(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	var schema models.Schema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	p.schema = &schema
	return nil
}

func (p *Parser) GetSchema() *models.Schema {
	return p.schema
}

func (p *Parser) LoadUserConfig(filePath string) (*models.UserConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &models.UserConfig{Values: make(map[string]interface{})}, nil
	}

	var config models.UserConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse user config: %w", err)
	}

	if config.Values == nil {
		config.Values = make(map[string]interface{})
	}

	return &config, nil
}

func (p *Parser) SaveUserConfig(config *models.UserConfig, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}