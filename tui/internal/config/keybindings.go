package config

import (
	"fmt"

	"multienv/tui/internal/keybindings"
)

func (c *Config) resolveKeybindings() error {
	common := keybindings.Common()
	c.Keybindings = keybindings.Merge(common, c.Keybindings)
	if err := c.Keybindings.Validate(common); err != nil {
		return fmt.Errorf("keybindings: %w", err)
	}

	scopes := []struct {
		name     string
		bindings *keybindings.Bindings
		defaults keybindings.Bindings
	}{
		{"tmux", &c.Tmux.Keybindings, keybindings.Tmux()},
		{"codex", &c.Codex.Keybindings, keybindings.Codex()},
		{"doctor", &c.Doctor.Keybindings, keybindings.Doctor()},
	}
	for _, scope := range scopes {
		supported := keybindings.Merge(common, scope.defaults)
		defaults := keybindings.Merge(c.Keybindings, scope.defaults)
		*scope.bindings = keybindings.Merge(defaults, *scope.bindings)
		if err := scope.bindings.Validate(supported); err != nil {
			return fmt.Errorf("%s.keybindings: %w", scope.name, err)
		}
	}

	return nil
}
