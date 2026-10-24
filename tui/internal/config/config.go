package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"multienv/common/util"
	"multienv/tui/internal/keybindings"
)

const AgentPath = "~/.config/multienv/agent.yaml"
const DefaultPath = "~/.config/multienv/tui.yaml"

type Config struct {
	HostAliases map[string]string    `yaml:"host_aliases"`
	Popup       Popup                `yaml:"popup"`
	LogFile     string               `yaml:"log_file"`
	Keybindings keybindings.Bindings `yaml:"keybindings"`
	Tmux        Tmux                 `yaml:"tmux"`
	Codex       Codex                `yaml:"codex"`
	Doctor      Doctor               `yaml:"doctor"`
}

type Tmux struct {
	ControlPort       int                  `yaml:"control_port"`
	RemotePanePreview bool                 `yaml:"remote_pane_preview"`
	Keybindings       keybindings.Bindings `yaml:"keybindings"`
}

type Codex struct {
	IgnoreFile  string               `yaml:"ignore_file"`
	Keybindings keybindings.Bindings `yaml:"keybindings"`
}

type Doctor struct {
	Keybindings keybindings.Bindings `yaml:"keybindings"`
}

func (c *Config) Default() {
	c.HostAliases = map[string]string{}
	c.Tmux.ControlPort = 7777
	c.Popup = Popup{Enabled: true, Width: "90%", Height: "90%"}
	c.Codex.IgnoreFile = "~/.config/multienv/.multienv-tui-ignore"
	c.LogFile = filepath.Join(os.TempDir(), "multienv-tui.log")
}

func Expand(path string) (string, error) {
	path = os.ExpandEnv(path)
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	return path, nil
}

func Load(ctx context.Context, path string) (Config, error) {
	var cfg Config
	cfg.Default()

	path, err := Expand(path)
	if err != nil {
		return cfg, err
	}

	if _, err = os.Stat(path); err != nil && !os.IsNotExist(err) {
		return cfg, err
	}

	if err == nil {
		if err = util.ReadYaml(ctx, path, &cfg); err != nil {
			return cfg, err
		}
	}

	err = cfg.prepare()
	return cfg, err
}

func (c *Config) prepare() error {
	if err := c.validate(); err != nil {
		return err
	}

	if err := c.Popup.validate(); err != nil {
		return err
	}

	if err := c.resolveKeybindings(); err != nil {
		return err
	}

	var err error
	c.Codex.IgnoreFile, err = Expand(c.Codex.IgnoreFile)
	if err != nil {
		return err
	}

	c.LogFile, err = Expand(c.LogFile)
	return err
}

func (c Config) Alias(host string) string {
	for alias, hostname := range c.HostAliases {
		if hostname == host {
			return alias
		}
	}

	return host
}

func (c Config) validate() error {
	if c.Tmux.ControlPort < 1 || c.Tmux.ControlPort > 65535 {
		return fmt.Errorf("tmux.control_port must be between 1 and 65535")
	}

	for alias, host := range c.HostAliases {
		if alias == "" || host == "" || strings.ContainsAny(host, "\t\r\n") {
			return fmt.Errorf("invalid host alias %q", alias)
		}
	}

	return nil
}
