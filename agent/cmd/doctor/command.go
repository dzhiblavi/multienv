package doctor

import (
	"fmt"

	"github.com/spf13/cobra"

	"multienv/agent/doctor/service"
)

func NewCommand(configPath *string) *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:     "doctor",
		Short:   "diagnose the local host or the entire workspace",
		Example: "  multienv-agent doctor global\n  multienv-agent doctor local --format json",
	}

	cmd.PersistentFlags().StringVar(&format, "format", "human", "output format: human or json")

	for _, scope := range []string{"local", "global"} {
		cmd.AddCommand(&cobra.Command{
			Use:   scope,
			Short: "run " + scope + " diagnostics",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				if format != "human" && format != "json" {
					return fmt.Errorf("unknown format %q: use human or json", format)
				}

				collect := service.Global
				if scope == "local" {
					collect = service.Local
				}

				result := collect(cmd.Context(), *configPath)
				if err := render(cmd.OutOrStdout(), format, result); err != nil {
					return err
				}

				if result.Failed() {
					return fmt.Errorf("doctor found failed checks")
				}

				return nil
			},
		})
	}

	return cmd
}
