# fastrun

A command launcher tool that helps you run npm scripts and make targets without remembering exact command names.

## Features

- Fuzzy finding for commands using fzf
- Support for npm scripts and Makefile targets
- Configurable through global and project-specific settings
- Extensible plugin system

## Installation

```bash
# macOS
brew install fastrun

# Windows
winget install fastrun

# Linux
apt-get install fastrun
```

## Usage

Simply run `fastrun` in a directory containing a `package.json` or `Makefile`:

```bash
fastrun
```

## Configuration

Global configuration: `~/.config/fastrun/config.json`
Project configuration: `.fastrun/config.json`

## License

MIT

