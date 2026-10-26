package app

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	tea "charm.land/bubbletea/v2"

	tmuxservice "multienv/agent/tmux/service"
	"multienv/tui/internal/browser"
	"multienv/tui/internal/dialog"
	"multienv/tui/internal/display"
	"multienv/tui/internal/navigation"
	"multienv/tui/internal/tmuxview"
)

func (m *model) activate() tea.Cmd {
	if m.busy {
		return nil
	}

	if m.mode == Codex {
		return m.openCodex()
	}

	session, ok := m.tmux.Selected()
	if !ok {
		return nil
	}

	target := session.Data.Id
	if m.tmux.WindowFocused {
		window, ok := m.tmux.Window()
		if !ok {
			return nil
		}

		target += ":" + window.Id
	}

	return m.finish(navigation.Action{Hostname: session.Host, Target: target})
}

func (m *model) renameTmux() tea.Cmd {
	s, ok := m.tmux.Selected()
	if !ok {
		return nil
	}

	if m.tmux.WindowFocused {
		window, ok := m.tmux.Window()
		if !ok {
			return nil
		}

		return m.text("Window name", window.Name, func(result dialog.Result) tea.Cmd {
			return m.performOn([]string{s.Host}, func(ctx context.Context) error {
				return m.services.Tmux.RenameWindow(ctx,
					tmuxservice.RenameWindowParams{Hostname: s.Host, SessionId: s.Data.Id, WindowId: window.Id, NewName: result.Text})
			})
		})
	}

	return m.text("Session name", s.Data.Name, func(result dialog.Result) tea.Cmd {
		return m.performOn([]string{s.Host}, func(ctx context.Context) error {
			return m.services.Tmux.RenameSession(ctx,
				tmuxservice.RenameSessionParams{Hostname: s.Host, SessionId: s.Data.Id, NewName: result.Text})
		})
	})
}

func (m *model) moveWindow() tea.Cmd {
	s, ok := m.tmux.Selected()
	if !ok {
		return nil
	}

	window, ok := m.tmux.Window()
	if !ok {
		return nil
	}

	return m.text("New window index", strconv.Itoa(window.Index), func(result dialog.Result) tea.Cmd {
		index, err := strconv.Atoi(result.Text)
		if err != nil || index < 0 {
			return m.error(fmt.Errorf("window index must be a non-negative integer"))
		}

		return m.performOn([]string{s.Host}, func(ctx context.Context) error {
			return m.services.Tmux.MoveWindow(ctx,
				tmuxservice.MoveWindowParams{Hostname: s.Host, SessionId: s.Data.Id, WindowId: window.Id, NewIndex: &index})
		})
	})
}

func (m *model) killTmux() tea.Cmd {
	if m.tmux.WindowFocused {
		s, ok := m.tmux.Selected()
		if !ok {
			return nil
		}

		window, ok := m.tmux.Window()
		if !ok {
			return nil
		}

		if len(s.Data.Windows) <= 1 {
			return m.error(fmt.Errorf("cannot kill the last window of a session"))
		}

		return m.perform(func(ctx context.Context) error {
			return m.services.Tmux.KillWindow(ctx,
				tmuxservice.KillWindowParams{Hostname: s.Host, SessionId: s.Data.Id, WindowId: window.Id})
		})
	}

	items := m.tmux.Browser.Targets()
	return m.perform(func(ctx context.Context) error {
		var failures []error
		for _, item := range items {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			s := item.Value.(tmuxview.Session)
			failures = append(failures, m.services.Tmux.KillSession(ctx,
				tmuxservice.KillSessionParams{Hostname: s.Host, SessionId: s.Data.Id}))
		}

		return errors.Join(failures...)
	})
}

func (m *model) createTmux() tea.Cmd {
	if m.tmux.WindowFocused {
		return m.windowDirectories()
	}

	hosts := []string{m.hostname}
	for _, host := range m.cfg.HostAliases {
		hosts = append(hosts, host)
	}

	for _, item := range m.tmux.Browser.Items {
		hosts = append(hosts, item.Value.(tmuxview.Session).Host)
	}

	slices.Sort(hosts)
	hosts = slices.Compact(hosts)
	items := make([]browser.Item, 0, len(hosts))
	for _, host := range hosts {
		items = append(items, browser.Item{ID: host, Label: display.HostAlias(host, m.cfg.Alias(host)),
			Detail: display.Clean(host), Search: host + " " + m.cfg.Alias(host), Value: host})
	}

	return m.choose("Choose host", items, func(result dialog.Result) tea.Cmd {
		host := result.Item.Value.(string)
		return m.text("Session name", "", func(result dialog.Result) tea.Cmd {
			return m.performOn([]string{host}, func(ctx context.Context) error {
				_, err := m.services.Tmux.NewSession(ctx, tmuxservice.NewSessionParams{Hostname: host, Name: result.Text})
				return err
			})
		})
	})
}
