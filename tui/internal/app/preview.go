package app

import (
	"context"

	tea "charm.land/bubbletea/v2"

	tmuxservice "multienv/agent/domains/tmux/service"
)

type previewLoaded struct {
	key     string
	id      int
	content string
}

func (m *model) cancelPreview() {
	if m.previewCancel != nil {
		m.previewCancel()
		m.previewCancel = nil
	}

	m.previewKey = ""
	m.tmux.Pane = ""
}

func (m *model) previewRequest() (string, tmuxservice.PanePreviewParams, bool) {
	var req tmuxservice.PanePreviewParams
	if m.mode != Tmux {
		return "", req, false
	}

	session, ok := m.tmux.Selected()
	if !ok {
		return "", req, false
	}

	window, ok := m.tmux.Window()
	if !ok {
		return "", req, false
	}

	if session.Host != m.hostname && !m.cfg.Tmux.RemotePanePreview {
		return "", req, false
	}

	req.Hostname = session.Host
	req.Target = session.Data.Id + ":" + window.Id
	return session.Key() + "\x00" + window.Id, req, true
}

func (m *model) requestPreview(force bool) tea.Cmd {
	key, req, ok := m.previewRequest()
	if !ok {
		m.cancelPreview()
		return nil
	}

	if !force && key == m.previewKey {
		return nil
	}

	if m.previewCancel != nil {
		m.previewCancel()
	}

	if key != m.previewKey {
		m.tmux.Pane = ""
	}

	m.previewKey = key
	m.previewID++
	id := m.previewID
	ctx, cancel := context.WithCancel(m.ctx)
	m.previewCancel = cancel
	return func() tea.Msg {
		response, err := m.services.Tmux.PanePreview(ctx, req)
		content := ""
		if err == nil {
			content = response.Content
		}

		return previewLoaded{key: key, id: id, content: content}
	}
}

func (m *model) setPreview(msg previewLoaded) tea.Cmd {
	if msg.key == m.previewKey && msg.id == m.previewID {
		m.tmux.Pane = msg.content
	}

	return nil
}

func (m *model) updateContents() {
	m.tmux.Preview.SetContent(m.tmux.PreviewText())
	if session, ok := m.codex.Selected(); ok {
		key := session.Hostname + "\x00" + session.Id
		if key != m.codexKey {
			m.codex.Preview.GotoTop()
			m.codexKey = key
		}
	}

	m.codex.Preview.SetContent(m.codex.PreviewText(m.cfg))
}

func (m *model) scrollPreview(action string) tea.Cmd {
	v := &m.tmux.Preview
	if m.mode == Codex {
		v = &m.codex.Preview
	}

	if action == "preview-up" {
		v.HalfPageUp()
	} else {
		v.HalfPageDown()
	}

	return nil
}
