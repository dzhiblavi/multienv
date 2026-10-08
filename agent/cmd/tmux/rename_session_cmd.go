package tmux

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	agent "multienv/agent/internal/tmuxagent"
)

func newRenameSessionCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename-session",
		Short: "rename a session in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params agent.RenameSessionParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %v", err)
			}

			return app.Agent.RenameSession(cmd.Context(), params)
		},
	}

	return cmd
}
