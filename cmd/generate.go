package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var shellType string

var generateCmd = &cobra.Command{
	Use:   "generate-shell-function",
	Short: "Generate shell function to add command to history",
	Long: `Generate shell function that wraps fastrun command and adds
the executed command to shell history.

Example:
  # Generate for bash
  fastrun generate-shell-function --shell=bash

  # Generate for zsh
  fastrun generate-shell-function --shell=zsh`,
	RunE: generateShellFunction,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&shellType, "shell", "", "Shell type (bash or zsh)")
}

func generateShellFunction(cmd *cobra.Command, args []string) error {
	// シェルタイプが指定されていない場合は自動検出を試みる
	if shellType == "" {
		shellType = detectShell()
	}

	var shellFunction string
	var configFile string
	var appendCmd string

	switch shellType {
	case "bash":
		shellFunction = bashFunctionWithT
		configFile = "~/.bash_profile"
		appendCmd = "cat << 'EOF' >> ~/.bash_profile\n" + bashFunctionWithT + "\nEOF"
	case "zsh":
		shellFunction = zshFunctionWithT
		configFile = "~/.zshrc"
		appendCmd = "cat << 'EOF' >> ~/.zshrc\n" + zshFunctionWithT + "\nEOF"
	default:
		return fmt.Errorf("unsupported shell type: %s (supported: bash, zsh)", shellType)
	}

	fmt.Println(shellFunction)
	fmt.Println()
	fmt.Printf("# To add this function to your %s, run:\n", configFile)
	fmt.Println(appendCmd)
	fmt.Println()
	fmt.Println("# Then reload your shell configuration with:")
	if shellType == "bash" {
		fmt.Println("source ~/.bash_profile")
	} else {
		fmt.Println("source ~/.zshrc")
	}

	return nil
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		// デフォルトはbash
		return "bash"
	}

	// シェルパスから種類を判断
	if shell == "/bin/bash" || shell == "/usr/bin/bash" {
		return "bash"
	} else if shell == "/bin/zsh" || shell == "/usr/bin/zsh" {
		return "zsh"
	}

	// デフォルトはbash
	return "bash"
}

const bashFunction = `# ========== fastrun wrapper function for bash - LEGACY MODE
# ========== THIS IS THE OLD IMPLEMENTATION AND IS NOT RECOMMENDED
# ========== Add this to your ~/.bash_profile
f() {
    # Run the f command and capture the output
    local cmd_output=$(command f "$@")
    local exit_code=$?
    
    # 特殊なマーカー（FASTRUN_CMD:）を使用してコマンドを抽出
    local cmd=$(echo "$cmd_output" | grep "FASTRUN_CMD:" | sed 's/FASTRUN_CMD://')
    
    # Add command to history if not empty
    if [ -n "$cmd" ]; then
        history -s "$cmd"
    fi
    
    # 特殊なマーカー行を除いて出力を表示
    echo "$cmd_output" | grep -v "FASTRUN_CMD:"
    
    return $exit_code
}
`

const zshFunction = `# ========== fastrun wrapper function for zsh - LEGACY MODE
# ========== THIS IS THE OLD IMPLEMENTATION AND IS NOT RECOMMENDED
# ========== Add this to your ~/.zshrc
f() {
    # Run the f command and capture the output
    local cmd_output=$(command f "$@")
    local exit_code=$?
    
    # 特殊なマーカー（FASTRUN_CMD:）を使用してコマンドを抽出
    local cmd=$(echo "$cmd_output" | grep "FASTRUN_CMD:" | sed 's/FASTRUN_CMD://')
    
    # Add command to history if not empty
    if [ -n "$cmd" ]; then
        print -s "$cmd"
    fi
    
    # 特殊なマーカー行を除いて出力を表示
    echo "$cmd_output" | grep -v "FASTRUN_CMD:"
    
    return $exit_code
}
`

// テキストモードを使用した新しいシェル関数
const bashFunctionWithT = `
# ========== Add this to your ~/.bash_profile ==========
f() {
    local cmd=$(command f -t "$@")
    if [ $? -eq 0 ] && [ -n "$cmd" ]; then
        history -s "$cmd"
        eval "$cmd"
    fi
}
`

const zshFunctionWithT = `
# ========== Add this to your ~/.zshrc ==========
f() {
    local cmd=$(command f -t "$@")
    if [ $? -eq 0 ] && [ -n "$cmd" ]; then
        print -s "$cmd"
        eval "$cmd"
    fi
}
`
