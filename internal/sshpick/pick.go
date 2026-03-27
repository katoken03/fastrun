package sshpick

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/katoken03/fastrun/internal/config"
	"github.com/katoken03/fastrun/internal/sshconfig"
	uiPkg "github.com/katoken03/fastrun/internal/ui"
)

const sshPreviewComment = "\033[2;32m"
const sshPreviewReset = "\033[0m"

// colorizeSSHComments wraps lines whose first non-space rune is '#' in dim green (fzf preview).
func colorizeSSHComments(block string) string {
	if block == "" {
		return block
	}
	lines := strings.Split(block, "\n")
	for i, line := range lines {
		if strings.Contains(line, "\033") {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			lines[i] = sshPreviewComment + line + sshPreviewReset
		}
	}
	return strings.Join(lines, "\n")
}

// Pick shows fzf with one line per candidate: index\\talias\\tdescription.
func Pick(cfg *config.Config, candidates []sshconfig.Candidate) (alias string, err error) {
	if len(candidates) == 0 {
		return "", fmt.Errorf("no ssh hosts to choose from")
	}

	tmpDir, err := os.MkdirTemp("", "fastrun-ssh-*")
	if err != nil {
		return "", fmt.Errorf("mkdir temp: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	maxAliasLen := 0
	for _, c := range candidates {
		if n := len(c.Alias); n > maxAliasLen {
			maxAliasLen = n
		}
	}

	var input strings.Builder
	for i, c := range candidates {
		desc := c.HostName
		if desc == "" {
			desc = "-"
		}
		desc = uiPkg.DimColorize(desc, "white")
		prev := sshconfig.BlockText(c.Block)
		if prev == "" {
			prev = "(empty block)"
		} else {
			prev = colorizeSSHComments(prev)
		}
		path := filepath.Join(tmpDir, fmt.Sprintf("%d.txt", i))
		if err := os.WriteFile(path, []byte(prev), 0o600); err != nil {
			return "", fmt.Errorf("write preview: %w", err)
		}
		padded := c.Alias + strings.Repeat(" ", maxAliasLen-len(c.Alias))
		fmt.Fprintf(&input, "%d\t%s\t%s\n", i, padded, desc)
	}

	safeDir := strings.ReplaceAll(tmpDir, "'", "'\\''")
	preview := "cat '" + safeDir + "/{1}.txt'"

	fzfArgs := []string{
		"--ansi",
		"--no-multi",
		"--delimiter=\t",
		"--with-nth=2,3",
		// After --with-nth, displayed columns are 1=alias, 2=HostName/IP; search both.
		"--nth=1,2",
		"--preview", preview,
	}
	if cfg != nil && cfg.FzfPosition == "top" {
		fzfArgs = append(fzfArgs, "--reverse")
	}

	cmd := exec.Command("fzf", fzfArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(input.String())

	out, err := cmd.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			if st, ok := exit.Sys().(syscall.WaitStatus); ok && st.ExitStatus() == 130 {
				return "", uiPkg.NewCancelledError("selection cancelled by user")
			}
		}
		return "", fmt.Errorf("fzf execution failed: %w", err)
	}

	line := strings.TrimSpace(string(out))
	fields := strings.Split(line, "\t")
	if len(fields) < 2 {
		return "", fmt.Errorf("unexpected fzf output: %q", line)
	}
	return strings.TrimSpace(fields[1]), nil
}

// FormatOutputLine returns the shell-invokable line for text-only mode (bash eval-safe).
func FormatOutputLine(alias string) string {
	return "ssh " + strconv.Quote(alias)
}
