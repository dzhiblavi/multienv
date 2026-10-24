package tmux

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/tmux/service"
)

func newRenameWindowCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "rename-window",
		Short: "rename a window in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params service.RenameWindowParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %w", err)
			}

			return app.Agent.RenameWindow(cmd.Context(), params)
		},
	}
}
