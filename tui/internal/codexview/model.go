package codexview

import (
	"cmp"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/viewport"

	"multienv/agent/codex/model"
	"multienv/tui/internal/browser"
	"multienv/tui/internal/config"
	"multienv/tui/internal/display"
)

type Model struct {
	Browser      browser.Model
	Preview      viewport.Model
	Snapshot     model.SessionsByHost
	ShowArchived bool
	ShowIgnored  bool
}

func New(cfg config.Codex) Model {
	preview := viewport.New()
	preview.SoftWrap = true
	return Model{Browser: browser.New(true, cfg.Keybindings), Preview: preview}
}

func (m *Model) Selected() (model.Session, bool) {
	item, ok := m.Browser.Selected()
	if !ok {
		return model.Session{}, false
	}

	session, ok := item.Value.(model.Session)
	return session, ok
}

func State(s model.Session) string {
	if s.State == model.Working && s.TmuxPaneId == "" {
		return "Working (bg)"
	}

	return string(s.State)
}

func (m *Model) Set(snapshot model.SessionsByHost, cfg config.Config) error {
	var patterns []*regexp.Regexp
	if !m.ShowIgnored {
		var err error
		patterns, err = ignorePatterns(cfg.Codex.IgnoreFile)
		if err != nil {
			return err
		}
	}

	m.Snapshot = snapshot
	items := []browser.Item{}
	for _, session := range sortedSessions(snapshot) {
		if m.visible(session, cfg, patterns) {
			items = append(items, sessionItem(session, cfg))
		}
	}

	m.Browser.Set(items)
	return nil
}

func sortedSessions(snapshot model.SessionsByHost) []model.Session {
	sessions := []model.Session{}
	for host, data := range snapshot {
		for id, session := range data {
			session.Hostname = host
			session.Id = id
			sessions = append(sessions, session)
		}
	}

	slices.SortFunc(sessions, func(a, b model.Session) int {
		return cmp.Or(cmp.Compare(statePriority(a.State), statePriority(b.State)),
			cmp.Compare(a.Hostname, b.Hostname), cmp.Compare(sessionTitle(a), sessionTitle(b)), cmp.Compare(a.Id, b.Id))
	})
	return sessions
}

func statePriority(state model.State) int {
	switch state {
	case model.ActionRequired:
		return 0
	case model.Idle:
		return 1
	case model.Working:
		return 2
	case model.NotLoaded:
		return 3
	case model.Archived:
		return 5
	default:
		return 4
	}
}

func (m *Model) visible(session model.Session, cfg config.Config, patterns []*regexp.Regexp) bool {
	if !m.ShowArchived && session.State == model.Archived {
		return false
	}

	plain := State(session) + "\t" + cfg.Alias(session.Hostname) + "\t" + sessionTitle(session)
	ignored := slices.ContainsFunc(patterns, func(p *regexp.Regexp) bool { return p.MatchString(plain) })
	return m.ShowIgnored || !ignored
}

func sessionItem(session model.Session, cfg config.Config) browser.Item {
	title := sessionTitle(session)
	state := display.Style(display.Clean(State(session)), stateColor(session.State))
	host := display.HostAlias(session.Hostname, cfg.Alias(session.Hostname))
	label := state + " · " + host + " · " + display.Clean(title)
	detail := display.Clean(session.WorkDir)
	indicator := display.Indicator(session.TmuxPaneId != "", session.SyncedAt)

	return browser.Item{
		ID:        session.Hostname + "\x00" + session.Id,
		Label:     indicator + " " + label,
		Detail:    detail,
		Search:    display.Clean(title + " " + cfg.Alias(session.Hostname)),
		AllSearch: display.Clean(State(session) + " " + cfg.Alias(session.Hostname) + " " + title + " " + detail),
		Value:     session,
	}
}

func sessionTitle(session model.Session) string {
	if session.Title != "" {
		return session.Title
	}

	return session.Id
}

func ignorePatterns(path string) ([]*regexp.Regexp, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var patterns []*regexp.Regexp
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		pattern, err := regexp.Compile(line)
		if err != nil {
			return nil, fmt.Errorf("ignore rule %q: %w", line, err)
		}

		patterns = append(patterns, pattern)
	}

	return patterns, nil
}

func (m *Model) PreviewText(cfg config.Config) string {
	s, ok := m.Selected()
	if !ok {
		return "No conversation selected"
	}

	title := s.Title
	if title == "" {
		title = s.Id
	}

	heading := func(text string) string {
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")).Render(display.Clean(text))
	}

	lines := []string{
		heading(title),
		display.Style(display.Clean(State(s)), stateColor(s.State)),
		"",
		display.Style("Host       ", "8") + display.HostAlias(s.Hostname, cfg.Alias(s.Hostname)),
		display.Style("Directory  ", "8") + display.Clean(s.WorkDir),
		display.Style("Pane       ", "8") + display.Clean(s.TmuxPaneId),
		display.Style("Synced     ", "8") + display.Age(s.SyncedAt),
		"", heading("First message"), "", display.Message(s.FirstMessage, m.Preview.Width()),
		"", heading("Last message"), "", display.Message(s.LastMessage, m.Preview.Width()),
	}

	return strings.Join(lines, "\n")
}

func stateColor(state model.State) string {
	switch state {
	case model.Working:
		return "6"
	case model.Idle:
		return "2"
	case model.ActionRequired:
		return "3"
	case model.Archived, model.Unknown:
		return "8"
	default:
		return "7"
	}
}
