package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/tmux-agent/internal/agent"
)

func newNewSessionCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-session",
		Short: "create a new session in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params agent.NewSessionParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %v", err)
			}

			session, err := app.Agent.NewSession(cmd.Context(), params)
			if err != nil {
				return err
			}

			return json.NewEncoder(cmd.OutOrStdout()).Encode(session)
		},
	}

	return cmd
}
