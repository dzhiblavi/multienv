package display

import (
	"strings"
	"time"
	"unicode"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func Clean(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}

		return r
	}, text)
}

func Style(text, color string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
}

func Indicator(attached bool, synced time.Time) string {
	if !synced.IsZero() && time.Since(synced) > 30*time.Second {
		return Style("!", "1")
	}

	if attached {
		return Style("●", "2")
	}

	return " "
}

func Age(synced time.Time) string {
	if synced.IsZero() {
		return "unknown"
	}

	return time.Since(synced).Round(time.Second).String() + " ago"
}

func Paragraph(text string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' {
			return r
		}

		if r == '\t' {
			return ' '
		}

		if unicode.IsControl(r) {
			return ' '
		}

		return r
	}, text)
}

func Message(text string, width int) string {
	return ansi.Wrap(Paragraph(text), max(1, width), "")
}
