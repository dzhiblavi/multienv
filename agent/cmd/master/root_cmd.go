package master

import (
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/config"
	"multienv/agent/master/client/http"
	"multienv/agent/master/service"
	"multienv/common/util"
)

func NewCommand(configPath *string) *cobra.Command {
	app := App{}
	cmd := &cobra.Command{
		Use:   "master",
		Short: "execute commands on the server host",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			var cfg config.Config
			if err := util.ReadYaml(cmd.Context(), *configPath, &cfg); err != nil {
				return fmt.Errorf("failed to parse config: %w", err)
			}
			app.Agent = service.New(http.New(cfg.Client))
			return nil
		},
	}

	cmd.AddCommand(newShowNotificationCmd(&app))
	cmd.AddCommand(newOpenCmd(&app))
	return cmd
}
