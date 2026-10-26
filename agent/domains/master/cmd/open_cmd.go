package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/common/master/operations"
)

func newOpenCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "open a URL, file or directory on the server host using JSON from stdin",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var req operations.OpenRequest
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&req); err != nil {
				return fmt.Errorf("failed to decode open request: %w", err)
			}
			return app.Agent.Open(cmd.Context(), req.Target)
		},
	}
}
