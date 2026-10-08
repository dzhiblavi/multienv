package tmux

import (
	"context"
	"strings"
	"time"

	"multienv/common/subprocess"
	"multienv/common/tmux/model"
)

func NewSession(ctx context.Context, name string) (model.Session, error) {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"-u", "new-session", "-d",
			"-s", name,
			"-P", "-F", sessionFormat,
		},
	}

	result, err := subprocess.RunSuccess(ctx, cmd)
	if err != nil {
		return model.Session{}, err
	}

	session, err := parseSession(strings.TrimSuffix(result.Stdout, "\n"), time.Now())
	if err != nil {
		return model.Session{}, err
	}

	windows, err := listWindows(ctx, session.Id)
	if err != nil {
		return model.Session{}, err
	}

	if list, ok := windows[session.Id]; ok {
		session.Windows = list
	}

	return session, nil
}
