package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/common/host/operations"
)

func newInvokeCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "invoke",
		Short: "execute a command on selected hosts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var req operations.InvokeRequest
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&req); err != nil {
				return fmt.Errorf("failed to decode request: %w", err)
			}

			response, err := app.Client.Invoke(cmd.Context(), req)
			if err != nil {
				return err
			}

			if err := json.NewEncoder(cmd.OutOrStdout()).Encode(response); err != nil {
				return err
			}

			return response.Err()
		},
	}
}
