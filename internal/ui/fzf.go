package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/katoken03/fastrun/internal/config"
	"github.com/katoken03/fastrun/internal/runner"
)

type UI struct {
	commands []runner.Command
	config   *config.Config
}

func NewUI(commands []runner.Command, config *config.Config) *UI {
	return &UI{
		commands: commands,
		config:   config,
	}
}

// ANSIカラーコードのマップ
var colorCodes = map[string]string{
	"black":   "30",
	"red":     "31",
	"green":   "32",
	"yellow":  "33",
	"blue":    "34",
	"magenta": "35",
	"cyan":    "36",
	"white":   "37",
}

func (u *UI) colorize(text, color string) string {
	if code, ok := colorCodes[color]; ok {
		return fmt.Sprintf("\033[%sm%s\033[0m", code, text)
	}
	return text
}

// DimColorize applies dim + named foreground color (e.g. white for secondary columns).
func DimColorize(text, color string) string {
	if code, ok := colorCodes[color]; ok {
		return fmt.Sprintf("\033[2;%sm%s\033[0m", code, text)
	}
	return fmt.Sprintf("\033[2m%s\033[0m", text)
}

func (u *UI) dimColorize(text, color string) string {
	return DimColorize(text, color)
}

// buildFzfInput はコマンド一覧を fzf への入力文字列に整形する
func (u *UI) buildFzfInput() string {
	maxNameLen := 0
	maxSourceLen := 0
	for _, cmd := range u.commands {
		if len(cmd.Name) > maxNameLen {
			maxNameLen = len(cmd.Name)
		}
		if len(cmd.Source) > maxSourceLen {
			maxSourceLen = len(cmd.Source)
		}
	}

	var input strings.Builder
	for _, cmd := range u.commands {
		coloredName := u.colorize(cmd.Name, u.config.CommandColor)
		namePadding := strings.Repeat(" ", maxNameLen-len(cmd.Name))

		coloredSource := u.dimColorize(cmd.Source, "white")
		sourcePadding := strings.Repeat(" ", maxSourceLen-len(cmd.Source))

		fmt.Fprintf(&input, "%s%s  %s%s  %s\n", coloredName, namePadding, coloredSource, sourcePadding, cmd.Description)
	}
	return input.String()
}

// runFzf は fzf を起動してユーザーが選択した行を返す
func (u *UI) runFzf(input string) (string, error) {
	fzfArgs := []string{
		"--ansi",
		"--no-multi",
		"--delimiter= ",
		"--nth=1",
	}
	if u.config.FzfPosition == "top" {
		fzfArgs = append(fzfArgs, "--reverse")
	}

	cmd := exec.Command("fzf", fzfArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(input)

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 130 {
					return "", NewCancelledError("selection cancelled by user")
				}
			}
		}
		return "", fmt.Errorf("fzf execution failed: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// findCommandByName は選択された行からコマンド名を取り出し、対応する Command を返す
func (u *UI) findCommandByName(selected string) (*runner.Command, error) {
	selectedName := strings.Split(selected, " ")[0]
	selectedName = strings.TrimPrefix(selectedName, "\033[36m")
	selectedName = strings.TrimSuffix(selectedName, "\033[0m")

	for _, cmd := range u.commands {
		if cmd.Name == selectedName {
			return &cmd, nil
		}
	}
	return nil, fmt.Errorf("command not found: %s", selectedName)
}

func (u *UI) Show() (*runner.Command, error) {
	input := u.buildFzfInput()

	selected, err := u.runFzf(input)
	if err != nil {
		return nil, err
	}

	return u.findCommandByName(selected)
}
