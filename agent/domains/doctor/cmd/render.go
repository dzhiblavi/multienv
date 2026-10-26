package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	doctor "multienv/common/doctor/model"
)

func render(out io.Writer, format string, report doctor.Report) error {
	if format == "json" {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	}

	for _, host := range report.Hosts {
		if _, err := fmt.Fprintf(out, "%s [%s]\n", singleLine(host.Hostname), singleLine(host.Source)); err != nil {
			return err
		}

		for _, check := range host.Checks {
			_, err := fmt.Fprintf(out, "  %-7s %-12s %s: %s (%d ms)\n",
				strings.ToUpper(string(check.Status)), check.Group, check.ID, singleLine(check.Message), check.DurationMS)

			if err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintln(out); err != nil {
			return err
		}
	}

	return nil
}

func singleLine(value string) string {
	return strings.Map(func(r rune) rune {
		if r < 32 || (r >= 127 && r <= 159) {
			return ' '
		}

		return r
	}, value)
}
