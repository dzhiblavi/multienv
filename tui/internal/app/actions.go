package app

import (
	"context"
	"errors"
	"fmt"

	tea "charm.land/bubbletea/v2"

	codexmodel "multienv/agent/domains/codex/model"
	codexservice "multienv/agent/domains/codex/service"
	"multienv/common/log"
	"multienv/tui/internal/navigation"
	"multienv/tui/internal/tmuxview"
)

type actionDone struct {
	mode   Mode
	action *navigation.Action
	err    error
}

func (m *model) perform(fn func(context.Context) error) tea.Cmd {
	return m.performOn(m.affectedHosts(), fn)
}

func (m *model) performOn(hosts []string, fn func(context.Context) error) tea.Cmd {
	mode := string(m.mode)
	return m.performAction(func(ctx context.Context) (*navigation.Action, error) {
		err := fn(ctx)
		refreshErr := m.services.Refresh(ctx, mode, hosts)
		if err == nil && refreshErr != nil {
			return nil, fmt.Errorf("operation completed; refresh failed: %w", refreshErr)
		}

		return nil, errors.Join(err, refreshErr)
	})
}

func (m *model) affectedHosts() []string {
	if m.mode == Tmux && m.tmux.WindowFocused {
		if session, ok := m.tmux.Selected(); ok {
			return []string{session.Host}
		}
	}

	var hosts []string
	for _, item := range m.browser().Targets() {
		if m.mode == Codex {
			hosts = append(hosts, item.Value.(codexmodel.Session).Hostname)
		} else {
			hosts = append(hosts, item.Value.(tmuxview.Session).Host)
		}
	}

	return hosts
}

func (m *model) performAction(fn func(context.Context) (*navigation.Action, error)) tea.Cmd {
	if m.busy {
		return nil
	}

	m.busy = true
	m.status = ""
	mode := m.mode
	return func() tea.Msg {
		action, err := fn(m.ctx)
		return actionDone{mode: mode, action: action, err: err}
	}
}

func (m *model) actionFinished(msg actionDone) tea.Cmd {
	m.busy = false
	if msg.err != nil {
		log.Warn(m.ctx, "TUI action failed", "error", msg.err)
		m.error(msg.err)
		return m.load(true)
	}

	if msg.action != nil {
		return m.finish(*msg.action)
	}

	if msg.mode == Tmux {
		m.tmux.Browser.ClearMarks()
	} else {
		m.codex.Browser.ClearMarks()
	}

	m.status = "Done"
	return m.load(true)
}

func (m *model) codexAction(action string) tea.Cmd {
	items := m.codex.Browser.Targets()
	return m.perform(func(ctx context.Context) error {
		var failures []error
		for _, item := range items {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			s := item.Value.(codexmodel.Session)
			var err error
			switch action {
			case "archive":
				err = m.services.Codex.ArchiveSession(ctx,
					codexservice.ArchiveSessionParams{Hostname: s.Hostname, SessionId: s.Id})
			case "unarchive":
				err = m.services.Codex.UnarchiveSession(ctx,
					codexservice.UnarchiveSessionParams{Hostname: s.Hostname, SessionId: s.Id})
			case "delete":
				err = m.services.Codex.DeleteSession(ctx,
					codexservice.DeleteSessionParams{Hostname: s.Hostname, SessionId: s.Id})
			}

			failures = append(failures, err)
		}

		return errors.Join(failures...)
	})
}
