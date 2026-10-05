package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"multienv/codex-agent/internal/agent"
	"multienv/codex-agent/internal/client/http"
	"multienv/common/util"
)

func newRootCmd() *cobra.Command {
	var configPath string
	app := App{}

	cmd := &cobra.Command{
		Use:           "codex-agent",
		Short:         "execute codex commands across connected workspace",
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
		"~/.config/multienv/codex-agent.yaml",
		"path to configuration file",
	)

	cmd.AddCommand(newListLocalCmd(&app))
	cmd.AddCommand(newListAllCmd(&app))
	cmd.AddCommand(newSyncLocalCmd(&app))
	cmd.AddCommand(newSyncAllCmd(&app))
	cmd.AddCommand(newDeleteSessionCmd(&app))
	cmd.AddCommand(newArchiveSessionCmd(&app))
	cmd.AddCommand(newUnarchiveSessionCmd(&app))
	cmd.AddCommand(newPrepareResumeCmd(&app))
	return cmd
}
