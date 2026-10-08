package codex

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	agent "multienv/agent/internal/codexagent"
)

func newArchiveSessionCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive-session",
		Short: "archive session in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params agent.ArchiveSessionParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %v", err)
			}

			return app.Agent.ArchiveSession(cmd.Context(), params)
		},
	}

	return cmd
}
