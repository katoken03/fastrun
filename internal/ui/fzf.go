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

func (u *UI) dimColorize(text, color string) string {
    if code, ok := colorCodes[color]; ok {
        return fmt.Sprintf("\033[2;%sm%s\033[0m", code, text)
    }
    return fmt.Sprintf("\033[2m%s\033[0m", text)
}

func (u *UI) Show() (*runner.Command, error) {
    // Find the longest command name and source for proper alignment
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

    // コマンド一覧を文字列に変換
    var input strings.Builder
    for _, cmd := range u.commands {
        // 色付きのコマンド名を作成
        coloredName := u.colorize(cmd.Name, u.config.CommandColor)
        namePadding := strings.Repeat(" ", maxNameLen-len(cmd.Name))

        // ソース種別（薄いグレー）
        coloredSource := u.dimColorize(cmd.Source, "white")
        sourcePadding := strings.Repeat(" ", maxSourceLen-len(cmd.Source))

        // 説明文
        description := cmd.Description

        fmt.Fprintf(&input, "%s%s  %s%s  %s\n", coloredName, namePadding, coloredSource, sourcePadding, description)
    }

    // fzfのオプションを設定
    fzfArgs := []string{
        "--ansi",            // ANSIカラーコードを解釈
        "--no-multi",        // 単一選択
        "--delimiter= ",     // スペースをデリミタとして使用
        "--nth=1",          // 最初のフィールドでマッチング
    }

    // 位置の設定を追加
    if u.config.FzfPosition == "top" {
        fzfArgs = append(fzfArgs, "--reverse")
    }

    cmd := exec.Command("fzf", fzfArgs...)
    cmd.Stderr = os.Stderr
    cmd.Stdin = strings.NewReader(input.String())

    output, err := cmd.Output()
    if err != nil {
        // Check if the error is due to exit status 130 (ESC key or Ctrl+C)
        if exitError, ok := err.(*exec.ExitError); ok {
            if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
                if status.ExitStatus() == 130 {
                    // User cancelled with ESC or Ctrl+C - this is not an error
                    return nil, NewCancelledError("selection cancelled by user")
                }
            }
        }
        return nil, fmt.Errorf("fzf execution failed: %w", err)
    }

    // 選択された行からコマンドを特定（色制御文字を除去）
    selected := strings.TrimSpace(string(output))
    selectedName := strings.Split(selected, " ")[0]
    selectedName = strings.TrimPrefix(selectedName, "\033[36m") // シアンの制御文字を除去
    selectedName = strings.TrimSuffix(selectedName, "\033[0m")  // リセット制御文字を除去

    // 対応するコマンドを探す
    for _, cmd := range u.commands {
        if cmd.Name == selectedName {
            return &cmd, nil
        }
    }

    return nil, fmt.Errorf("command not found: %s", selectedName)
}
