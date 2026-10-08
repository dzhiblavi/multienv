package codex

import (
	"github.com/spf13/cobra"
)

func newSyncLocalCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-local",
		Short: "sync local sessions with the server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return app.Agent.SyncLocal(cmd.Context())
		},
	}

	return cmd
}
