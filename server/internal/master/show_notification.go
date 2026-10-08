package master

import (
	"context"
	"fmt"
	"strings"

	"multienv/common/master/operations"
	"multienv/common/subprocess"
	"multienv/server/internal/dispatcher"
)

func (s Dispatcher) ShowNotification(
	ctx context.Context,
	req operations.ShowNotificationRequest,
) (*operations.ShowNotificationResponse, error) {
	notification := req.Notification
	if notification.Text == "" {
		return nil, fmt.Errorf("notification text is required: %w", dispatcher.ErrInvalidParams)
	}

	var script strings.Builder
	script.WriteString("on run argv\n    display notification (item 1 of argv)")
	args := []string{"-", notification.Text}

	for _, parameter := range []struct {
		name  string
		value string
	}{
		{"with title", notification.Title},
		{"subtitle", notification.Subtitle},
		{"sound name", notification.SoundName},
	} {
		if parameter.value == "" {
			continue
		}
		args = append(args, parameter.value)
		script.WriteString(fmt.Sprintf(" %s (item %d of argv)", parameter.name, len(args)-1))
	}

	script.WriteString("\nend run\n")

	_, err := subprocess.RunSuccess(ctx, subprocess.Command{
		Name:  "osascript",
		Stdin: script.String(),
		Args:  args,
	})
	if err != nil {
		return nil, fmt.Errorf("show notification: %w", err)
	}

	return &operations.ShowNotificationResponse{}, nil
}
