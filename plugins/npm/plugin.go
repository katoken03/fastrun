package npm

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/tidwall/gjson"
    "github.com/kato/fastrun/internal/runner"
)

type Runner struct{}

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

    var commands []runner.Command
    scripts.ForEach(func(key, value gjson.Result) bool {
        cmd := runner.Command{
            Name:           key.String(),
            Description:    value.String(),
            ExecuteCommand: fmt.Sprintf("npm run %s", key.String()),
        }
        commands = append(commands, cmd)
        return true
    })

    return commands, nil
}

func (r *Runner) RunCommand(cmd runner.Command) error {
    c := exec.Command("npm", "run", cmd.Name)
    c.Stdout = os.Stdout
    c.Stderr = os.Stderr
    c.Stdin = os.Stdin
    return c.Run()
}
