package codex

import (
	"fmt"

	"github.com/spf13/cobra"

	agenthttp "multienv/agent/internal/client/agent/http"
	"multienv/agent/internal/codex/client/http"
	"multienv/agent/internal/codex/service"
	"multienv/agent/internal/config"
	"multienv/common/util"
)

func NewCommand(configPath *string) *cobra.Command {
	app := App{}

	cmd := &cobra.Command{
		Use:   "codex",
		Short: "execute codex commands across connected workspace",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var cfg config.Config
			err := util.ReadYaml(ctx, *configPath, &cfg)
			if err != nil {
				return fmt.Errorf("failed to parse config: %v", err)
			}

			client := http.New(cfg.Client)
			instance, err := service.New(client, agenthttp.New(cfg.Client))
			if err != nil {
				return err
			}

			app.Agent = instance
			return nil
		},
	}

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
