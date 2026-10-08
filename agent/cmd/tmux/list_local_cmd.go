package tmux

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

func newListLocalCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-local",
		Short: "list local tmux sessions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sessions, err := app.Agent.ListLocal(cmd.Context())
			if err != nil {
				return err
			}
			return json.NewEncoder(cmd.OutOrStdout()).Encode(sessions)
		},
	}
	return cmd
}
