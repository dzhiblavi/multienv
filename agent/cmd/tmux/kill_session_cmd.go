package tmux

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	agent "multienv/agent/internal/tmuxagent"
)

func newKillSessionCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kill-session",
		Short: "kill a new session in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params agent.KillSessionParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %v", err)
			}

			return app.Agent.KillSession(cmd.Context(), params)
		},
	}

	return cmd
}
