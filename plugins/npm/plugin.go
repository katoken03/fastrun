package npm

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/tidwall/gjson"
    "github.com/kato/fastrun/internal/runner"
)

type Runner struct{
	packageManager string
}

// detectPackageManager detects the appropriate package manager based on lock files
func (r *Runner) detectPackageManager(path string) string {
	// Check for pnpm-lock.yaml
	if _, err := os.Stat(filepath.Join(path, "pnpm-lock.yaml")); err == nil {
		return "pnpm"
	}
	
	// Check for bun.lockb
	if _, err := os.Stat(filepath.Join(path, "bun.lockb")); err == nil {
		return "bun"
	}
	
	// Default to npm
	return "npm"
}

func (r *Runner) ParseCommands(path string) ([]runner.Command, error) {
    packageJSONPath := filepath.Join(path, "package.json")
    content, err := os.ReadFile(packageJSONPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read package.json: %w", err)
    }

    scripts := gjson.Get(string(content), "scripts")
    if !scripts.Exists() {
        return nil, fmt.Errorf("no scripts found in package.json")
    }

    // Detect the appropriate package manager
    r.packageManager = r.detectPackageManager(path)

    var commands []runner.Command
    scripts.ForEach(func(key, value gjson.Result) bool {
        cmd := runner.Command{
            Name:           key.String(),
            Description:    value.String(),
            ExecuteCommand: fmt.Sprintf("%s run %s", r.packageManager, key.String()),
        }
        commands = append(commands, cmd)
        return true
    })

    return commands, nil
}

func (r *Runner) RunCommand(cmd runner.Command) error {
    fullCmd := fmt.Sprintf("%s run %s", r.packageManager, cmd.Name)
    runner.DisplayCommand(fullCmd, "cyan")

    c := exec.Command(r.packageManager, "run", cmd.Name)
    c.Stdout = os.Stdout
    c.Stderr = os.Stderr
    c.Stdin = os.Stdin
    return c.Run()
}
