package app

import (
	"context"
	"fmt"
	"regexp"

	tea "charm.land/bubbletea/v2"

	codexmodel "multienv/agent/codex/model"
	codexservice "multienv/agent/codex/service"
	tmuxmodel "multienv/agent/tmux/model"
	tmuxservice "multienv/agent/tmux/service"
	"multienv/common/log"
	"multienv/tui/internal/browser"
	"multienv/tui/internal/dialog"
	"multienv/tui/internal/display"
	"multienv/tui/internal/navigation"
)

type codexReady struct {
	session codexmodel.Session
	tmux    tmuxmodel.SessionsByHost
	err     error
}

func (m *model) openCodex() tea.Cmd {
	selected, ok := m.codex.Selected()
	if !ok {
		return nil
	}

	m.busy = true
	return func() tea.Msg {
		result := codexReady{}
		snapshot, err := m.services.Codex.ListAll(m.ctx)
		if err != nil {
			result.err = err
			return result
		}

		session, ok := snapshot[selected.Hostname][selected.Id]
		if !ok {
			result.err = fmt.Errorf("conversation no longer available")
			return result
		}

		session.Hostname = selected.Hostname
		session.Id = selected.Id
		result.session = session
		if session.TmuxPaneId != "" {
			return result
		}

		result.session, result.err = m.services.Codex.PrepareResume(m.ctx, codexservice.PrepareResumeParams{
			Hostname: session.Hostname, SessionId: session.Id,
		})
		if result.err != nil {
			return result
		}

		if m.options.ResumeInPlace && session.Hostname == m.hostname {
			return result
		}

		result.tmux, result.err = m.services.Tmux.ListAll(m.ctx)
		return result
	}
}

func (m *model) resumeDialog(msg codexReady) tea.Cmd {
	m.busy = false
	if msg.err != nil {
		return m.error(msg.err)
	}

	s := msg.session
	if s.TmuxPaneId != "" {
		return m.finish(navigation.Action{Hostname: s.Hostname, Target: s.TmuxPaneId})
	}

	if m.options.ResumeInPlace && s.Hostname == m.hostname {
		return m.finish(navigation.Action{ResumeID: s.Id, WorkDir: s.WorkDir})
	}

	items := []browser.Item{}
	for id, session := range msg.tmux[s.Hostname] {
		items = append(items, browser.Item{
			ID: id, Label: display.Clean(session.Name), Detail: display.Clean(session.CurrentPath),
			Search: session.Name, Value: id})
	}

	if len(items) == 0 {
		return m.error(fmt.Errorf("create a tmux session on %s first", s.Hostname))
	}

	return m.choose("Resume in tmux session · "+m.cfg.Alias(s.Hostname), items, func(result dialog.Result) tea.Cmd {
		return m.resumeWindow(s, result.Item.Value.(string))
	})
}

var paneID = regexp.MustCompile(`^%[0-9]+$`)

func (m *model) resumeWindow(session codexmodel.Session, tmuxID string) tea.Cmd {
	return m.performAction(func(ctx context.Context) (*navigation.Action, error) {
		s, err := m.services.Codex.PrepareResume(ctx, codexservice.PrepareResumeParams{
			Hostname: session.Hostname, SessionId: session.Id,
		})
		if err != nil {
			return nil, err
		}

		name := session.Id
		if len(name) > 8 {
			name = name[:8]
		}

		result, err := m.services.Tmux.NewWindow(ctx, tmuxservice.NewWindowParams{
			Hostname: session.Hostname, SessionId: tmuxID, Name: "codex-" + name, WorkDir: s.WorkDir,
			Command: []string{"multienv-tui", "--resume-when-visible", session.Id},
		})
		if err != nil {
			return nil, err
		}

		if !paneID.MatchString(result.TmuxPaneId) {
			return nil, fmt.Errorf("new window response has no pane ID")
		}

		if err := m.services.Refresh(ctx, "tmux", []string{session.Hostname}); err != nil {
			log.Warn(ctx, "window created but snapshot refresh failed", "error", err)
		}

		return &navigation.Action{Hostname: session.Hostname, Target: result.TmuxPaneId}, nil
	})
}
