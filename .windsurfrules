# f - Command Launcher Instructions

## Project Overview
This is a command launcher tool written in Go that helps developers run npm scripts and make targets without remembering exact command names. It uses fzf for fuzzy finding and supports multiple build tools through a plugin system.

## Code Style and Conventions

### General
- Follow standard Go conventions and idioms
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and small
- Error handling should use Go's standard `error` type with meaningful messages

### File Structure
```
/cmd             - CLI implementation
/internal        - Internal packages
/plugins         - Build tool plugins
```

### Plugin Development Guidelines
- Each plugin must implement the `CommandRunner` interface
- Plugins should be self-contained
- Error messages should be user-friendly
- Handle edge cases gracefully

## Configuration
The tool supports two levels of configuration:
1. Global: `~/.config/f/config.json`
2. Project: `.f/config.json`

Configuration files support comments and use JSON format.

## Best Practices

### Error Handling
```go
// Preferred
if err != nil {
    return fmt.Errorf("failed to read config: %w", err)
}

// Avoid
if err != nil {
    return err
}
```

### Command Execution
- Always handle STDIN/STDOUT/STDERR properly
- Use `os/exec` package for running commands
- Preserve terminal settings for interactive commands

### Plugin Interface
```go
type CommandRunner interface {
    ParseCommands(path string) ([]Command, error)
    RunCommand(cmd Command) error
}
```

## Testing
- Write unit tests for core functionality
- Use table-driven tests where appropriate
- Mock external commands in tests
- Test error conditions explicitly

## Documentation
- Keep comments up to date
- Document public APIs
- Include examples in documentation
- Update README.md when adding features

## Dependencies
- Minimize external dependencies
- Use standard library when possible
- Required external tools (fzf) should be documented

## Security Considerations
- Don't execute arbitrary commands
- Validate all user input
- Handle file permissions correctly
- Follow principle of least privilege

# Copilot(Generative AI) Instructions
When generating code, do not omit anything and generate everything.
