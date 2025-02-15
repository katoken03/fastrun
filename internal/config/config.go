package config

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/tidwall/gjson"
)

// LoadConfig loads configuration from files
func LoadConfig() (*Config, error) {
    config := DefaultConfig()

    // まずグローバル設定を読み込む
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("failed to get home directory: %w", err)
    }

    globalConfigPath := filepath.Join(homeDir, ".config", "fastrun", "config.json")
    if err := loadConfigFile(globalConfigPath, config); err != nil {
        // グローバル設定がない場合は無視
        if !os.IsNotExist(err) {
            return nil, fmt.Errorf("failed to load global config: %w", err)
        }
    }

    // 次にプロジェクト固有の設定を読み込む
    cwd, err := os.Getwd()
    if err != nil {
        return nil, fmt.Errorf("failed to get current directory: %w", err)
    }

    projectConfigPath := filepath.Join(cwd, ".fastrun", "config.json")
    if err := loadConfigFile(projectConfigPath, config); err != nil {
        // プロジェクト設定がない場合は無視
        if !os.IsNotExist(err) {
            return nil, fmt.Errorf("failed to load project config: %w", err)
        }
    }

    return config, nil
}

func loadConfigFile(path string, config *Config) error {
    content, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    // fzf_positionの値を読み込む
    position := gjson.Get(string(content), "fzf_position")
    if position.Exists() {
        pos := position.String()
        if pos == "top" || pos == "bottom" {
            config.FzfPosition = pos
        }
    }

    // command_colorの値を読み込む
    commandColor := gjson.Get(string(content), "command_color")
    if commandColor.Exists() {
        config.CommandColor = commandColor.String()
    }

    return nil
}
