package make

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/katoken03/fastrun/internal/runner"
)

type Runner struct{}

// 無視するターゲットのリスト
var ignoredTargets = map[string]bool{
	".PHONY": true,
	"@echo":  true,
}

// extractComment はコメント行からコメントテキストを返す。コメント行でなければ false を返す。
func extractComment(line string) (string, bool) {
	if !strings.HasPrefix(line, "#") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, "#")), true
}

// isVariableAssignment は Makefile の変数代入行かどうかを判定する。
// VAR = value / VAR := value / VAR ?= value 等をすべて対象とする。
func isVariableAssignment(line string) bool {
	idx := strings.Index(line, ":")
	if idx == -1 {
		return false
	}
	// = が : より前にある場合は変数代入（VAR = value, VAR ?= value 等）
	if eqIdx := strings.Index(line, "="); eqIdx != -1 && eqIdx < idx {
		return true
	}
	// := は変数代入
	return idx+1 < len(line) && line[idx+1] == '='
}

// extractTargetName はターゲット行からターゲット名を返す。ターゲット行でなければ false を返す。
func extractTargetName(line string) (string, bool) {
	if isVariableAssignment(line) {
		return "", false
	}
	idx := strings.Index(line, ":")
	if idx == -1 {
		return "", false
	}
	return strings.TrimSpace(line[:idx]), true
}

// shouldIgnoreTarget は表示対象外のターゲットかどうかを判定する。
func shouldIgnoreTarget(target string) bool {
	if ignoredTargets[target] {
		return true
	}
	return strings.HasPrefix(target, ".") ||
		strings.HasPrefix(target, "@")
}

func (r *Runner) ParseCommands(path string) ([]runner.Command, error) {
	makefilePath := filepath.Join(path, "Makefile")
	file, err := os.Open(makefilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Makefile: %w", err)
	}
	defer file.Close()

	var commands []runner.Command
	seenTargets := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	var currentComment string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if comment, ok := extractComment(line); ok {
			currentComment = comment
			continue
		}

		if target, ok := extractTargetName(line); ok && !shouldIgnoreTarget(target) && !seenTargets[target] {
			seenTargets[target] = true
			commands = append(commands, runner.Command{
				Name:           target,
				Source:         "make",
				Description:    currentComment,
				ExecuteCommand: fmt.Sprintf("make %s", target),
			})
		}
		currentComment = ""
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
