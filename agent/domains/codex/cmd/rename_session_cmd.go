package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/domains/codex/service"
)

func newRenameSessionCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "rename-session",
		Short: "rename a session in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params service.RenameSessionParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %w", err)
			}

			return app.Agent.RenameSession(cmd.Context(), params)
		},
	}
}
