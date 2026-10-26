package app

import (
	"context"

	tea "charm.land/bubbletea/v2"

	codexservice "multienv/agent/domains/codex/service"
	"multienv/tui/internal/dialog"
)

func (m *model) renameCodex() tea.Cmd {
	session, ok := m.codex.Selected()
	if !ok {
		return nil
	}

	return m.text("Conversation name", session.Title, func(result dialog.Result) tea.Cmd {
		return m.performOn([]string{session.Hostname}, func(ctx context.Context) error {
			return m.services.Codex.RenameSession(ctx, codexservice.RenameSessionParams{
				Hostname: session.Hostname, SessionId: session.Id, NewName: result.Text,
			})
		})
	})
}
