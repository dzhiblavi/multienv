package tmux

import (
	"context"
	"strings"

	"multienv/common/subprocess"
	"multienv/common/tmux/model"
	"multienv/common/util"
)

func ListSessions(ctx context.Context) (model.SessionsById, error) {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"list-sessions",
			"-F",
			`{"id":"#{session_id}","name":"#{session_name}","window_count":#{session_windows}}`,
		},
	}

	result, err := subprocess.RunSuccess(ctx, cmd)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(result.Stdout, "\n")
	if len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}

	list, err := util.BulkUnmarshal[model.Session](lines)
	if err != nil {
		return nil, err
	}

	sessions := make(model.SessionsById)
	for _, session := range list {
		sessions[session.Id] = session
	}

	return sessions, nil
}
