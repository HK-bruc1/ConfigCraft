# CLAUDE.md

This file provides technical guidance to Claude Code (claude.ai/code) when working with the ConfigCraft repository.

## Project Overview

**ConfigCraft** is a universal visual configuration management tool that transforms complex configuration files into user-friendly GUI interfaces using YAML schemas. Originally developed for DHF firmware configurations, it has evolved into a general-purpose configuration tool suitable for any schema-driven configuration management.

## Current Version: v0.3.5

- **Application Name**: ConfigCraft
- **Module Name**: configcraft
- **Executable**: configcraft.exe
- **Version Management**: `internal/version/version.go`

## Build and Run Commands

### Building the Application
```bash
# Build GUI version (Windows with TDM-GCC)
build\build.bat

# Manual build (requires TDM-GCC in PATH)
go build -ldflags "-s -w -H windowsgui" -o build\configcraft.exe main.go

# Development mode (run directly)
go run main.go
```

### CLI Version
```bash
cd cmd
go run cli.go
```

### Dependencies Management
```bash
# Install/update dependencies
go mod tidy

# Verify dependencies
go list -m all
```

## Architecture Overview

### Core Data Flow
**YAML Schema/Config** → **Parser** (`internal/config/`) → **Models** (`internal/models/`) → **UI Components** (`internal/ui/components/`) → **Configuration Output**

### Key Components

**Schema-Driven UI System**: Dynamic UI generation based on YAML schema files:
- `ConfigSection`: Top-level configuration groups
- `ConfigField`: Individual configuration parameters with validation
- `ConfigGroup`: Sub-groupings within sections

**MVC Architecture**:
- **Models** (`internal/models/types.go`): Core data structures
- **Controller** (`internal/config/parser.go`): File I/O, parsing, generation
- **View** (`internal/ui/`): Fyne-based GUI with custom components

**UI Components**:
- **Custom Tree** (`tree.go`): Solves Fyne Tree flickering issues
- **Dynamic Editor** (`editor.go`): Auto-generates forms based on schema
- **Smart Toolbar** (`toolbar.go`): Integrates zenity for native dialogs
- **Status Bar**: Shows current file path and version info

### Supported Field Types
- `select`: Dropdown with predefined options
- `combo`: Editable dropdown (v0.3.3+)
- `number`: Numeric input with validation
- `boolean`: Checkbox
- `text`: Text entry field

### Configuration Schema Structure
```yaml
sections:
  section_key:
    name: "Display Name"
    groups:
      group_key:
        name: "Group Name"
        fields:
          field_key:
            type: "select|combo|number|boolean|text"
            label: "UI Label"
            description: "Help text"
            tooltip: "Additional info"
            placeholder: "Input hint"
            options: [{value: "VALUE", label: "Display"}]
            default: value
```

## File Processing Pipeline

1. **File Detection**: Auto-detect schema vs configuration files
2. **Schema Loading**: Load/generate schema for UI rendering
3. **Configuration Parsing**: Parse YAML into `UserConfig.Values` map
4. **UI Generation**: Create form controls based on field types
5. **Output Generation**: Save YAML and generate corresponding config files

## Development Guidelines

### Version Management
- **Central Version**: Update only `internal/version/version.go`
- **Auto-sync**: Version appears in status bar and About dialog
- **Format**: "ConfigCraft v0.3.5"

### Adding New Features
1. Update schema structure if needed
2. Add field types in `editor.go` if required
3. Update parser logic for new formats
4. Test with various configuration files

### File Naming Conventions
- Input: `my_config.yaml`
- Output: `my_config.yaml` (updated) + `my_config.conf` (generated)

### Status Bar Behavior
- Shows relative path (max 2 levels: `../parent/file.yaml`)
- Schema mode: `Schema模式: schema.yaml`
- Config mode: `当前文件: config.yaml`
- Initial state: `请打开配置文件...`

## Known Issues and Solutions

### ✅ RESOLVED ISSUES

**Tree Widget Flickering (v0.3.1)**
- **Problem**: Fyne Tree widget caused flickering during expand/collapse
- **Solution**: Custom tree implementation using VBox containers
- **Location**: `internal/ui/components/tree.go`

**Configuration Group Random Positioning (v0.3.2)**
- **Problem**: Go map iteration randomness caused UI inconsistency
- **Solution**: Unified sorting logic with priority ordering
- **Locations**: `parser.go`, `app.go`, `tree.go`

**File Dialog Path Issues (v0.3.4)**
- **Problem**: Fyne dialogs couldn't start from current directory on Windows
- **Solution**: Integrated zenity library for native dialogs
- **Location**: `internal/ui/components/zenity_dialog.go`

**Configuration Editor Layout Chaos (v0.3.3)**
- **Problem**: Mixed configuration items without visual separation
- **Solution**: Individual card design for each config item
- **Location**: `internal/ui/components/editor.go`

**Chinese Font Display Issues (v0.2.1)**
- **Problem**: Chinese characters appeared as squares
- **Solution**: Force font via `FYNE_FONT` environment variable
- **Location**: `main.go` (sets `simhei.ttf`)

### 🚧 PENDING FEATURES
- Configuration file import (conf → YAML)
- Advanced validation and error checking
- Configuration templates system
- Batch configuration processing

## Technical Constraints

### GUI Requirements
- **Window Size**: 900x650 pixels (do not modify)
- **Font Support**: Chinese via `FYNE_FONT=C:\Windows\Fonts\simhei.ttf`
- **Native Dialogs**: Uses zenity library on Windows

### Build Requirements
- **Go Version**: 1.21+
- **CGO Compiler**: TDM-GCC 10.3.0 for Windows builds
- **Dependencies**: Fyne v2.4.3+, zenity v0.10.14+

### File Structure
```
configcraft/
├── internal/
│   ├── version/           # Version constants
│   ├── config/           # Parser and generator
│   ├── models/           # Data structures
│   └── ui/               # GUI components
├── assets/schemas/       # Schema files
├── build/               # Build artifacts
├── docs/                # Documentation
└── cmd/                 # CLI version
```

## Code Quality Standards

- **Error Handling**: Always return meaningful errors
- **Logging**: Use structured logging for debugging
- **UI Responsiveness**: Avoid blocking operations in GUI thread
- **Memory Management**: Close files and resources properly
- **Cross-platform**: Consider path separators and file permissions

## Testing Strategies

- **Schema Validation**: Test with various schema formats
- **File I/O**: Test read/write permissions and error cases  
- **UI Components**: Verify responsive behavior with different data sizes
- **Configuration Output**: Validate generated files match expected format

## Integration Points

- **File Dialog**: `ShowOpenDialog()` and `ShowSaveDialog()` in zenity_dialog.go
- **Configuration Parsing**: `LoadUserConfig()` and `SaveConfigWithConf()` in parser.go
- **UI Updates**: `updateStatusBar()` and tree refresh in app.go
- **Version Display**: Use `version.GetVersionString()` for all version references