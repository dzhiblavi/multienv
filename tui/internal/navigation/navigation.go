package navigation

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"multienv/common/log"
	"multienv/common/subprocess"
	"multienv/tui/internal/config"
)

type Action struct {
	Hostname string `json:"hostname,omitempty"`
	Target   string `json:"target,omitempty"`
	ResumeID string `json:"resume_id,omitempty"`
	WorkDir  string `json:"work_dir,omitempty"`
}

func Execute(ctx context.Context, cfg config.Config, action *Action) error {
	if action == nil {
		return nil
	}

	if action.ResumeID != "" {
		return Resume(action.ResumeID, action.WorkDir)
	}

	log.Info(ctx, "switching session", "host", action.Hostname, "target", action.Target)

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	if action.Hostname == hostname {
		return tmux(ctx, "switch-client", "-t", action.Target)
	}

	if err := sendToLauncher(ctx, cfg, *action); err != nil {
		return err
	}

	return tmux(ctx, "detach-client")
}

func sendToLauncher(ctx context.Context, cfg config.Config, action Action) error {
	if strings.ContainsAny(action.Hostname+action.Target, "\t\r\n") {
		return fmt.Errorf("invalid launcher target")
	}

	address := fmt.Sprintf("127.0.0.1:%d", cfg.Tmux.ControlPort)
	conn, err := (&net.Dialer{Timeout: 3 * time.Second}).DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connect to launcher %s: %w", address, err)
	}

	defer conn.Close()
	if err = conn.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
		return err
	}

	if _, err = io.WriteString(conn, action.Hostname+"\t"+action.Target+"\n"); err != nil {
		return err
	}

	return nil
}

func tmux(ctx context.Context, args ...string) error {
	_, err := subprocess.RunSuccess(ctx, subprocess.Command{Name: "tmux", Args: args})
	return err
}
