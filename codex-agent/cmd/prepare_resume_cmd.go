package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/codex-agent/internal/agent"
)

func newPrepareResumeCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "prepare-resume",
		Short: "unarchive and validate a session before resuming it",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params agent.PrepareResumeParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %w", err)
			}

			session, err := app.Agent.PrepareResume(cmd.Context(), params)
			if err != nil {
				return err
			}

			return json.NewEncoder(cmd.OutOrStdout()).Encode(session)
		},
	}
}
