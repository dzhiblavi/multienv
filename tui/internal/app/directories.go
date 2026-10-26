package app

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"

	filesystemservice "multienv/agent/filesystem/service"
	tmuxservice "multienv/agent/tmux/service"
	"multienv/tui/internal/browser"
	"multienv/tui/internal/dialog"
	"multienv/tui/internal/display"
	"multienv/tui/internal/tmuxview"
)

type directoryChoices struct {
	session tmuxview.Session
	current string
	paths   []string
	err     error
}

func (m *model) windowDirectories() tea.Cmd {
	s, ok := m.tmux.Selected()
	if !ok || m.busy {
		return nil
	}

	window, ok := m.tmux.Window()
	if !ok {
		return nil
	}

	if window.CurrentPath == "" {
		return m.error(fmt.Errorf("selected window has no working directory"))
	}

	m.busy = true
	return func() tea.Msg {
		choices := directoryChoices{session: s, current: window.CurrentPath}
		request := filesystemservice.DirectoryCandidatesParams{Hostname: s.Host}
		response, err := m.services.Filesystem.DirectoryCandidates(m.ctx, request)
		choices.err = err
		if err == nil {
			choices.paths = response.Directories
		}

		return choices
	}
}

func (m *model) directoryDialog(msg directoryChoices) tea.Cmd {
	m.busy = false
	paths := append([]string{msg.current}, msg.paths...)
	items := []browser.Item{}
	seen := map[string]bool{}
	for _, path := range paths {
		if path == "" || seen[path] {
			continue
		}

		seen[path] = true
		items = append(items, browser.Item{ID: path, Label: display.Clean(path), Search: path, Value: path})
	}

	cmd := m.choose("New window directory · "+m.cfg.Alias(msg.session.Host), items, func(result dialog.Result) tea.Cmd {
		return m.performOn([]string{msg.session.Host}, func(ctx context.Context) error {
			_, err := m.services.Tmux.NewWindow(ctx, tmuxservice.NewWindowParams{Hostname: msg.session.Host,
				SessionId: msg.session.Data.Id, WorkDir: result.Item.Value.(string)})
			return err
		})
	})
	if msg.err != nil && m.modal != nil {
		m.modal.Note = "Directory lookup failed; current directory is available"
	}

	return cmd
}
