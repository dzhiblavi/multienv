package codex

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	agent "multienv/agent/internal/codexagent"
)

func newUnarchiveSessionCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unarchive-session",
		Short: "unarchive session in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params agent.UnarchiveSessionParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %v", err)
			}

			return app.Agent.UnarchiveSession(cmd.Context(), params)
		},
	}

	return cmd
}
