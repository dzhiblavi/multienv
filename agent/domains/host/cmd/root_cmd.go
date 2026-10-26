package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/config"
	"multienv/agent/domains/host/client/http"
	"multienv/common/util"
)

func NewCommand(configPath *string) *cobra.Command {
	app := App{}
	cmd := &cobra.Command{
		Use:   "host",
		Short: "execute commands across connected hosts",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			var cfg config.Config
			if err := util.ReadYaml(cmd.Context(), *configPath, &cfg); err != nil {
				return fmt.Errorf("failed to parse config: %w", err)
			}

			app.Client = http.New(cfg.Client)
			return nil
		},
	}

	cmd.AddCommand(newInvokeCmd(&app))
	return cmd
}
