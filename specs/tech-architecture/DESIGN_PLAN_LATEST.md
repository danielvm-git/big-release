# Design Plan

## CLI Interface
- `big-release` — run release
- `big-release validate` — validate config
- `big-release version` — show version

## Configuration
- YAML-based configuration
- Auto-detection of project type
- Validation with clear error messages

## Plugin System
- Registry-based plugin loading
- Built-in changelog plugin
- Extensible for custom plugins

## Publisher System
- Interface-based publisher pattern
- Language-specific implementations
- Registry for publisher discovery
