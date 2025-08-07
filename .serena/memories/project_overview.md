# Fastrun Project Overview

## Project Purpose
Fastrun is a command launcher tool written in Go that makes it faster to execute npm scripts and make targets. Instead of typing full commands like `npm run dev`, users can simply type `f` followed by Enter to get an interactive fuzzy-search interface for available commands.

## Key Features
- Quick command selection with fuzzy search using fzf
- Support for npm scripts and make targets
- Project-specific and global configurations
- Extensible plugin system
- Interactive command selection with arrow keys
- No need to remember exact command names

## Tech Stack
- **Language**: Go 1.21+
- **CLI Framework**: Cobra (spf13/cobra)
- **JSON Parsing**: tidwall/gjson
- **Fuzzy Finder**: fzf (external dependency)
- **Build System**: Make + Go modules

## Architecture
- `main.go` - Entry point, calls cmd.Execute()
- `/cmd` - CLI implementation using Cobra
  - `root.go` - Main command logic and execution
  - `generate.go` - Shell function generation for integration
- `/internal` - Internal packages
  - `/config` - Configuration management
  - `/ui` - User interface (fzf integration)
  - `/runner` - Command execution logic
  - `/history` - Command history functionality
- `/plugins` - Build tool plugins
  - `/npm` - NPM scripts plugin
  - `/make` - Makefile targets plugin

## Target Platform
- Primary: macOS (Darwin)
- Built with cross-platform support in mind
- Uses Homebrew for installation on macOS