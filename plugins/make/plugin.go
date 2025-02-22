package make

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kato/fastrun/internal/runner"
)

type Runner struct{}

// 無視するターゲットのリスト
var ignoredTargets = map[string]bool{
	".PHONY": true,
	"@echo":  true,
}

func shouldIgnoreTarget(target string) bool {
	if ignoredTargets[target] {
		return true
	}
	return strings.HasPrefix(target, ".") ||
		strings.HasPrefix(target, "@") ||
		strings.HasPrefix(target, "VERSION") ||
		strings.HasPrefix(target, "BUILD_TIME")
}

func (r *Runner) ParseCommands(path string) ([]runner.Command, error) {
	makefilePath := filepath.Join(path, "Makefile")
	file, err := os.Open(makefilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Makefile: %w", err)
	}
	defer file.Close()

	var commands []runner.Command
	seenTargets := make(map[string]bool) // 重複チェック用

	scanner := bufio.NewScanner(file)
	var currentComment string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// コメント行の処理
		if strings.HasPrefix(line, "#") {
			currentComment = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			continue
		}

		// ターゲット行の処理
		if idx := strings.Index(line, ":"); idx != -1 {
			target := strings.TrimSpace(line[:idx])
			// 重複と無視すべきターゲットをスキップ
			if !shouldIgnoreTarget(target) && !seenTargets[target] {
				seenTargets[target] = true
				cmd := runner.Command{
					Name:           target,
					Description:    currentComment,
					ExecuteCommand: fmt.Sprintf("make %s", target),
				}
				commands = append(commands, cmd)
			}
			currentComment = ""
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Makefile: %w", err)
	}

	return commands, nil
}

func (r *Runner) RunCommand(cmd runner.Command) error {
	fullCmd := fmt.Sprintf("make %s", cmd.Name)
	runner.DisplayCommand(fullCmd, "cyan")

	c := exec.Command("make", cmd.Name)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}
