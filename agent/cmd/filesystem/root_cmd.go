package filesystem

import (
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/config"
	"multienv/agent/filesystem/service"
	"multienv/agent/host/client/http"
	"multienv/common/util"
)

func NewCommand(configPath *string) *cobra.Command {
	app := App{}
	cmd := &cobra.Command{
		Use:   "filesystem",
		Short: "query files and directories across connected hosts",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			var cfg config.Config
			if err := util.ReadYaml(cmd.Context(), *configPath, &cfg); err != nil {
				return fmt.Errorf("failed to parse config: %w", err)
			}

			instance, err := service.New(http.New(cfg.Client))
			if err != nil {
				return err
			}

			app.Agent = instance
			return nil
		},
	}

	cmd.AddCommand(newDirectoryCandidatesCmd(&app))
	return cmd
}
