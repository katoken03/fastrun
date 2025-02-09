# fastrun

Instead of typing...
```bash
npm run dev
```
Isn't it much faster to just type:
```bash
f [Enter] [Enter]
```
This command launcher tool makes it possible!

Even with multiple commands in your project, you can use incremental search to filter them or select with arrow keys - no need to remember the exact command names.

## Features

- Quick command selection with fuzzy search using fzf
- Support for npm scripts and make targets
- Project-specific and global configurations
- Extensible plugin system

## Installation

### Prerequisites

- [Homebrew](https://brew.sh/) (for macOS users)
- OpenJDK 17 or later

### Quick Install

```bash
brew install katoken03/fastrun/fastrun
```

### Manual Installation

If you prefer to install from source:

1. Clone the repository
2. Build using Go 1.x
3. Install dependencies

## Usage

Basic command syntax:

```bash
f           # List and select available commands
f <pattern> # Filter commands by pattern
```

Examples:

```bash
f              # Shows all available commands
f test         # Shows commands containing "test"
f "run dev"    # Shows commands containing "run dev"
```

## Configuration

fastrun can be configured at two levels:

1. Global configuration: `~/.config/f/config.json`
2. Project configuration: `.f/config.json`

Example configuration:
```json
{
  "plugins": ["npm", "make"],
  "defaultTool": "npm"
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT License](LICENSE)

This tool is designed to improve developer productivity by eliminating the need to type long commands repeatedly, allowing you to focus more on development.