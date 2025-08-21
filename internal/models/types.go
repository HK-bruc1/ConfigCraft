package models

type ConfigSection struct {
	Name   string                 `yaml:"name"`
	Icon   string                 `yaml:"icon"`
	Fields map[string]ConfigField `yaml:"fields"`
	Groups map[string]ConfigGroup `yaml:"groups"`
}

type ConfigGroup struct {
	Name   string                 `yaml:"name"`
	Fields map[string]ConfigField `yaml:"fields"`
}

type ConfigField struct {
	Type        string                 `yaml:"type"`
	Label       string                 `yaml:"label"`
	Description string                 `yaml:"description,omitempty"`  // 字段描述信息
	Tooltip     string                 `yaml:"tooltip,omitempty"`      // 鼠标悬停提示
	Placeholder string                 `yaml:"placeholder,omitempty"`  // 输入框占位符
	Options     []ConfigOption         `yaml:"options,omitempty"`
	Default     interface{}            `yaml:"default,omitempty"`
	Required    bool                   `yaml:"required,omitempty"`
	Min         *int                   `yaml:"min,omitempty"`
	Max         *int                   `yaml:"max,omitempty"`
}

type ConfigOption struct {
	Value interface{} `yaml:"value"`
	Label string      `yaml:"label"`
}

type Schema struct {
	SchemaVersion string                    `yaml:"schema_version"`
	DisplayName   string                    `yaml:"display_name"`
	Sections      map[string]ConfigSection `yaml:"sections"`
}

type UserConfig struct {
	Values map[string]interface{} `json:"values"`
}