package codex

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

func newListAllCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-all",
		Short: "list all codex sessions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			sessions, err := app.Agent.ListAll(ctx)
			if err != nil {
				return err
			}

			return json.NewEncoder(cmd.OutOrStdout()).Encode(sessions)
		},
	}

	return cmd
}
