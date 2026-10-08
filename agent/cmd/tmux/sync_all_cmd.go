package tmux

import (
	"github.com/spf13/cobra"
)

func newSyncAllCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-all",
		Short: "sync all workspace sessions with the server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return app.Agent.SyncAll(cmd.Context())
		},
	}

	return cmd
}
