# Configuration System

## Configuration Levels
The tool supports two levels of configuration with a priority hierarchy:

### 1. Global Configuration
- **Location**: `~/.config/f/config.json`
- **Scope**: Applied across all projects
- **Use case**: Default plugins, user preferences

### 2. Project Configuration  
- **Location**: `.f/config.json` (in project root)
- **Scope**: Project-specific settings
- **Priority**: Overrides global configuration

## Configuration Format
```json
{
  "plugins": ["npm", "make"],
  "defaultTool": "npm"
}
```

## Configuration Features
- JSON format with comments support
- Plugin selection and ordering
- Default tool specification
- Extensible for future options

## Configuration Management
- Handled in `/internal/config/` package
- Type definitions in `config/types.go`
- Loading logic in `config/config.go`

## Plugin System
- Located in `/plugins/` directory
- Each plugin implements `CommandRunner` interface
- Current plugins:
  - `npm`: NPM scripts support
  - `make`: Makefile targets support

## Setup Directory
During installation, the tool creates:
- `~/.config/fastrun/` directory for configuration
- This ensures proper configuration file location