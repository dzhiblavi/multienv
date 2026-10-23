package doctor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"multienv/common/doctor/model"
	"multienv/common/subprocess"
)

type AgentRunner struct{}

func (AgentRunner) Diagnose(ctx context.Context, hostname string) (model.Report, error) {
	command := subprocess.Command{
		Name: "~/.bin/multienv-agent", Args: []string{"doctor", "local", "--format", "json"},
	}

	localHostname, err := os.Hostname()
	if err != nil {
		return model.Report{}, err
	}

	var result subprocess.RunResult
	if hostname == localHostname {
		result, err = subprocess.Run(ctx, command)
	} else {
		result, err = subprocess.RunRemote(ctx, subprocess.RemoteCommand{Hostname: hostname, Command: command})
	}
	if err != nil {
		return model.Report{}, err
	}

	if err := ctx.Err(); err != nil {
		return model.Report{}, err
	}

	var report model.Report
	if err := json.Unmarshal([]byte(result.Stdout), &report); err != nil {
		return report, fmt.Errorf("agent report unavailable: %w; stderr: %s", err, result.Stderr)
	}

	if len(report.Hosts) != 1 || report.Hosts[0].Hostname != hostname || len(report.Hosts[0].Checks) == 0 {
		return model.Report{}, fmt.Errorf("invalid local report for %s", hostname)
	}

	if result.ExitCode != 0 && (result.ExitCode != 1 || !report.Failed()) {
		return model.Report{}, fmt.Errorf("agent exited with code %d: %s", result.ExitCode, result.Stderr)
	}

	return report, nil
}
