package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"multienv/common/util"
	"multienv/tmux-agent/internal/agent"
	"multienv/tmux-agent/internal/client/http"
)

func newRootCmd() *cobra.Command {
	var configPath string
	app := App{}

	cmd := &cobra.Command{
		Use:           "tmux-agent",
		Short:         "execute tmux commands across connected workspace",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var config Config
			err := util.ReadYaml(ctx, configPath, &config)
			if err != nil {
				return fmt.Errorf("failed to parse config: %v", err)
			}

			client := http.New(config.Client)
			agent, err := agent.New(config.Agent, client)
			if err != nil {
				return err
			}

			app.Agent = agent
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(
		&configPath,
		"config",
		"c",
		"~/.config/multienv/tmux-agent.yaml",
		"path to configuration file",
	)

	cmd.AddCommand(newListAllCmd(&app))
	cmd.AddCommand(newSyncLocalCmd(&app))
	cmd.AddCommand(newSyncAllCmd(&app))
	cmd.AddCommand(newNewSessionCmd(&app))
	cmd.AddCommand(newNewWindowCmd(&app))
	cmd.AddCommand(newKillSessionCmd(&app))
	cmd.AddCommand(newRenameSessionCmd(&app))
	return cmd
}
