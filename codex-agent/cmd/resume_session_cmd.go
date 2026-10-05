package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/codex-agent/internal/agent"
)

func newResumeSessionCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume-session",
		Short: "resume session in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var params agent.ResumeSessionParams
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&params); err != nil {
				return fmt.Errorf("failed to decode request: %v", err)
			}

			return app.Agent.ResumeSession(cmd.Context(), params)
		},
	}

	return cmd
}
