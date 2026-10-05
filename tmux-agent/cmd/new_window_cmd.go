package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/tmux-agent/internal/agent"
)

func newNewWindowCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "new-window",
		Short: "create a window in a tmux session in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params agent.NewWindowParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %w", err)
			}

			result, err := app.Agent.NewWindow(cmd.Context(), params)
			if err != nil {
				return err
			}
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		},
	}
}
