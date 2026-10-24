package keybindings

import (
	"fmt"
	"slices"
	"strings"
)

type Bindings map[string][]string

func Merge(base, overrides Bindings) Bindings {
	merged := Bindings{}
	for action, keys := range base {
		merged[action] = slices.Clone(keys)
	}

	for action, keys := range overrides {
		merged[action] = slices.Clone(keys)
	}

	return merged
}

func (b Bindings) Action(input string) string {
	for action, keys := range b {
		if slices.Contains(keys, input) {
			return action
		}
	}

	return ""
}

func (b Bindings) Hint(action, label string) string {
	if len(b[action]) == 0 {
		return ""
	}

	return strings.Join(b[action], "/") + " " + label
}

func Hints(hints ...string) string {
	hints = slices.DeleteFunc(hints, func(hint string) bool { return hint == "" })
	return strings.Join(hints, " · ")
}

func (b Bindings) Validate(supported Bindings) error {
	actions := make([]string, 0, len(b))
	for action := range b {
		actions = append(actions, action)
	}

	slices.Sort(actions)
	used := map[string]string{}
	for _, action := range actions {
		if _, ok := supported[action]; !ok {
			return fmt.Errorf("unknown action %q", action)
		}

		for _, key := range b[action] {
			if key == "" || key != strings.TrimSpace(key) || key != strings.ToLower(key) {
				return fmt.Errorf("invalid key %q for %q: use lowercase key names", key, action)
			}

			if previous, ok := used[key]; ok {
				return fmt.Errorf("key %q is assigned to both %q and %q", key, previous, action)
			}

			used[key] = action
		}
	}

	return nil
}
