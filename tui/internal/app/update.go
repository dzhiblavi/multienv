package app

import (
	tea "charm.land/bubbletea/v2"

	"multienv/tui/internal/keybindings"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case snapshot:
		cmd = m.setSnapshot(msg)
	case synced:
		cmd = m.setSynced(msg)
	case refreshTick:
		cmd = tick()
		if m.mode != Doctor {
			cmd = tea.Batch(cmd, m.load(false), m.sync())
		}
	case previewLoaded:
		cmd = m.setPreview(msg)
	case actionDone:
		cmd = m.actionFinished(msg)
	case directoryChoices:
		cmd = m.directoryDialog(msg)
	case codexReady:
		cmd = m.resumeDialog(msg)
	default:
		cmd = m.input(msg)
	}

	m.resize(m.width, m.height)
	m.updateContents()
	return m, tea.Batch(cmd, m.requestPreview(false))
}

func (m *model) input(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyPressMsg); ok && m.bindings().Action(key.String()) == "quit" {
		return tea.Quit
	}

	if m.modal != nil {
		return m.modalInput(msg)
	}

	if key, ok := msg.(tea.KeyPressMsg); ok {
		if cmd, handled := m.key(m.bindings().Action(key.String())); handled {
			return cmd
		}
	}

	switch m.mode {
	case Doctor:
		return m.doctor.Update(msg)
	case Codex:
		return m.codex.Browser.Update(msg)
	default:
		return m.tmux.Update(msg)
	}
}

func (m *model) modalInput(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch m.modal.Keybindings.Action(key.String()) {
		case "back":
			m.modal = nil
			m.submit = nil
			return nil
		case "open":
			result, ok := m.modal.Result()
			if !ok {
				return nil
			}

			submit := m.submit
			m.modal = nil
			m.submit = nil
			return submit(result)
		}
	}

	return m.modal.Update(msg)
}

func (m *model) key(action string) (tea.Cmd, bool) {
	switch action {
	case "back":
		return m.back(), true
	case "switch-view":
		return m.toggleMode(), true
	case "doctor":
		return m.toggleDoctor(), true
	case "help":
		m.help = !m.help
		return nil, true
	case "preview-up", "preview-down":
		if m.mode != Doctor {
			return m.scrollPreview(action), true
		}
	}

	if m.mode == Doctor {
		if action == "refresh" {
			return m.load(true), true
		}

		return nil, false
	}

	return m.domainKey(action)
}

func (m *model) domainKey(action string) (tea.Cmd, bool) {
	if m.busy {
		switch action {
		case "open", "create", "rename", "move-window", "delete", "archive", "unarchive":
			return nil, true
		}
	}

	if m.windowFocusKey(action) {
		return nil, true
	}

	if action == "open" {
		return m.activate(), true
	}

	if m.mode == Codex {
		return m.codexKeyAction(action)
	}

	switch action {
	case "create":
		return m.createTmux(), true
	case "rename":
		return m.renameTmux(), true
	case "move-window":
		if m.tmux.WindowFocused {
			return m.moveWindow(), true
		}
	case "delete":
		return m.killTmux(), true
	}

	return nil, false
}

func (m *model) codexKeyAction(action string) (tea.Cmd, bool) {
	switch action {
	case "archive":
		return m.codexAction("archive"), true
	case "unarchive":
		return m.codexAction("unarchive"), true
	case "delete":
		return m.codexAction("delete"), true
	case "toggle-archived":
		m.codex.ShowArchived = !m.codex.ShowArchived
	case "toggle-ignored":
		m.codex.ShowIgnored = !m.codex.ShowIgnored
	default:
		return nil, false
	}

	return m.error(m.codex.Set(m.codex.Snapshot, m.cfg)), true
}

func (m *model) back() tea.Cmd {
	if m.mode == Doctor {
		return m.activateMode(m.previous)
	}
	if m.mode == Tmux && m.tmux.WindowFocused {
		m.tmux.WindowFocused = false
		return nil
	}

	return tea.Quit
}

func (m *model) toggleMode() tea.Cmd {
	if m.mode == Doctor {
		return m.activateMode(m.previous)
	}
	if m.mode == Tmux {
		return m.activateMode(Codex)
	}
	return m.activateMode(Tmux)
}

func (m *model) toggleDoctor() tea.Cmd {
	if m.mode == Doctor {
		return m.activateMode(m.previous)
	}
	m.previous = m.mode
	return m.activateMode(Doctor)
}

func (m *model) windowFocusKey(action string) bool {
	if m.mode != Tmux {
		return false
	}
	if action == "focus-windows" {
		m.tmux.WindowFocused = true
		if window, ok := m.tmux.Window(); ok {
			m.tmux.WindowID = window.Id
		}
		return true
	}
	if action == "focus-sessions" && m.tmux.WindowFocused {
		m.tmux.WindowFocused = false
		return true
	}

	return false
}

func (m *model) bindings() keybindings.Bindings {
	switch m.mode {
	case Codex:
		return m.cfg.Codex.Keybindings
	case Doctor:
		return m.cfg.Doctor.Keybindings
	default:
		return m.cfg.Tmux.Keybindings
	}
}
