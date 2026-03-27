package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/katoken03/fastrun/internal/config"
	"github.com/katoken03/fastrun/internal/runner"
	"github.com/katoken03/fastrun/internal/sshconfig"
	"github.com/katoken03/fastrun/internal/sshpick"
	uiPkg "github.com/katoken03/fastrun/internal/ui"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Select a Host from SSH config and connect",
	Long:  `Lists Host entries from your SSH config (see ssh_config(5)), lets you pick one with fzf, then runs ssh.`,
	RunE:  runSSH,
}

func init() {
	rootCmd.AddCommand(sshCmd)
}

func runSSH(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unexpected arguments (try: %s ssh)", cmd.Root().Name())
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	path, err := sshconfig.ResolvePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("ssh config %q: %w", path, err)
	}

	blocks, err := sshconfig.ParseFile(path)
	if err != nil {
		return err
	}

	candidates := sshconfig.CandidatesFromBlocks(blocks)
	if len(candidates) == 0 {
		return fmt.Errorf("no selectable Host entries in %s (wildcard-only stanzas are skipped)", path)
	}

	alias, err := sshpick.Pick(cfg, candidates)
	if err != nil {
		if uiPkg.IsCancelled(err) {
			return nil
		}
		return err
	}

	textOnly, err := cmd.Root().Flags().GetBool("text-only")
	if err != nil {
		return fmt.Errorf("text-only flag: %w", err)
	}

	if textOnly {
		fmt.Println(sshpick.FormatOutputLine(alias))
		return nil
	}

	full := fmt.Sprintf("ssh %s", strconv.Quote(alias))
	runner.DisplayCommand(full, cfg.CommandColor)

	c := exec.Command("ssh", alias)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
