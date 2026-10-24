package navigation

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"multienv/tui/internal/config"
)

func Popup(ctx context.Context, settings config.Popup, args []string) (*Action, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}

	directory, err := os.MkdirTemp("", "multienv-tui.")
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(directory)
	resultPath := filepath.Join(directory, "result.json")
	args = append(args, "--no-popup", "--popup-child", "--result-file", resultPath)
	command := shellQuote(executable)
	for _, arg := range args {
		command += " " + shellQuote(arg)
	}

	cmd := exec.CommandContext(ctx, "tmux", "display-popup", "-E", "-w", settings.Width, "-h", settings.Height, command)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("open popup: %w", err)
	}

	data, err := os.ReadFile(resultPath)
	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var action Action
	if err = json.Unmarshal(data, &action); err != nil {
		return nil, err
	}

	return &action, nil
}

func WriteResult(path string, action *Action) error {
	if action == nil {
		return nil
	}

	data, err := json.Marshal(action)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func shellQuote(value string) string { return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'" }
