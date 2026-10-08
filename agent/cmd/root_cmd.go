package main

import (
	"github.com/spf13/cobra"

	"multienv/agent/cmd/codex"
	"multienv/agent/cmd/tmux"
)

func newRootCmd() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:               "multienv-agent",
		Short:             "manage tmux and Codex sessions across connected workspace",
		SilenceErrors:     true,
		SilenceUsage:      true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	cmd.PersistentFlags().StringVarP(
		&configPath,
		"config",
		"c",
		"~/.config/multienv/agent.yaml",
		"path to configuration file",
	)
	cmd.AddCommand(tmux.NewCommand(&configPath), codex.NewCommand(&configPath))
	return cmd
}
