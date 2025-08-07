# Code Style and Conventions

## General Go Conventions
- Follow standard Go conventions and idioms
- Use meaningful variable and function names
- Add comments for complex logic (in Japanese as per rules)
- Keep functions focused and small
- Error handling should use Go's standard `error` type with meaningful messages

## Language Requirements
- **Comments**: Japanese (as specified in .windsurfrules)
- **Error messages**: English (as specified in .windsurfrules)
- **AI chat responses**: Japanese

## Error Handling Best Practices
```go
// Preferred pattern
if err != nil {
    return fmt.Errorf("failed to read config: %w", err)
}

// Avoid simple error forwarding
if err != nil {
    return err
}
```

## Command Execution Guidelines
- Always handle STDIN/STDOUT/STDERR properly
- Use `os/exec` package for running commands
- Preserve terminal settings for interactive commands

## Plugin Development Guidelines
- Each plugin must implement the `CommandRunner` interface
- Plugins should be self-contained
- Error messages should be user-friendly
- Handle edge cases gracefully

## Plugin Interface
```go
type CommandRunner interface {
    ParseCommands(path string) ([]Command, error)
    RunCommand(cmd Command) error
}
```

## File Structure Standards
```
/cmd             - CLI implementation
/internal        - Internal packages  
/plugins         - Build tool plugins
```

## Dependencies Policy
- Minimize external dependencies
- Use standard library when possible
- Required external tools (like fzf) should be documented

## Security Considerations
- Don't execute arbitrary commands
- Validate all user input
- Handle file permissions correctly
- Follow principle of least privilege