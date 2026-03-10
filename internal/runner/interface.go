package runner

// Command represents a runnable command with its metadata
type Command struct {
    Name           string
    Source         string // "npm" | "make"
    Description    string
    ExecuteCommand string
}

// CommandRunner is the interface that all plugins must implement
type CommandRunner interface {
    // ParseCommands parses the configuration file and returns a list of available commands
    ParseCommands(path string) ([]Command, error)
    
    // RunCommand executes the specified command
    RunCommand(cmd Command) error
}
