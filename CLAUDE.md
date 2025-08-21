# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DHF Configuration Manager is a visual configuration tool for DHF AC710N-V300P03 SDK that converts complex `.conf` files into user-friendly GUI interfaces using YAML schemas. The application transforms complex firmware configuration syntax into intuitive form controls.

## Build and Run Commands

### Building the Application
```bash
# Build GUI version (Windows with TDM-GCC)
build\build.bat

# Manual build (requires TDM-GCC in PATH and proxy settings)
go build -ldflags "-s -w -H windowsgui" -o build\dhf-config-manager.exe main.go

# Development mode (run directly)
go run main.go
```

### CLI Version
```bash
cd cmd
go run cli.go
```

### Dependencies
```bash
# Install/update dependencies
go mod tidy

# Verify dependencies
go list -m all
```

## Architecture Overview

### Core Data Flow
1. **YAML Schema** (`assets/schemas/`) → **Parser** (`internal/config/`) → **Models** (`internal/models/`) → **UI Components** (`internal/ui/components/`) → **DHF conf output**

### Key Components

**Schema-Driven UI System**: The application uses YAML schema files to dynamically generate UI components. Each schema defines:
- `ConfigSection`: Top-level configuration groups (basic, key_actions, led_config, etc.)  
- `ConfigField`: Individual configuration parameters with type, options, validation
- `ConfigGroup`: Sub-groupings within sections (e.g., call_scenario, music_scenario)

**MVC Architecture**:
- **Models** (`internal/models/types.go`): Core data structures (Schema, ConfigField, UserConfig)
- **Controller** (`internal/config/parser.go`): YAML parsing, conf generation, file I/O
- **View** (`internal/ui/`): Fyne-based GUI with tree navigation and form editors

**Dynamic UI Generation**: UI components are created at runtime based on field types:
- `select` → Dropdown with predefined options
- `number` → Numeric input with min/max validation  
- `boolean` → Checkbox
- `text` → Text entry field

### Configuration Schema Structure
The schema format maps directly to DHF conf file syntax:
```yaml
sections:
  section_key:
    name: "Display Name"
    fields:
      field_key:
        type: "select|number|boolean|text"
        label: "UI Label"
        options: [{value: "CONF_VALUE", label: "Display Text"}]
        default: value
```

### File Processing Pipeline
1. Load schema from `assets/schemas/dhf-real-schema.yaml`
2. Parse user input into `UserConfig.Values` map (key paths like "section.group.field")
3. Generate DHF conf format with proper macro naming (e.g., `_SECTION_FIELD=value`)

## Development Notes

### Adding New Configuration Options
1. Add field definition to appropriate schema file (`assets/schemas/`)
2. UI components auto-generate based on field type
3. conf output automatically includes new fields with proper naming conventions

### Known Issues to Address
- ✅ ~~Chinese character display in GUI~~ **FIXED** - Using FYNE_FONT environment variable with simhei.ttf
- ✅ ~~Tree widget flickering and position drift~~ **FIXED** - Custom tree implementation replaces problematic Fyne Tree
- ✅ ~~Configuration groups random positioning~~ **FIXED** - Unified sorting logic with priority ordering
- ✅ ~~Configuration editor layout chaos~~ **FIXED** - Individual card design with clear visual hierarchy
- Missing conf-to-YAML import functionality  
- No configuration validation/error checking

### Configuration Editor UI Revolution (v0.3.3)
- **Problem**: Right-side configuration editor had chaotic layout with unclear item boundaries
- **Root Cause**: Multiple configuration items mixed in single cards without proper visual separation
- **Solution**: Individual card design for each configuration item with unified visual hierarchy
- **Benefits**: Clear visual separation, consistent interaction patterns, professional appearance

### Enhanced Schema Support (v0.3.3)
- **New Field Type**: `combo` - Editable dropdown combining preset selection with manual input
- **New Attributes**: `description`, `tooltip`, `placeholder` for comprehensive user guidance
- **Smart Detection**: Automatic recognition of schema files vs configuration files
- **Backwards Compatible**: Existing configuration files continue to work with dynamic schema generation

### Configuration Group Order Fix (v0.3.2)
- **Problem**: Configuration groups appear in random positions each time the same YAML file is loaded
- **Root Cause**: Go map iteration is random, causing inconsistent UI rendering across three key locations
- **Solution**: Unified sorting logic in parser.go, app.go, and tree.go with priority-based ordering
- **Benefits**: Stable, predictable group positioning with logical priority sequence (basic → key_actions → led_config → factory → advanced)

### Tree Widget Solution (v0.3.1)
- **Problem**: Fyne Tree widget caused flickering and position drift during expand/collapse
- **Root Cause**: Fyne Tree's internal refresh mechanism rebuilds entire node structure
- **Solution**: Custom tree implementation using VBox container and dynamic rendering
- **Benefits**: Smooth animations, no flickering, better performance, elegant design

### Chinese Font Support Solution
- **Environment Variable**: `FYNE_FONT=C:\Windows\Fonts\simhei.ttf`
- **Key Discovery**: Fyne doesn't support TTC font collections, requires individual TTF files
- **Working Font**: SimHei (黑体) provides excellent Chinese character support
- **Implementation**: Set in main.go before app initialization

### GUI Design Constraints
- **Window Size**: Keep the default window size at 900x650 pixels - DO NOT modify this size
- The current modern layout design should work within this constraint

### Build Dependencies
- **TDM-GCC 10.3.0**: Required for CGO compilation (Fyne → OpenGL)
- **Proxy settings**: May be needed for Go module downloads
- **Windows-specific**: Build script assumes Windows with specific TDM-GCC path

### Configuration File References
- `customer.conf`: Real-world reference configuration
- `dhf-real-schema.yaml`: Primary schema based on customer.conf analysis
- `dhf-schema-en.yaml`: English version for development
- `dhf-schema.yaml`: Chinese version (has display issues)