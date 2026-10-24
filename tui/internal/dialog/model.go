package dialog

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"multienv/tui/internal/browser"
	"multienv/tui/internal/keybindings"
)

type Result struct {
	Text string
	Item browser.Item
}

type Model struct {
	Keybindings keybindings.Bindings
	Title       string
	Note        string
	Choice      *browser.Model
	Input       textinput.Model
}

func Text(title, initial string, bindings keybindings.Bindings) *Model {
	input := textinput.New()
	input.Prompt = "❯ "
	input.SetValue(initial)
	input.CursorEnd()
	return &Model{Title: title, Input: input, Keybindings: bindings}
}

func Choose(title string, items []browser.Item, bindings keybindings.Bindings) *Model {
	choice := browser.New(false, bindings)
	choice.Set(items)
	return &Model{Title: title, Choice: &choice, Keybindings: bindings}
}

func (m *Model) Init() tea.Cmd {
	if m.Choice != nil {
		return m.Choice.Init()
	}

	return m.Input.Focus()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if m.Choice != nil {
		return m.Choice.Update(msg)
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return cmd
}

func (m *Model) Result() (Result, bool) {
	if m.Choice != nil {
		item, ok := m.Choice.Selected()
		return Result{Item: item}, ok
	}

	text := m.Input.Value()
	return Result{Text: text}, strings.TrimSpace(text) != ""
}

func (m *Model) View(width, height int) string {
	screenWidth := width
	width = max(1, min(90, width-6))
	content := m.Input.View()
	if m.Choice != nil {
		content = m.Choice.View()
	}

	content = m.Title + "\n\n" + content
	if m.Note != "" {
		content += "\n\n" + m.Note
	}

	content += "\n\n" + keybindings.Hints(m.Keybindings.Hint("open", "select"), m.Keybindings.Hint("back", "cancel"))
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Width(width).Render(content)
	return lipgloss.Place(screenWidth, height, lipgloss.Center, lipgloss.Center, box)
}

func (m *Model) Resize(width, height int) {
	width = max(1, min(90, width-6))
	m.Input.SetWidth(max(1, width-6))
	if m.Choice != nil {
		m.Choice.Resize(width-4, max(3, height-10))
	}
}
