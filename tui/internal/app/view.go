package app

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"multienv/tui/internal/display"
	"multienv/tui/internal/keybindings"
)

func (m *model) resize(width, height int) {
	m.width, m.height = max(1, width), max(1, height)
	if m.modal != nil {
		m.modal.Resize(width, height)
	}
	left := max(1, (width-6)/2)
	right := max(1, width-6-left)
	body := m.bodyHeight()
	m.tmux.Browser.Resize(left, body)
	m.codex.Browser.Resize(left, body)
	m.tmux.Preview.SetWidth(right)
	m.tmux.Preview.SetHeight(body)
	m.codex.Preview.SetWidth(right)
	m.codex.Preview.SetHeight(body)
	m.doctor.Resize(max(1, width-4), body)
}

func (m *model) View() tea.View {
	if m.modal != nil {
		view := tea.NewView(m.modal.View(m.width, m.height))
		view.AltScreen = true
		return view
	}

	title := fmt.Sprintf("Multienv · %s", m.mode)
	if m.mode == Codex {
		title += fmt.Sprintf(" · archived %t · ignored %t", m.codex.ShowArchived, m.codex.ShowIgnored)
	}

	activity := " "
	if m.busy {
		activity = "Working…"
	} else if m.loading[m.mode] {
		activity = "Refreshing…"
	}

	body := m.body()
	content := display.Style(title, "6") + "\n" + display.Style(activity, "8") + "\n" + body + "\n" + m.footer()
	view := tea.NewView(lipgloss.NewStyle().Margin(1, 2).MaxWidth(m.width).MaxHeight(m.height).Render(content))
	view.AltScreen = true
	return view
}

func (m *model) body() string {
	if m.mode == Doctor {
		return m.doctor.Content()
	}

	left := m.browser().View()
	right := m.tmux.Preview.View()
	if m.mode == Codex {
		right = m.codex.Preview.View()
	}

	width := max(1, (m.width-6)/2)
	height := m.bodyHeight()
	panel := lipgloss.NewStyle().Width(width).MaxWidth(width).Height(height).MaxHeight(height)
	separator := display.Style(" │ ", "8")
	return lipgloss.JoinHorizontal(lipgloss.Top, panel.Render(left), separator, panel.Render(right))
}

func (m *model) helpText() string {
	keys := m.bindings()
	actions := []string{"open", "previous", "next", "mark-next", "mark-previous", "search-all",
		"preview-up", "preview-down", "quit"}
	switch m.mode {
	case Doctor:
		actions = []string{"previous", "next", "preview-up", "preview-down", "refresh", "quit"}
	case Codex:
		actions = append(actions, "rename", "archive", "unarchive", "delete", "toggle-archived", "toggle-ignored")
	default:
		actions = append(actions, "focus-windows", "focus-sessions", "previous-window", "next-window",
			"create", "rename", "move-window", "delete")
	}

	hints := make([]string, 0, len(actions))
	for _, action := range actions {
		hints = append(hints, keys.Hint(action, action))
	}

	return keybindings.Hints(hints...)
}

func (m *model) bodyHeight() int {
	return max(3, m.height-4-lipgloss.Height(m.footer()))
}

func (m *model) footer() string {
	keys := m.bindings()
	footer := keybindings.Hints(keys.Hint("switch-view", "tmux/codex"), keys.Hint("doctor", "doctor"),
		keys.Hint("help", "help"), keys.Hint("back", "back/quit"))
	if m.help {
		footer += "\n" + m.helpText()
	}

	if m.status != "" {
		footer += "\n" + display.Style(display.Clean(m.status), "1")
	}

	return ansi.Wrap(footer, max(1, m.width-4), "")
}
