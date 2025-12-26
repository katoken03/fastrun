package config

type Config struct {
    FzfPosition      string `json:"fzf_position"` // "top" or "bottom"
    CommandColor     string `json:"command_color"` // ANSI color code for command names
}

// DefaultConfig provides default settings
func DefaultConfig() *Config {
    return &Config{
        FzfPosition:  "top",
        CommandColor: "cyan", // デフォルトは水色
    }
}
