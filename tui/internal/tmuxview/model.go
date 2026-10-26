package tmuxview

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"multienv/agent/domains/tmux/model"
	"multienv/tui/internal/browser"
	"multienv/tui/internal/config"
	"multienv/tui/internal/display"
	"multienv/tui/internal/keybindings"
)

type Session struct {
	Host string
	Data model.Session
}

func (s Session) Key() string { return s.Host + "\x00" + s.Data.Id }

type Model struct {
	Browser       browser.Model
	Preview       viewport.Model
	WindowFocused bool
	WindowID      string
	Pane          string
}

func New(cfg config.Tmux) Model {
	return Model{Browser: browser.New(true, cfg.Keybindings), Preview: viewport.New()}
}

func (m *Model) Selected() (Session, bool) {
	item, ok := m.Browser.Selected()
	if !ok {
		return Session{}, false
	}

	session, ok := item.Value.(Session)
	return session, ok
}

func (m *Model) Window() (model.Window, bool) {
	session, ok := m.Selected()
	if !ok {
		return model.Window{}, false
	}

	windows := sortedWindows(session.Data.Windows)
	for _, window := range windows {
		if m.WindowFocused && window.Id == m.WindowID {
			return window, true
		}
	}

	for _, window := range windows {
		if window.Active {
			return window, true
		}
	}

	if len(windows) > 0 {
		return windows[0], true
	}

	return model.Window{}, false
}

func (m *Model) MoveWindow(delta int) {
	session, ok := m.Selected()
	if !ok {
		return
	}

	windows := sortedWindows(session.Data.Windows)
	current, _ := m.Window()
	index := slices.IndexFunc(windows, func(w model.Window) bool { return w.Id == current.Id })
	if len(windows) > 0 {
		m.WindowID = windows[max(0, min(len(windows)-1, index+delta))].Id
	}

	m.Preview.GotoTop()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if m.WindowFocused {
		if key, ok := msg.(tea.KeyPressMsg); ok {
			switch m.Browser.Keybindings.Action(key.String()) {
			case "previous", "previous-window":
				m.MoveWindow(-1)
			case "next", "next-window":
				m.MoveWindow(1)
			}
		}

		return nil
	}

	return m.Browser.Update(msg)
}

func sortedWindows(windows []model.Window) []model.Window {
	windows = slices.Clone(windows)
	slices.SortFunc(windows, func(a, b model.Window) int { return cmp.Compare(a.Index, b.Index) })
	return windows
}

func (m *Model) Set(snapshot model.SessionsByHost, cfg config.Config) {
	sessions := []Session{}
	for host, data := range snapshot {
		for id, session := range data {
			session.Id = id
			sessions = append(sessions, Session{host, session})
		}
	}

	slices.SortFunc(sessions, func(a, b Session) int {
		return cmp.Or(cmp.Compare(b.Data.LastAttached, a.Data.LastAttached), cmp.Compare(a.Host, b.Host),
			cmp.Compare(a.Data.Name, b.Data.Name), cmp.Compare(a.Data.Id, b.Data.Id))
	})
	items := make([]browser.Item, 0, len(sessions))
	for _, s := range sessions {
		plain := display.Clean(s.Data.Name + " · " + cfg.Alias(s.Host))
		label := display.Clean(s.Data.Name) + " · " + display.HostAlias(s.Host, cfg.Alias(s.Host))
		detail := display.Clean(fmt.Sprintf("%d windows · %s · %s",
			s.Data.WindowCount, s.Data.CurrentCommand, s.Data.CurrentPath))
		indicator := display.Indicator(s.Data.Attached > 0, s.Data.SyncedAt)
		items = append(items, browser.Item{ID: s.Key(), Label: indicator + " " + label,
			Detail: detail, Search: plain, AllSearch: plain + " " + detail, Value: s})
	}

	m.Browser.Set(items)
	if window, ok := m.Window(); ok {
		m.WindowID = window.Id
	} else {
		m.WindowFocused = false
		m.WindowID = ""
	}
}

func (m *Model) PreviewText() string {
	session, ok := m.Selected()
	if !ok {
		return "No session selected"
	}

	window, _ := m.Window()
	keys := m.Browser.Keybindings
	hints := keybindings.Hints(keys.Hint("focus-windows", "windows"), keys.Hint("focus-sessions", "sessions"),
		keys.Hint("preview-up", "scroll up"), keys.Hint("preview-down", "scroll down"))
	lines := []string{display.Style(display.Clean(session.Data.Name), "6"), hints, ""}
	for _, w := range sortedWindows(session.Data.Windows) {
		prefix := "  "
		if w.Active {
			prefix = "● "
		}

		if m.WindowFocused && w.Id == window.Id {
			prefix = "> "
		}

		label := display.Clean(fmt.Sprintf("%d: %s · %s · %s", w.Index, w.Name, w.CurrentCommand, w.CurrentPath))
		label = prefix + label
		if m.WindowFocused && w.Id == window.Id {
			label = display.Style(label, "6")
		}

		lines = append(lines, label)
	}

	if m.Pane != "" {
		lines = append(lines, "", "Pane preview", "", m.Pane)
	}

	return strings.Join(lines, "\n")
}
