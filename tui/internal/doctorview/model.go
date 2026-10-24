package doctorview

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"multienv/common/doctor/model"
	"multienv/tui/internal/display"
	"multienv/tui/internal/keybindings"
)

type Model struct {
	Keybindings keybindings.Bindings
	Search      textinput.Model
	Viewport    viewport.Model
	Report      model.Report
}

func New(bindings keybindings.Bindings) Model {
	input := textinput.New()
	input.Prompt = "Search checks ❯ "
	preview := viewport.New()
	preview.SoftWrap = true
	return Model{Search: input, Viewport: preview, Keybindings: bindings}
}

func (m *Model) Set(report model.Report) {
	m.Report = report
	m.filter()
}

func (m *Model) Resize(width, height int) {
	m.Search.SetWidth(width)
	m.Viewport.SetWidth(width)
	m.Viewport.SetHeight(height - 2)
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch m.Keybindings.Action(key.String()) {
		case "next":
			m.Viewport.ScrollDown(1)
			return nil
		case "previous":
			m.Viewport.ScrollUp(1)
			return nil
		case "preview-up":
			m.Viewport.HalfPageUp()
			return nil
		case "preview-down":
			m.Viewport.HalfPageDown()
			return nil
		}
	}

	previous := m.Search.Value()
	var cmd tea.Cmd
	m.Search, cmd = m.Search.Update(msg)
	if previous != m.Search.Value() {
		m.filter()
		m.Viewport.GotoTop()
	}

	return cmd
}

func (m *Model) filter() {
	lines := []string{}
	query := strings.ToLower(m.Search.Value())
	for _, host := range m.Report.Hosts {
		for _, check := range host.Checks {
			line := fmt.Sprintf("%s [%s]  %-7s %-12s %s: %s (%d ms)", host.Hostname, host.Source,
				strings.ToUpper(string(check.Status)), check.Group, check.ID, check.Message, check.DurationMS)
			if query != "" && !strings.Contains(strings.ToLower(line), query) {
				continue
			}

			color := map[model.Status]string{model.OK: "2", model.Warn: "3", model.Fail: "1", model.Skipped: "8"}[check.Status]
			lines = append(lines, display.Style(display.Clean(line), color))
		}
	}

	m.Viewport.SetContent(strings.Join(lines, "\n"))
}

func (m *Model) Content() string { return m.Search.View() + "\n\n" + m.Viewport.View() }

func Plain(report model.Report) string {
	var b strings.Builder
	for _, host := range report.Hosts {
		fmt.Fprintf(&b, "%s [%s]\n", host.Hostname, host.Source)
		for _, c := range host.Checks {
			fmt.Fprintf(&b, "  %-7s %-12s %s: %s (%d ms)\n",
				strings.ToUpper(string(c.Status)), c.Group, c.ID, c.Message, c.DurationMS)
		}

		b.WriteString("\n")
	}

	return b.String()
}

func WithConfig(report model.Report, err error) model.Report {
	hostname, _ := os.Hostname()
	check := model.Check{ID: "config", Group: "tui", Status: model.OK, Message: "TUI configuration loaded"}
	if err != nil {
		check.Status = model.Fail
		check.Message = err.Error()
	}

	host := model.HostReport{Hostname: hostname, Source: "tui", Checks: []model.Check{check}}
	report.Hosts = append([]model.HostReport{host}, report.Hosts...)
	return report
}
