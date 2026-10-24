package codex

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/internal/codex/service"
)

func newDeleteSessionCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-session",
		Short: "delete session in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params service.DeleteSessionParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %v", err)
			}

			return app.Agent.DeleteSession(cmd.Context(), params)
		},
	}

	return cmd
}
