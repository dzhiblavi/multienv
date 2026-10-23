package tmux

import (
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/internal/client/tmux/http"
	"multienv/agent/internal/config"
	agent "multienv/agent/internal/tmuxagent"
	"multienv/common/util"
)

func NewCommand(configPath *string) *cobra.Command {
	app := App{}

	cmd := &cobra.Command{
		Use:   "tmux",
		Short: "execute tmux commands across connected workspace",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var cfg config.Config
			err := util.ReadYaml(ctx, *configPath, &cfg)
			if err != nil {
				return fmt.Errorf("failed to parse config: %v", err)
			}

			client := http.New(cfg.Client)
			agent, err := agent.New(client)
			if err != nil {
				return err
			}

			app.Agent = agent
			return nil
		},
	}

	cmd.AddCommand(newListAllCmd(&app))
	cmd.AddCommand(newListLocalCmd(&app))
	cmd.AddCommand(newPanePreviewCmd(&app))
	cmd.AddCommand(newSyncLocalCmd(&app))
	cmd.AddCommand(newSyncAllCmd(&app))
	cmd.AddCommand(newNewSessionCmd(&app))
	cmd.AddCommand(newNewWindowCmd(&app))
	cmd.AddCommand(newMoveWindowCmd(&app))
	cmd.AddCommand(newKillSessionCmd(&app))
	cmd.AddCommand(newRenameSessionCmd(&app))
	return cmd
}
