package browser

import (
	"slices"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"multienv/tui/internal/display"
	"multienv/tui/internal/keybindings"
)

type Item struct {
	ID        string
	Label     string
	Detail    string
	Search    string
	AllSearch string
	Value     any
	marked    bool
}

func (i Item) Title() string {
	prefix := ""
	if i.marked {
		prefix = "✓ "
	}

	return prefix + i.Label
}

func (i Item) Description() string { return i.Detail }

func (i Item) FilterValue() string { return i.Search }

type Model struct {
	Keybindings keybindings.Bindings
	Input       textinput.Model
	List        list.Model
	Items       []Item
	Marks       map[string]bool
	AllFields   bool
	Multi       bool
}

func New(multi bool, bindings keybindings.Bindings) Model {
	input := textinput.New()
	input.Prompt = "Search ❯ "
	l := list.New(nil, list.NewDefaultDelegate(), 40, 16)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	return Model{Input: input, List: l, Marks: map[string]bool{}, Multi: multi, Keybindings: bindings}
}

func (m *Model) Init() tea.Cmd { return m.Input.Focus() }

func (m *Model) Resize(width, height int) {
	m.Input.SetWidth(max(1, width-2))
	m.List.SetSize(max(1, width), max(1, height-2))
}

func (m *Model) Selected() (Item, bool) {
	item, ok := m.List.SelectedItem().(Item)
	return item, ok
}

func (m *Model) Move(delta int) {
	if delta < 0 {
		m.List.CursorUp()
	} else {
		m.List.CursorDown()
	}
}

func (m *Model) Set(items []Item) {
	m.Items = items
	for id := range m.Marks {
		if !slices.ContainsFunc(items, func(i Item) bool { return i.ID == id }) {
			delete(m.Marks, id)
		}
	}

	m.Filter(true)
}

func (m *Model) Filter(preserve bool) {
	selected, _ := m.Selected()
	items := m.filtered()
	m.List.SetItems(items)
	m.List.ResetSelected()
	if preserve {
		for i, item := range items {
			if item.(Item).ID == selected.ID {
				m.List.Select(i)
				break
			}
		}
	}
}

func (m *Model) filtered() []list.Item {
	targets := make([]string, len(m.Items))
	for i, item := range m.Items {
		targets[i] = item.Search
		if m.AllFields {
			targets[i] = item.AllSearch
		}
	}

	ranks := list.DefaultFilter(m.Input.Value(), targets)
	if m.Input.Value() == "" {
		ranks = make([]list.Rank, len(m.Items))
		for i := range ranks {
			ranks[i].Index = i
		}
	}

	result := make([]list.Item, 0, len(ranks))
	for _, rank := range ranks {
		item := m.Items[rank.Index]
		item.marked = m.Marks[item.ID]
		result = append(result, item)
	}

	return result
}

func (m *Model) Targets() []Item {
	if len(m.Marks) == 0 {
		if item, ok := m.Selected(); ok {
			return []Item{item}
		}

		return nil
	}

	var items []Item
	for _, item := range m.Items {
		if m.Marks[item.ID] {
			items = append(items, item)
		}
	}

	return items
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		action := m.Keybindings.Action(key.String())
		switch action {
		case "previous":
			m.Move(-1)
			return nil
		case "next":
			m.Move(1)
			return nil
		case "mark-next", "mark-previous":
			if m.Multi {
				m.mark()
				if action == "mark-next" {
					m.Move(1)
				} else {
					m.Move(-1)
				}
			}

			return nil
		case "search-all":
			m.AllFields = !m.AllFields
			m.Filter(true)
			return nil
		}
	}

	previous := m.Input.Value()
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	if previous != m.Input.Value() {
		m.Filter(false)
	}

	return cmd
}

func (m *Model) mark() {
	if item, ok := m.Selected(); ok {
		if m.Marks[item.ID] {
			delete(m.Marks, item.ID)
		} else {
			m.Marks[item.ID] = true
		}

		m.Filter(true)
	}
}

func (m *Model) ClearMarks() {
	m.Marks = map[string]bool{}
	m.Filter(true)
}

func (m *Model) View() string {
	prompt := m.Input.View()
	if m.AllFields {
		prompt = display.Style("[all] ", "6") + prompt
	}

	return prompt + "\n\n" + m.List.View()
}
