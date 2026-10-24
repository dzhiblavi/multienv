package navigation

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"multienv/common/subprocess"
)

func Resume(id, workDir string) error {
	executable, err := exec.LookPath("codex")
	if err != nil {
		return err
	}

	args := []string{"codex", "--no-daemon"}
	if workDir != "" {
		args = append(args, "--cd", workDir)
	}

	args = append(args, "resume", id)
	return syscall.Exec(executable, args, os.Environ())
}

func ResumeWhenVisible(ctx context.Context, id string) error {
	pane := os.Getenv("TMUX_PANE")
	if pane == "" {
		return fmt.Errorf("run inside a tmux pane")
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		result, err := subprocess.RunSuccess(ctx, subprocess.Command{Name: "tmux", Args: []string{
			"display-message", "-p", "-t", pane, "#{pane_active}:#{window_active}:#{session_attached}",
		}})
		if err != nil {
			return err
		}

		fields := strings.Split(strings.TrimSpace(result.Stdout), ":")
		if len(fields) == 3 && fields[0] == "1" && fields[1] == "1" && fields[2] != "0" {
			return Resume(id, "")
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for visible pane: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}
