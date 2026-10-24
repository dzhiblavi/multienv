package tmux

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/tmux/service"
)

func newMoveWindowCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "move-window",
		Short: "move a window to a new index in its session",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params service.MoveWindowParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %w", err)
			}

			return app.Agent.MoveWindow(cmd.Context(), params)
		},
	}
}
