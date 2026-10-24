package app

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"multienv/tui/internal/backend"
	"multienv/tui/internal/browser"
	"multienv/tui/internal/codexview"
	"multienv/tui/internal/config"
	"multienv/tui/internal/dialog"
	"multienv/tui/internal/doctorview"
	"multienv/tui/internal/navigation"
	"multienv/tui/internal/tmuxview"
)

type Mode string

const (
	Tmux   Mode = "tmux"
	Codex  Mode = "codex"
	Doctor Mode = "doctor"
)

type Options struct {
	Mode          Mode
	ResumeInPlace bool
}

type model struct {
	ctx           context.Context
	services      backend.Services
	cfg           config.Config
	options       Options
	hostname      string
	mode          Mode
	previous      Mode
	tmux          tmuxview.Model
	codex         codexview.Model
	doctor        doctorview.Model
	width, height int
	modal         *dialog.Model
	submit        func(dialog.Result) tea.Cmd
	busy          bool
	status        string
	help          bool
	result        *navigation.Action
	loading       map[Mode]bool
	syncing       map[Mode]bool
	loadID        map[Mode]int
	loadCancel    map[Mode]context.CancelFunc
	previewCancel context.CancelFunc
	previewID     int
	previewKey    string
	codexKey      string
}

func Run(
	ctx context.Context,
	services backend.Services,
	cfg config.Config,
	options Options,
) (*navigation.Action, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	m := &model{ctx: ctx, services: services, cfg: cfg, options: options, hostname: hostname,
		mode: options.Mode, previous: Tmux, tmux: tmuxview.New(cfg.Tmux), codex: codexview.New(cfg.Codex),
		doctor:  doctorview.New(cfg.Doctor.Keybindings),
		loading: map[Mode]bool{}, syncing: map[Mode]bool{}, loadID: map[Mode]int{}, loadCancel: map[Mode]context.CancelFunc{}}
	m.resize(100, 30)
	_, err = tea.NewProgram(m, tea.WithContext(ctx)).Run()
	return m.result, err
}

func (m *model) browser() *browser.Model {
	if m.mode == Codex {
		return &m.codex.Browser
	}

	return &m.tmux.Browser
}

func (m *model) activateMode(mode Mode) tea.Cmd {
	m.mode = mode
	m.status = ""
	m.cancelPreview()
	var focus tea.Cmd
	if mode == Doctor {
		focus = m.doctor.Search.Focus()
	} else {
		focus = m.browser().Init()
	}

	return tea.Batch(focus, m.load(true), m.sync(), m.requestPreview(true))
}

func (m *model) Init() tea.Cmd { return tea.Batch(m.activateMode(m.mode), tick()) }

func (m *model) error(err error) tea.Cmd {
	if err != nil {
		m.status = err.Error()
	}

	return nil
}

func (m *model) choose(title string, items []browser.Item, submit func(dialog.Result) tea.Cmd) tea.Cmd {
	if len(items) == 0 {
		return m.error(fmt.Errorf("no choices available"))
	}

	m.modal = dialog.Choose(title, items, m.bindings())
	m.submit = submit
	return m.modal.Init()
}

func (m *model) text(title, initial string, submit func(dialog.Result) tea.Cmd) tea.Cmd {
	m.modal = dialog.Text(title, initial, m.bindings())
	m.submit = submit
	return m.modal.Init()
}

func (m *model) finish(action navigation.Action) tea.Cmd {
	m.result = &action
	return tea.Quit
}
