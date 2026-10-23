package main

import (
	"github.com/spf13/cobra"

	"multienv/agent/cmd/codex"
	"multienv/agent/cmd/doctor"
	"multienv/agent/cmd/master"
	"multienv/agent/cmd/tmux"
)

func newRootCmd() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:               "multienv-agent",
		Short:             "manage workspace sessions and execute commands on the server host",
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

	cmd.AddCommand(tmux.NewCommand(&configPath))
	cmd.AddCommand(codex.NewCommand(&configPath))
	cmd.AddCommand(master.NewCommand(&configPath))
	cmd.AddCommand(doctor.NewCommand(&configPath))
	return cmd
}
