package tmux

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	agent "multienv/agent/internal/tmuxagent"
)

func newKillWindowCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "kill-window",
		Short: "kill a window in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params agent.KillWindowParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %w", err)
			}

			return app.Agent.KillWindow(cmd.Context(), params)
		},
	}
}
