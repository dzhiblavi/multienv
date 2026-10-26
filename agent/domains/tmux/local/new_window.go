package local

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"multienv/common/subprocess"
)

func NewWindow(
	ctx context.Context,
	sessionId, name, workDir string,
	command []string,
) (string, error) {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"new-window", "-d", "-P", "-F", "#{pane_id}",
			"-t", sessionId + ":",
		},
	}

	if name != "" {
		cmd.Args = append(cmd.Args, "-n", name)
	}

	if workDir != "" {
		info, err := os.Stat(workDir)
		if err != nil {
			return "", fmt.Errorf("check window directory: %w", err)
		}
		if !info.IsDir() {
			return "", fmt.Errorf("window directory is not a directory: %s", workDir)
		}
		cmd.Args = append(cmd.Args, "-c", workDir)
	}
	if len(command) > 0 {
		program, err := exec.LookPath(command[0])
		if err != nil {
			return "", fmt.Errorf("find window command: %w", err)
		}
		// Use the agent's environment rather than the tmux server's stale PATH.
		cmd.Args = append(cmd.Args, "-e", "PATH="+os.Getenv("PATH"))
		if codexHome, exists := os.LookupEnv("CODEX_HOME"); exists {
			cmd.Args = append(cmd.Args, "-e", "CODEX_HOME="+codexHome)
		}
		cmd.Args = append(cmd.Args, "--", program)
		cmd.Args = append(cmd.Args, command[1:]...)
	}

	result, err := subprocess.RunSuccess(ctx, cmd)
	if err != nil {
		return "", fmt.Errorf("create tmux window: %w", err)
	}
	return strings.TrimSpace(result.Stdout), nil
}
