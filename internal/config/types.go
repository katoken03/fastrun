package config

type Config struct {
	FzfPosition  string `json:"fzf_position"`  // "top" or "bottom"
	CommandColor string `json:"command_color"` // ANSI color code for command names
	// UseNr is nil when unset (legacy: try nr if on PATH). If false, never use nr.
	UseNr *bool `json:"use_nr,omitempty"`
}

// DefaultConfig provides default settings
func DefaultConfig() *Config {
	return &Config{
		FzfPosition:  "top",
		CommandColor: "cyan", // デフォルトは水色
	}
}
