package app

import (
	"context"
	"time"

	tea "charm.land/bubbletea/v2"

	codexmodel "multienv/agent/domains/codex/model"
	tmuxmodel "multienv/agent/domains/tmux/model"
	doctormodel "multienv/common/doctor/model"
	"multienv/common/log"
	"multienv/tui/internal/doctorview"
)

type refreshTick struct{}
type snapshot struct {
	mode   Mode
	id     int
	tmux   tmuxmodel.SessionsByHost
	codex  codexmodel.SessionsByHost
	doctor doctormodel.Report
	err    error
}

type synced struct {
	mode Mode
	err  error
}

func tick() tea.Cmd { return tea.Tick(3*time.Second, func(time.Time) tea.Msg { return refreshTick{} }) }

func (m *model) load(force bool) tea.Cmd {
	mode := m.mode
	if m.loading[mode] && !force {
		return nil
	}

	if cancel := m.loadCancel[mode]; cancel != nil {
		cancel()
	}

	ctx, cancel := context.WithCancel(m.ctx)
	m.loadCancel[mode] = cancel
	m.loadID[mode]++
	id := m.loadID[mode]
	m.loading[mode] = true
	return func() tea.Msg {
		result := snapshot{mode: mode, id: id}
		switch mode {
		case Tmux:
			result.tmux, result.err = m.services.Tmux.ListAll(ctx)
		case Codex:
			result.codex, result.err = m.services.Codex.ListAll(ctx)
		case Doctor:
			result.doctor = doctorview.WithConfig(m.services.Diagnose(ctx), nil)
		}

		return result
	}
}

func (m *model) sync() tea.Cmd {
	mode := m.mode
	if mode == Doctor || m.syncing[mode] {
		return nil
	}

	m.syncing[mode] = true
	return func() tea.Msg {
		var err error
		if mode == Tmux {
			err = m.services.Tmux.SyncAll(m.ctx)
		} else {
			err = m.services.Codex.SyncAll(m.ctx)
		}

		return synced{mode: mode, err: err}
	}
}

func (m *model) setSnapshot(msg snapshot) tea.Cmd {
	if msg.id != m.loadID[msg.mode] {
		return nil
	}

	m.loading[msg.mode] = false
	if msg.err != nil {
		return m.error(msg.err)
	}

	switch msg.mode {
	case Tmux:
		m.tmux.Set(msg.tmux, m.cfg)
	case Codex:
		if err := m.codex.Set(msg.codex, m.cfg); err != nil {
			return m.error(err)
		}
	case Doctor:
		m.doctor.Set(msg.doctor)
	}

	m.updateContents()
	return m.requestPreview(true)
}

func (m *model) setSynced(msg synced) tea.Cmd {
	m.syncing[msg.mode] = false
	if msg.err != nil {
		log.Warn(m.ctx, "failed to sync sessions", "domain", msg.mode, "error", msg.err)
	}

	return nil
}
