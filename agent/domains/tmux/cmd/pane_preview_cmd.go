package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	tmuxservice "multienv/agent/domains/tmux/service"
)

func newPanePreviewCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "pane-preview",
		Short: "capture the visible contents of a local or remote tmux pane",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var req tmuxservice.PanePreviewParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&req); err != nil {
				return fmt.Errorf("failed to decode request: %w", err)
			}

			response, err := app.Agent.PanePreview(cmd.Context(), req)
			if err != nil {
				return err
			}

			return json.NewEncoder(cmd.OutOrStdout()).Encode(response)
		},
	}
}
