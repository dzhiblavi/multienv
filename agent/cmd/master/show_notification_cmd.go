package master

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"multienv/common/master/model"
)

func newShowNotificationCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "show-notification",
		Short: "show a notification on the server host using JSON from stdin",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var notification model.Notification
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&notification); err != nil {
				return fmt.Errorf("failed to decode notification: %w", err)
			}
			return app.Agent.ShowNotification(cmd.Context(), notification)
		},
	}
}
