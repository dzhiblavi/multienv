package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	filesystemservice "multienv/agent/domains/filesystem/service"
)

func newDirectoryCandidatesCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "directory-candidates",
		Short: "list suggested directories on a local or remote host",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var req filesystemservice.DirectoryCandidatesParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&req); err != nil {
				return fmt.Errorf("failed to decode request: %w", err)
			}

			response, err := app.Agent.DirectoryCandidates(cmd.Context(), req)
			if err != nil {
				return err
			}

			return json.NewEncoder(cmd.OutOrStdout()).Encode(response)
		},
	}
}
