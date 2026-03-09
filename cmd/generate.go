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


// テキストモードを使用した新しいシェル関数
const bashFunctionWithT = `
# ========== Add this to your ~/.bash_profile ==========
f() {
    local cmd=$(command f -t "$@")
    if [ $? -eq 0 ] && [ -n "$cmd" ]; then
        # Display command in cyan color (DisplayCommand equivalent)
        echo -e "\033[36m$cmd\033[0m"
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
        # Display command in cyan color (DisplayCommand equivalent)
        echo -e "\033[36m$cmd\033[0m"
        print -s "$cmd"
        eval "$cmd"
    fi
}
`
