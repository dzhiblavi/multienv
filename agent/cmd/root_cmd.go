package main

import (
	"github.com/spf13/cobra"

	codexcmd "multienv/agent/domains/codex/cmd"
	doctorcmd "multienv/agent/domains/doctor/cmd"
	filesystemcmd "multienv/agent/domains/filesystem/cmd"
	hostcmd "multienv/agent/domains/host/cmd"
	mastercmd "multienv/agent/domains/master/cmd"
	tmuxcmd "multienv/agent/domains/tmux/cmd"
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

	cmd.AddCommand(tmuxcmd.NewCommand(&configPath))
	cmd.AddCommand(codexcmd.NewCommand(&configPath))
	cmd.AddCommand(mastercmd.NewCommand(&configPath))
	cmd.AddCommand(doctorcmd.NewCommand(&configPath))
	cmd.AddCommand(filesystemcmd.NewCommand(&configPath))
	cmd.AddCommand(hostcmd.NewCommand(&configPath))
	return cmd
}
