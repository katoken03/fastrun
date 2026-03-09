package cmd

import (
	"fmt"
	"os"

	"github.com/katoken03/fastrun/internal/config"
	"github.com/katoken03/fastrun/internal/runner"
	uiPkg "github.com/katoken03/fastrun/internal/ui"
	"github.com/katoken03/fastrun/plugins/make"
	"github.com/katoken03/fastrun/plugins/npm"
	"github.com/spf13/cobra"
)

var (
	textOnly bool
	version  string
	commit   string
	date     string
)

var rootCmd = &cobra.Command{
	Use:   "fastrun",
	Short: "fastrun is a command launcher",
	Long: `fastrun is a command launcher that helps you run npm scripts and make targets
without remembering the exact command names.`,
	RunE: runCommand,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

func init() {
	// テキストモードのフラグを追加
	rootCmd.Flags().BoolVarP(&textOnly, "text-only", "t", false, "Print selected command text only without execution")
}

func runCommand(cmd *cobra.Command, args []string) error {
	// 設定を読み込む
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Initialize runners
	runners := []runner.CommandRunner{
		&npm.Runner{},
		&make.Runner{},
	}

	// Collect all available commands
	var commands []runner.Command
	for _, r := range runners {
		cmds, err := r.ParseCommands(cwd)
		if err == nil {
			commands = append(commands, cmds...)
		}
	}

	if len(commands) == 0 {
		return fmt.Errorf("no commands found in current directory")
	}

	// Show UI and get selected command
	ui := uiPkg.NewUI(commands, cfg)
	selectedCmd, err := ui.Show()
	if err != nil {
		// Check if the error is a user cancellation (ESC key or Ctrl+C)
		if uiPkg.IsCancelled(err) {
			// User cancelled - exit silently without error message
			return nil
		}
		return fmt.Errorf("UI error: %w", err)
	}

	// テキストのみモードの場合は、コマンドのテキストを出力して終了
	if textOnly {
		// 実行コマンドだけをテキストとして出力
		for _, r := range runners {
			cmds, err := r.ParseCommands(cwd)
			if err != nil {
				continue
			}
			for _, cmd := range cmds {
				if cmd.Name == selectedCmd.Name {
					fmt.Println(cmd.ExecuteCommand)
					return nil
				}
			}
		}
		return fmt.Errorf("command not found: %s", selectedCmd.Name)
	}

	// Find the appropriate runner
	for _, r := range runners {
		cmds, err := r.ParseCommands(cwd)
		if err != nil {
			continue
		}
		for _, cmd := range cmds {
			if cmd.Name == selectedCmd.Name {
				return r.RunCommand(*selectedCmd)
			}
		}
	}

	return fmt.Errorf("runner not found for command: %s", selectedCmd.Name)
}
