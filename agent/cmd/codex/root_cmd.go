package codex

import (
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/codex/client/http"
	"multienv/agent/codex/service"
	"multienv/agent/config"
	hosthttp "multienv/agent/host/client/http"
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
			instance, err := service.New(client, hosthttp.New(cfg.Client))
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
	cmd.AddCommand(newRenameSessionCmd(&app))
	cmd.AddCommand(newArchiveSessionCmd(&app))
	cmd.AddCommand(newUnarchiveSessionCmd(&app))
	cmd.AddCommand(newPrepareResumeCmd(&app))
	return cmd
}
