package invoke

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/config"
	"multienv/agent/host/client/http"
	"multienv/common/host/operations"
	"multienv/common/util"
)

func NewCommand(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "invoke",
		Short: "execute a command on selected hosts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var cfg config.Config
			if err := util.ReadYaml(cmd.Context(), *configPath, &cfg); err != nil {
				return fmt.Errorf("failed to parse config: %w", err)
			}

			var req operations.InvokeRequest
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&req); err != nil {
				return fmt.Errorf("failed to decode request: %w", err)
			}

			response, err := http.New(cfg.Client).Invoke(cmd.Context(), req)
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
